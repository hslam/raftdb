package node

import (
	"github.com/hslam/raft"
	"net/url"
	"strconv"
)

type Service struct {
	node *Node
}

func (s *Service) Set(req *Request, res *Response) error {
	if s.node.raft_node.IsLeader() {
		setCommand := newSetCommand(req.Key, req.Value)
		_, err := s.node.raft_node.Do(setCommand)
		setCommandPool.Put(setCommand)
		if err == nil {
			res.Ok = true
			return nil
		}
		return err
	} else {
		leader := s.node.raft_node.Leader()
		if leader != "" {
			leader_url, err := url.Parse("http://" + leader)
			if err != nil {
				panic(err)
			}
			port, err := strconv.Atoi(leader_url.Port())
			if err != nil {
				panic(err)
			}
			leader_host := leader_url.Hostname() + ":" + strconv.Itoa(port-1000)
			return s.node.rpc_transport.Call(leader_host, "S.Set", req, res)
		}
		return raft.ErrNotLeader
	}
}

func (s *Service) Get(req *Request, res *Response) error {
	if s.node.raft_node.IsLeader() {
		if ok := s.node.raft_node.LeaseRead(); ok {
			value := s.node.db.Get(req.Key)
			res.Result = []byte(value)
			res.Ok = true
			return nil
		}
		return nil
	} else {
		leader := s.node.raft_node.Leader()
		if leader != "" {
			leader_url, err := url.Parse("http://" + leader)
			if err != nil {
				panic(err)
			}
			port, err := strconv.Atoi(leader_url.Port())
			if err != nil {
				panic(err)
			}
			leader_host := leader_url.Hostname() + ":" + strconv.Itoa(port-1000)
			return s.node.rpc_transport.Call(leader_host, "S.Get", req, res)
		}
		return raft.ErrNotLeader
	}
}

func (s *Service) ReadIndexGet(req *Request, res *Response) error {
	if s.node.raft_node.IsLeader() {
		if ok := s.node.raft_node.ReadIndex(); ok {
			value := s.node.db.Get(req.Key)
			res.Result = []byte(value)
			res.Ok = true
			return nil
		}
		return nil
	} else {
		leader := s.node.raft_node.Leader()
		if leader != "" {
			leader_url, err := url.Parse("http://" + leader)
			if err != nil {
				panic(err)
			}
			port, err := strconv.Atoi(leader_url.Port())
			if err != nil {
				panic(err)
			}
			leader_host := leader_url.Hostname() + ":" + strconv.Itoa(port-1000)
			return s.node.rpc_transport.Call(leader_host, "S.ReadIndexGet", req, res)
		}
		return raft.ErrNotLeader
	}
}
