package main

import (
	"flag"
	"github.com/hslam/raftdb/node"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"strings"
)

var (
	host      string
	port      int
	rpcPort   int
	raftPort  int
	debug     bool
	debugPort int
	addrs     string
	join      string
	dataDir   string
	max       int
)

func init() {
	flag.StringVar(&host, "h", "localhost", "hostname")
	flag.IntVar(&port, "p", 7001, "port")
	flag.IntVar(&rpcPort, "c", 8001, "port")
	flag.IntVar(&raftPort, "f", 9001, "port")
	flag.StringVar(&addrs, "peers", "", "host:port,host:port")
	flag.BoolVar(&debug, "debug", true, "debug: -debug=false")
	flag.IntVar(&debugPort, "d", 6061, "debug_port: -dp=6060")
	flag.StringVar(&join, "join", "", "host:port")
	flag.StringVar(&dataDir, "path", "raft.1", "path")
	flag.IntVar(&max, "m", 8, "MaxConnsPerHost: -m=8")
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.Parse()
	go http.ListenAndServe(":"+strconv.Itoa(debugPort), nil)
	var peers []string
	if addrs != "" {
		peers = strings.Split(addrs, ",")
	}
	s := node.NewNode(dataDir, host, port, rpcPort, raftPort, peers, join)
	s.InitRPCProxy(max, 0)
	log.Fatal(s.ListenAndServe())
}
