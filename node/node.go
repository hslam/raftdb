package node

import (
	"fmt"
	"github.com/hslam/handler/proxy"
	"github.com/hslam/handler/render"
	"github.com/hslam/mux"
	"github.com/hslam/raft"
	"github.com/hslam/rpc"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

const (
	network             = "tcp"
	codec               = "pb"
	MaxConnsPerHost     = 1
	MaxIdleConnsPerHost = 0
)

var (
	setCommandPool *sync.Pool
)

func init() {
	setCommandPool = &sync.Pool{
		New: func() interface{} {
			return &SetCommand{}
		},
	}
}

const LeaderPrefix = "LEADER:"

type Node struct {
	mu           sync.RWMutex
	host         string
	httpPort     int
	rpcPort      int
	dataDir      string
	mux          *mux.Mux
	render       *render.Render
	raftNode     raft.Node
	httpServer   *http.Server
	rpcServer    *rpc.Server
	rpcTransport *rpc.Transport
	db           *DB
}

func NewNode(dataDir string, host string, httpPort, rpcPort, raftPort int, peers []string, join string) *Node {
	n := &Node{
		host:     host,
		httpPort: httpPort,
		rpcPort:  rpcPort,
		dataDir:  dataDir,
		db:       newDB(),
		mux:      mux.New(),
		render:   render.NewRender(),
	}
	var err error
	nodes := make([]*raft.NodeInfo, len(peers))
	for i, v := range peers {
		nodes[i] = &raft.NodeInfo{Address: v, Data: nil}
	}
	n.raftNode, err = raft.NewNode(host, raftPort, n.dataDir, n.db, false, nodes)
	if err != nil {
		log.Fatal(err)
	}
	raft.SetLogLevel(0)
	n.raftNode.RegisterCommand(&SetCommand{})
	n.raftNode.SetSnapshot(NewSnapshot(n.db))
	n.raftNode.SetSyncTypes([]*raft.SyncType{
		{86400, 1},
		{14400, 1000},
		{3600, 10000},
		{1800, 50000},
		{900, 200000},
		{60, 5000000},
	})
	n.raftNode.SetCodec(&raft.GOGOPBCodec{})
	n.raftNode.SetGzipSnapshot(true)
	n.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", n.httpPort),
		Handler: n.mux,
	}
	n.mux.Group("/cluster", func(m *mux.Mux) {
		m.HandleFunc("/status", n.statusHandler).All()
		m.HandleFunc("/leader", n.leaderHandler).All()
		m.HandleFunc("/ready", n.readyHandler).All()
		m.HandleFunc("/address", n.addressHandler).All()
		m.HandleFunc("/isleader", n.isLeaderHandler).All()
		m.HandleFunc("/peers", n.peersHandler).All()
		m.HandleFunc("/nodes", n.nodesHandler).All()
	})
	n.mux.HandleFunc("/db/:key", n.leaderHandle(n.getHandler)).GET()
	n.mux.HandleFunc("/db/:key", n.leaderHandle(n.setHandler)).POST()
	return n
}

func (n *Node) ListenAndServe() error {
	n.raftNode.Start()
	log.Println("HTTP listening at:", n.uri())
	log.Println("RPC listening at:", fmt.Sprintf("%s:%d", n.host, n.rpcPort))
	service := new(Service)
	service.node = n
	n.rpcServer = rpc.NewServer()
	n.rpcServer.RegisterName("S", service)
	n.rpcServer.SetPoll(true)
	rpc.SetLogLevel(rpc.OffLogLevel)
	if n.rpcTransport == nil {
		n.InitRPCProxy(MaxConnsPerHost, MaxIdleConnsPerHost)
	}
	go func() {
		fmt.Println("raftdb.node.rpc :", n.rpcServer.Listen("tcp", fmt.Sprintf(":%d", n.rpcPort), codec))
	}()
	return n.httpServer.ListenAndServe()
}
func (n *Node) InitRPCProxy(MaxConnsPerHost int, MaxIdleConnsPerHost int) {
	n.rpcTransport = &rpc.Transport{
		MaxConnsPerHost:     MaxConnsPerHost,
		MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		Options:             &rpc.Options{Network: network, Codec: codec},
	}
}
func (n *Node) uri() string {
	return fmt.Sprintf("http://%s:%d", n.host, n.httpPort)
}

func (n *Node) leaderHandle(hander http.HandlerFunc) http.HandlerFunc {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	return func(w http.ResponseWriter, r *http.Request) {
		if n.raftNode.IsLeader() {
			hander(w, r)
		} else {
			leader := n.raftNode.Leader()
			if leader != "" {
				leader_url, err := url.Parse("http://" + leader)
				if err != nil {
					panic(err)
				}
				port, err := strconv.Atoi(leader_url.Port())
				if err != nil {
					panic(err)
				}
				leader_url.Host = leader_url.Hostname() + ":" + strconv.Itoa(port-2000)
				proxy.Proxy(w, r, leader_url.String()+r.URL.Path)
			}
		}
	}
}

func (n *Node) setHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	params := n.mux.Params(req)
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value := string(b)
	setCommand := newSetCommand(params["key"], value)
	_, err = n.raftNode.Do(setCommand)
	setCommandPool.Put(setCommand)
	if err != nil {
		if err == raft.ErrNotLeader {
			leader := n.raftNode.Leader()
			if leader != "" {
				http.Error(w, LeaderPrefix+leader, http.StatusBadRequest)
				return
			}
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
func (n *Node) getHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	params := n.mux.Params(req)
	if ok := n.raftNode.ReadIndex(); ok {
		value := n.db.Get(params["key"])
		w.Write([]byte(value))
	}
}
func (n *Node) statusHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	status := &Status{
		IsLeader: n.raftNode.IsLeader(),
		Leader:   n.raftNode.Leader(),
		Node:     n.raftNode.Address(),
		Peers:    n.raftNode.Peers(),
	}
	n.render.JSON(w, req, status, http.StatusOK)
}
func (n *Node) leaderHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	w.Write([]byte(n.raftNode.Leader()))
}
func (n *Node) readyHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	n.render.JSON(w, req, n.raftNode.Ready(), http.StatusOK)
}
func (n *Node) isLeaderHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	n.render.JSON(w, req, n.raftNode.IsLeader(), http.StatusOK)
}
func (n *Node) addressHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	w.Write([]byte(n.raftNode.Address()))
}
func (n *Node) peersHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	n.render.JSON(w, req, n.raftNode.Peers(), http.StatusOK)
}
func (n *Node) nodesHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	nodes := n.raftNode.Peers()
	nodes = append(nodes, n.raftNode.Address())
	n.render.JSON(w, req, nodes, http.StatusOK)
}
