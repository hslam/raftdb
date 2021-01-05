// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package node

import (
	"github.com/hslam/rpc"
	"sync"
)

type Client struct {
	lock   sync.Mutex
	leader string
	client *rpc.Client
}

func NewClient(targets ...string) *Client {
	client := &Client{}
	opts := &rpc.Options{Network: network, Codec: codec}
	client.client = rpc.NewClient(opts, targets...)
	client.client.Scheduling = rpc.LeastTimeScheduling
	client.client.Director = client.director
	return client
}

func (c *Client) director() (target string) {
	c.lock.Lock()
	target = c.leader
	c.lock.Unlock()
	return
}

func (c *Client) setDirector(target string) {
	c.lock.Lock()
	c.leader = target
	c.lock.Unlock()
}

func (c *Client) Set(key, value string) bool {
	req := &Request{Key: key, Value: value}
	var res Response
	c.client.Call("S.Set", req, &res)
	if len(res.Leader) > 0 {
		c.setDirector(res.Leader)
	}
	return res.Ok
}

func (c *Client) ReadIndexGet(key string) (value string, ok bool) {
	req := &Request{Key: key}
	var res Response
	c.client.Call("S.RGet", req, &res)
	if len(res.Leader) > 0 {
		c.setDirector(res.Leader)
	}
	return res.Result, res.Ok
}

func (c *Client) LeaseReadGet(key string) (value string, ok bool) {
	req := &Request{Key: key}
	var res Response
	c.client.Call("S.LGet", req, &res)
	if len(res.Leader) > 0 {
		c.setDirector(res.Leader)
	}
	return res.Result, res.Ok
}
