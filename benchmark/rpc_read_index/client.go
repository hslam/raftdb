package main

import (
	"flag"
	"fmt"
	"github.com/hslam/raftdb/node"
	"github.com/hslam/rpc"
	"github.com/hslam/stats"
	"log"
	"strings"
)

var network string
var addrs string
var codec string
var clients int
var total int
var parallel int
var bar bool

func init() {
	flag.StringVar(&network, "network", "tcp", "-network=tcp")
	flag.StringVar(&addrs, "addrs", ":8001,:8002,:8003", "-addr=:8001,:8002,:8003")
	flag.StringVar(&codec, "codec", "pb", "-codec=code")
	flag.IntVar(&total, "total", 1000000, "-total=100000")
	flag.IntVar(&parallel, "parallel", 512, "-parallel=1")
	flag.IntVar(&clients, "clients", 1, "-clients=1")
	flag.BoolVar(&bar, "bar", true, "-bar=true")
	flag.Parse()
	stats.SetBar(bar)
	fmt.Printf("./client -network=%s -addrs=%s -codec=%s -total=%d -parallel=%d -clients=%d\n", network, addrs, codec, total, parallel, clients)
}

func main() {
	if clients < 1 || parallel < 1 || total < 1 {
		return
	}
	var wrkClients []stats.Client
	opts := &rpc.Options{Network: network, Codec: codec}
	peers := strings.Split(addrs, ",")
	for i := 0; i < clients; i++ {
		if conn, err := rpc.Dials(opts, peers...); err != nil {
			log.Fatalln("dailing error: ", err)
		} else {
			conn.Scheduling = rpc.LeastTimeScheduling
			req := &node.Request{Key: "foo", Value: "bar"}
			var res node.Response
			conn.Call("S.Set", req, &res)
			wrkClients = append(wrkClients, &WrkClient{conn})
		}
	}
	stats.StartPrint(parallel, total, wrkClients)
}

type WrkClient struct {
	*rpc.ReverseProxy
}

func (c *WrkClient) Call() (int64, int64, bool) {
	A := "foo"
	req := &node.Request{Key: A}
	var res node.Response
	c.ReverseProxy.Call("S.ReadIndexGet", req, &res)
	if res.Ok == true {
		return int64(len(A)), int64(len(res.Result)), true
	}
	return int64(len(A)), int64(len(res.Result)), false
}
