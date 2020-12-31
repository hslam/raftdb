// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

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
			err := s.node.rpcTransport.Call(leaderRPCAddress, "S.Set", req, res)
			res.Leader = leaderRPCAddress
			return err
		}
		return raft.ErrNotLeader
	}
}

func (s *Service) LGet(req *Request, res *Response) error {
	if s.node.raftNode.IsLeader() {
		if ok := s.node.raftNode.LeaseRead(); ok {
			value := s.node.db.Get(req.Key)
			res.Result = value
			res.Ok = true
			return nil
		}
		return nil
	} else {
		leaderRPCAddress := s.node.leaderRPCAddress()
		if leaderRPCAddress != "" {
			err := s.node.rpcTransport.Call(leaderRPCAddress, "S.LGet", req, res)
			res.Leader = leaderRPCAddress
			return err
		}
		return raft.ErrNotLeader
	}
}

func (s *Service) RGet(req *Request, res *Response) error {
	if s.node.raftNode.IsLeader() {
		if ok := s.node.raftNode.ReadIndex(); ok {
			value := s.node.db.Get(req.Key)
			res.Result = value
			res.Ok = true
			return nil
		}
		return nil
	} else {
		leaderRPCAddress := s.node.leaderRPCAddress()
		if leaderRPCAddress != "" {
			err := s.node.rpcTransport.Call(leaderRPCAddress, "S.RGet", req, res)
			res.Leader = leaderRPCAddress
			return err
		}
		return raft.ErrNotLeader
	}
}
