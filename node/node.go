// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package node

import (
	"encoding/json"
	"fmt"
	"github.com/hslam/handler/proxy"
	"github.com/hslam/handler/render"
	"github.com/hslam/log"
	"github.com/hslam/raft"
	"github.com/hslam/rpc"
	"github.com/hslam/rum"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	network             = "tcp"
	codec               = "pb"
	MaxConnsPerHost     = 1
	MaxIdleConnsPerHost = 1
)

var (
	setCommandPool *sync.Pool
	logger         = log.New()
)

func init() {
	setCommandPool = &sync.Pool{
		New: func() interface{} {
			return &SetCommand{}
		},
	}
	logger.SetPrefix("raftdb")
}

const LeaderPrefix = "LEADER:"

type Node struct {
	mu           sync.RWMutex
	host         string
	httpPort     int
	rpcPort      int
	dataDir      string
	render       *render.Render
	raftNode     raft.Node
	rum          *rum.Rum
	rpcServer    *rpc.Server
	rpcTransport *rpc.Transport
	db           *DB
	leader       *address
}

type address struct {
	HTTP string `json:"h,omitempty"`
	RPC  string `json:"r,omitempty"`
}

func NewNode(dataDir string, host string, httpPort, rpcPort, raftPort int, members []string, join bool) *Node {
	n := &Node{
		host:     host,
		httpPort: httpPort,
		rpcPort:  rpcPort,
		dataDir:  dataDir,
		db:       newDB(),
		rum:      rum.New(),
		render:   render.NewRender(),
	}
	n.InitRPCProxy(1, 0)
	var err error
	m := make([]*raft.Member, 0, len(members))
	for _, v := range members {
		if len(v) > 0 {
			m = append(m, &raft.Member{Address: v})
		}
	}
	n.raftNode, err = raft.NewNode(host, raftPort, n.dataDir, n.db, join, m)
	if err != nil {
		logger.Fatal(err)
	}
	n.raftNode.SetLogLevel(raft.AllLogLevel)
	n.raftNode.RegisterCommand(&SetCommand{})
	n.raftNode.SetSnapshot(NewSnapshot(n.db))
	n.raftNode.SetSyncTypes([]*raft.SyncType{
		{86400, 1},
		{14400, 1000},
		{3600, 50000},
		{1800, 200000},
		{900, 2000000},
		{60, 5000000},
	})
	n.raftNode.SetCodec(&raft.GOGOPBCodec{})
	n.raftNode.SetGzipSnapshot(true)
	meta, err := json.Marshal(address{
		HTTP: fmt.Sprintf("%s:%d", host, httpPort),
		RPC:  fmt.Sprintf("%s:%d", host, rpcPort),
	})
	if err != nil {
		logger.Fatal(err)
	}
	n.raftNode.SetNodeMeta(n.raftNode.Address(), meta)
	n.raftNode.LeaderChange(func() {
		n.resetLeader()
	})
	return n
}

func (n *Node) ListenAndServe() error {
	n.raftNode.Start()
	logger.Infoln("HTTP listening at:", fmt.Sprintf("http://%s:%d", n.host, n.httpPort))
	logger.Infoln("RPC listening at:", fmt.Sprintf("%s:%d", n.host, n.rpcPort))
	service := new(Service)
	service.node = n
	n.rpcServer = rpc.NewServer()
	n.rpcServer.RegisterName("S", service)
	n.rpcServer.SetLogLevel(rpc.OffLogLevel)
	if n.rpcTransport == nil {
		n.InitRPCProxy(MaxConnsPerHost, MaxIdleConnsPerHost)
	}
	go func() {
		fmt.Println("raftdb.node.rpc :", n.rpcServer.Listen("tcp", fmt.Sprintf(":%d", n.rpcPort), codec))
	}()
	n.rum.Recovery(rum.Recovery)
	n.rum.HandleFunc("/db/:key", n.leaderHandle(n.getHandler)).GET()
	n.rum.HandleFunc("/db/:key", n.leaderHandle(n.setHandler)).POST()
	return n.rum.Run(fmt.Sprintf(":%d", n.httpPort))
}

func (n *Node) InitRPCProxy(MaxConnsPerHost int, MaxIdleConnsPerHost int) {
	n.rpcTransport = &rpc.Transport{
		MaxConnsPerHost:     MaxConnsPerHost,
		MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		Options:             &rpc.Options{Network: network, Codec: codec},
	}
}

func (n *Node) resetLeader() {
	leader := n.raftNode.Leader()
	if leader != "" {
		meta, ok := n.raftNode.GetNodeMeta(leader)
		if ok && len(meta) > 0 {
			var addr address
			err := json.Unmarshal(meta, &addr)
			if err == nil {
				n.mu.Lock()
				n.leader = &addr
				n.mu.Unlock()
			}
		}
	}
}

func (n *Node) leaderRPCAddress() (addr string) {
	n.mu.Lock()
	if n.leader != nil {
		addr = n.leader.RPC
	}
	n.mu.Unlock()
	return
}

func (n *Node) leaderHTTPAddress() (addr string) {
	n.mu.Lock()
	if n.leader != nil {
		addr = n.leader.HTTP
	}
	n.mu.Unlock()
	return
}

func (n *Node) leaderHandle(hander http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if n.raftNode.IsLeader() {
			hander(w, r)
		} else {
			leaderHTTPAddress := n.leaderHTTPAddress()
			if leaderHTTPAddress != "" {
				proxy.Proxy(w, r, "http://"+leaderHTTPAddress+r.URL.Path)
			}
		}
	}
}

func (n *Node) setHandler(w http.ResponseWriter, req *http.Request) {
	params := n.rum.Params(req)
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
	params := n.rum.Params(req)
	if ok := n.raftNode.ReadIndex(); ok {
		value := n.db.Get(params["key"])
		w.Write([]byte(value))
	}
}
