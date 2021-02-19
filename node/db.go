// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package node

import (
	"encoding/json"
	"sync"
)

type DB struct {
	mutex sync.RWMutex
	data  map[string]string
}

func newDB() *DB {
	return &DB{
		data: make(map[string]string),
	}
}

func (db *DB) Data() (raw []byte, err error) {
	db.mutex.RLock()
	raw, err = json.Marshal(db.data)
	db.mutex.RUnlock()
	return
}

func (db *DB) SetData(raw []byte) (err error) {
	var data map[string]string
	err = json.Unmarshal(raw, &data)
	if err == nil {
		db.mutex.Lock()
		db.data = data
		db.mutex.Unlock()
	}
	return
}

func (db *DB) Set(key string, value string) {
	db.mutex.Lock()
	db.data[key] = value
	db.mutex.Unlock()
}

func (db *DB) Get(key string) (value string) {
	db.mutex.RLock()
	value = db.data[key]
	db.mutex.RUnlock()
	return
}
