// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package node

import (
	"github.com/hslam/raft"
)

func newSetCommand(key string, value string) raft.Command {
	c := setCommandPool.Get().(*SetCommand)
	c.Key = key
	c.Value = value
	return c
}

func (c *SetCommand) Type() int32 {
	return 1
}

func (c *SetCommand) Do(context interface{}) (interface{}, error) {
	db := context.(*DB)
	db.Set(c.Key, c.Value)
	return nil, nil
}
