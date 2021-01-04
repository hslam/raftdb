package main

import (
	"flag"
	"github.com/hslam/raftdb/node"
	"log"
	"strings"
)

var host, peers, path string
var httpPort, rpcPort, raftPort int
var join bool

func init() {
	flag.StringVar(&host, "h", "localhost", "")
	flag.IntVar(&httpPort, "p", 7001, "")
	flag.IntVar(&rpcPort, "c", 8001, "")
	flag.IntVar(&raftPort, "f", 9001, "")
	flag.StringVar(&peers, "peers", "", "")
	flag.StringVar(&path, "path", "raftdb.1", "")
	flag.BoolVar(&join, "join", false, "")
	flag.Parse()
}

func main() {
	nodes := strings.Split(peers, ",")
	n := node.NewNode(path, host, httpPort, rpcPort, raftPort, nodes, join)
	log.Fatal(n.ListenAndServe())
}
