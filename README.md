# raftdb
The raftdb implements a simple distributed key-value datastore, using the [raft](https://github.com/hslam/raft  "raft") distributed consensus protocol.

## Get started

### Install
```
go get github.com/hslam/raftdb
```

### Build
```
go build -o raftdb main.go
```

#### Three nodes
```sh
./raftdb -h=localhost -p=7001 -c=8001 -f=9001 -path=./raftdb.1 \
-peers=localhost:9001,localhost:9002,localhost:9003

./raftdb -h=localhost -p=7002 -c=8002 -f=9002  -path=./raftdb.2 \
-peers=localhost:9001,localhost:9002,localhost:9003

./raftdb -h=localhost -p=7003 -c=8003 -f=9003  -path=./raftdb.3 \
-peers=localhost:9001,localhost:9002,localhost:9003
```

##### HTTP SET
```
curl -XPOST http://localhost:7001/db/foo -d 'bar'
```

##### HTTP GET
```
curl http://localhost:7001/db/foo
```

##### Client example
```go
package main

import (
	"fmt"
	"github.com/hslam/raftdb/node"
)

func main() {
	client := node.NewClient("localhost:8001", "localhost:8002", "localhost:8003")
	key := "foo"
	value := "Hello World"
	if ok := client.Set(key, value); !ok {
		panic("set failed")
	}
	if result, ok := client.LeaseReadGet(key); ok && result != value {
		panic(result)
	}
	if result, ok := client.ReadIndexGet(key); ok {
		fmt.Println(result)
	}
}
```

##### Output
```
Hello World
```

### [Benchmark](https://github.com/hslam/raft-benchmark  "raft-benchmark")
Running on a three nodes cluster.
##### Write

<img src="https://raw.githubusercontent.com/hslam/raft-benchmark/master/raft-write-qps.png" width = "400" height = "300" alt="write-qps" align=center><img src="https://raw.githubusercontent.com/hslam/raft-benchmark/master/raft-write-p99.png" width = "400" height = "300" alt="write-p99" align=center>

##### Read Index

<img src="https://raw.githubusercontent.com/hslam/raft-benchmark/master/raft-read-qps.png" width = "400" height = "300" alt="read-qps" align=center><img src="https://raw.githubusercontent.com/hslam/raft-benchmark/master/raft-read-p99.png" width = "400" height = "300" alt="read-p99" align=center>

## License
This package is licensed under a MIT license (Copyright (c) 2019 Meng Huang)

## Author
raftdb was written by Meng Huang.