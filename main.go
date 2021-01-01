package main

import (
	"flag"
	"github.com/hslam/raftdb/node"
	"log"
	"strings"
)

var (
	host     string
	httpPort int
	rpcPort  int
	raftPort int
	peers    string
	join     string
	path     string
)

func init() {
	flag.StringVar(&host, "h", "localhost", "hostname")
	flag.IntVar(&httpPort, "p", 7001, "port")
	flag.IntVar(&rpcPort, "c", 8001, "port")
	flag.IntVar(&raftPort, "f", 9001, "port")
	flag.StringVar(&peers, "peers", "", "host:port,host:port")
	flag.StringVar(&join, "join", "", "host:port")
	flag.StringVar(&path, "path", "raft.1", "path")
	flag.Parse()
}

func main() {
	n := node.NewNode(path, host, httpPort, rpcPort, raftPort, strings.Split(peers, ","), join)
	log.Fatal(n.ListenAndServe())
}
