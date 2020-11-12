package main

import (
	"bytes"
	"flag"
	"github.com/hslam/stats"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
)

var addr string
var clients int
var total_calls int
var parallel int
var bar bool

func init() {
	flag.StringVar(&addr, "addr", ":7001", "-addr=:7001")

	flag.IntVar(&clients, "clients", 200, "num: -clients=1")
	flag.IntVar(&total_calls, "total", 100000, "total_calls: -total=10000")
	flag.IntVar(&parallel, "parallel", 1, "total_calls: -total=10000")
	flag.BoolVar(&bar, "bar", false, "bar: -bar=true")
	flag.Parse()
	stats.SetBar(bar)
}
func main() {
	var wrkClients = make([]stats.Client, clients)
	for i := 0; i < clients; i++ {
		var conn = &WrkClient{}
		conn.client = &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: false,
				MaxConnsPerHost:   1,
			},
		}
		conn.url = "http://" + addr + "/db/"
		conn.meth = "POST"
		wrkClients[i] = conn
	}
	stats.StartPrint(parallel, total_calls, wrkClients)
}

type WrkClient struct {
	client *http.Client
	url    string
	meth   string
}

func (c *WrkClient) Call() (int64, int64, bool) {
	key := RandString(4)
	value := RandString(32)
	var requestBody = []byte(value)
	var requestBodyReader io.Reader
	if requestBody != nil {
		requestBodyReader = bytes.NewReader(requestBody)
	}
	req, _ := http.NewRequest(c.meth, c.url+key, requestBodyReader)
	resp, err := c.client.Do(req)
	if err != nil {
		return int64(len(key) + len(value)), 0, false
	}
	Body, err := ioutil.ReadAll(resp.Body)
	length := len(Body)
	if err != nil {
		return int64(len(key) + len(value)), 0, false
	}
	return int64(len(key) + len(value)), int64(length), true
}

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
