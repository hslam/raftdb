package main

import (
	"flag"
	"fmt"
	"github.com/hslam/raftdb/node"
	"github.com/hslam/stats"
	"math/rand"
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
	flag.IntVar(&parallel, "parallel", 2048, "-parallel=1")
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
		wrkClients = append(wrkClients, &WrkClient{node.NewClient(peers...)})
	}
	stats.StartPrint(parallel, total, wrkClients)
}

type WrkClient struct {
	*node.Client
}

func (c *WrkClient) Call() (int64, int64, bool) {
	A := RandString(4)
	B := RandString(32)
	if c.Client.Set(A, B) {
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
