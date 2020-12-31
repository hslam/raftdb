package main

import (
	"flag"
	"fmt"
	"github.com/hslam/raftdb/node"
	"github.com/hslam/stats"
	"strings"
)

var addrs string
var clients int
var total int
var parallel int
var bar bool

func init() {
	flag.StringVar(&addrs, "addrs", ":8001,:8002,:8003", "-addr=:8001,:8002,:8003")
	flag.IntVar(&total, "total", 1000000, "-total=100000")
	flag.IntVar(&parallel, "parallel", 512, "-parallel=1")
	flag.IntVar(&clients, "clients", 1, "-clients=1")
	flag.BoolVar(&bar, "bar", true, "-bar=true")
	flag.Parse()
	stats.SetBar(bar)
	fmt.Printf("./client -addrs=%s -total=%d -parallel=%d -clients=%d\n", addrs, total, parallel, clients)
}

func main() {
	if clients < 1 || parallel < 1 || total < 1 {
		return
	}
	var wrkClients []stats.Client
	peers := strings.Split(addrs, ",")
	for i := 0; i < clients; i++ {
		client := node.NewClient(peers...)
		client.Set("foo", "bar")
		wrkClients = append(wrkClients, &WrkClient{client})
	}
	stats.StartPrint(parallel, total, wrkClients)
}

type WrkClient struct {
	*node.Client
}

func (c *WrkClient) Call() (int64, int64, bool) {
	A := "foo"
	result, ok := c.Client.ReadIndexGet(A)
	if ok {
		return int64(len(A)), int64(len(result)), true
	}
	return int64(len(A)), int64(len(result)), false
}
