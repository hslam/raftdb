package main

import (
	"flag"
	"fmt"
	"github.com/hslam/raftdb/node"
	"github.com/hslam/rpc"
	"github.com/hslam/stats"
	"log"
	"math/rand"
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
	flag.IntVar(&parallel, "parallel", 1024, "-parallel=1")
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
			conn.Ping()
			wrkClients = append(wrkClients, &WrkClient{conn})
		}
	}
	stats.StartPrint(parallel, total, wrkClients)
}

type WrkClient struct {
	*rpc.ReverseProxy
}

func (c *WrkClient) Call() (int64, int64, bool) {
	A := RandString(4)
	B := RandString(32)
	req := &node.Request{Key: A, Value: B}
	var res node.Response
	c.ReverseProxy.Call("S.Set", req, &res)
	if res.Ok == true {
		return int64(len(A) + len(B)), 0, true
	}
	return int64(len(A) + len(B)), 0, false
}

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
