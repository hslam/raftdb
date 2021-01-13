// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package node

import (
	"github.com/hslam/atomic"
	"github.com/hslam/rpc"
	"time"
)

type Client struct {
	leader   *atomic.String
	ready    *atomic.Bool
	client   *rpc.Client
	Fallback time.Duration
}

func NewClient(targets ...string) *Client {
	client := &Client{
		leader:   atomic.NewString(""),
		ready:    atomic.NewBool(true),
		Fallback: time.Second * 3,
	}
	opts := &rpc.Options{Network: network, Codec: codec}
	client.client = rpc.NewClient(opts, targets...)
	client.client.Scheduling = rpc.LeastTimeScheduling
	client.client.Director = client.director
	return client
}

func (c *Client) director() (target string) {
	target = c.leader.Load()
	return
}

func (c *Client) setDirector(target string) {
	c.leader.Store(target)
	if len(target) > 0 {
		c.ready.Store(true)
	} else {
		c.ready.Store(false)
		if c.Fallback > 0 {
			c.client.Fallback(c.Fallback)
		}
	}
}

func (c *Client) Set(key, value string) bool {
	if !c.ready.Load() {
		time.Sleep(time.Millisecond * 200)
	}
	req := &Request{Key: key, Value: value}
	var res Response
	err := c.client.Call("S.Set", req, &res)
	if len(res.Leader) > 0 {
		c.setDirector(res.Leader)
	} else if err == rpc.ErrShutdown || !res.Ok {
		c.setDirector("")
	}
	return res.Ok
}

func (c *Client) ReadIndexGet(key string) (value string, ok bool) {
	if !c.ready.Load() {
		time.Sleep(time.Millisecond * 200)
	}
	req := &Request{Key: key}
	var res Response
	err := c.client.Call("S.RGet", req, &res)
	if len(res.Leader) > 0 {
		c.setDirector(res.Leader)
	} else if err == rpc.ErrShutdown || !res.Ok {
		c.setDirector("")
	}
	return res.Result, res.Ok
}

func (c *Client) LeaseReadGet(key string) (value string, ok bool) {
	if !c.ready.Load() {
		time.Sleep(time.Millisecond * 200)
	}
	req := &Request{Key: key}
	var res Response
	err := c.client.Call("S.LGet", req, &res)
	if len(res.Leader) > 0 {
		c.setDirector(res.Leader)
	} else if err == rpc.ErrShutdown || !res.Ok {
		c.setDirector("")
	}
	return res.Result, res.Ok
}
