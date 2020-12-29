package node

import (
	"github.com/hslam/raft"
)

type Service struct {
	node *Node
}

func (s *Service) Set(req *Request, res *Response) error {
	if s.node.raftNode.IsLeader() {
		setCommand := newSetCommand(req.Key, req.Value)
		_, err := s.node.raftNode.Do(setCommand)
		setCommandPool.Put(setCommand)
		if err == nil {
			res.Ok = true
			return nil
		}
		return err
	} else {
		leaderRPCAddress := s.node.leaderRPCAddress()
		if leaderRPCAddress != "" {
			return s.node.rpcTransport.Call(leaderRPCAddress, "S.Set", req, res)
		}
		return raft.ErrNotLeader
	}
}

func (s *Service) Get(req *Request, res *Response) error {
	if s.node.raftNode.IsLeader() {
		if ok := s.node.raftNode.LeaseRead(); ok {
			value := s.node.db.Get(req.Key)
			res.Result = []byte(value)
			res.Ok = true
			return nil
		}
		return nil
	} else {
		leaderRPCAddress := s.node.leaderRPCAddress()
		if leaderRPCAddress != "" {
			return s.node.rpcTransport.Call(leaderRPCAddress, "S.Get", req, res)
		}
		return raft.ErrNotLeader
	}
}

func (s *Service) ReadIndexGet(req *Request, res *Response) error {
	if s.node.raftNode.IsLeader() {
		if ok := s.node.raftNode.ReadIndex(); ok {
			value := s.node.db.Get(req.Key)
			res.Result = []byte(value)
			res.Ok = true
			return nil
		}
		return nil
	} else {
		leaderRPCAddress := s.node.leaderRPCAddress()
		if leaderRPCAddress != "" {
			return s.node.rpcTransport.Call(leaderRPCAddress, "S.ReadIndexGet", req, res)
		}
		return raft.ErrNotLeader
	}
}
