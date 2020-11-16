# raftdb
The raftdb implements a simple distributed key-value datastore, using the [raft](https://github.com/hslam/raft  "raft") distributed consensus protocol.

## Get started

### Install
```
go get github.com/hslam/raftdb
```

### Build
```
go build -tags=use_cgo main.go
```

### Singleton
```sh
./raftdb -h=localhost -p=7001 -c=8001 -f=9001 -d=6061 -m=1 -peers="" -path=./raftdb.1
```
### Three nodes
```sh
./raftdb -h=localhost -p=7001 -c=8001 -f=9001 -d=6061 -m=1 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./raftdb.1
./raftdb -h=localhost -p=7002 -c=8002 -f=9002 -d=6062 -m=1 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./raftdb.2
./raftdb -h=localhost -p=7003 -c=8003 -f=9003 -d=6063 -m=1 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./raftdb.3
```
**HTTP SET**
```
curl -XPOST http://localhost:7001/db/foo -d 'bar'
```
**HTTP GET**
```
curl http://localhost:7001/db/foo
```

## Benchmark
```
pkg     cluster     operation   transport   requests/s  p99
RAFT    ThreeNodes  ReadIndex   RPC         143726      7.21ms
RAFT    ThreeNodes  Write       RPC         33874       62.46ms
ETCD    ThreeNodes  ReadIndex   GRPC        31808       50.97ms
ETCD    ThreeNodes  Write       GRPC        9635        307.00ms
```

## License
This package is licenced under a MIT license (Copyright (c) 2019 Meng Huang)

## Author
raftdb was written by Meng Huang.