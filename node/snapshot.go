// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package node

import (
	"github.com/hslam/raft"
	"io"
	"io/ioutil"
)

type Snapshot struct {
	db *DB
}

func NewSnapshot(db *DB) raft.Snapshot {
	return &Snapshot{db: db}
}

func (s *Snapshot) Save(w io.Writer) (int, error) {
	raw, err := s.db.Data()
	if err != nil {
		return 0, err
	}
	return w.Write(raw)
}

func (s *Snapshot) Recover(r io.Reader) (int, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return len(raw), err
	}
	err = s.db.SetData(raw)
	return len(raw), err
}
