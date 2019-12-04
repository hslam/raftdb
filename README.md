# raftdb

[raftdb](https://hslam.com/git/x/raftdb  "raftdb") is an example usage of [raft](https://hslam.com/git/x/raft  "raft") library.

## Get started

### Install
```
go get hslam.com/git/x/raftdb
```

### Build
```
go build -tags=use_cgo main.go
```

### Singleton
```sh
./raftdb -h=localhost -p=7001 -c=8001 -f=9001 -d=6061 -m=8 -peers="" -path=./raftdb.1
```
### Three nodes
```sh
./raftdb -h=localhost -p=7001 -c=8001 -f=9001 -d=6061 -m=8 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./raftdb.1
./raftdb -h=localhost -p=7002 -c=8002 -f=9002 -d=6062 -m=8 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./raftdb.2
./raftdb -h=localhost -p=7003 -c=8003 -f=9003 -d=6063 -m=8 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./raftdb.3
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

#### Linux Environment
* **CPU** 12 Cores 3.1 GHz
* **Memory** 24 GiB

```
cluster    operation transport requests/s average(ms) fastest(ms) median(ms) p99(ms) slowest(ms)
Singleton  ReadIndex HTTP      74456      6.63        2.62        6.23       12.12   110.90
Singleton  ReadIndex RPC       293865     13.14       4.09        12.14      31.22   35.09
Singleton  Write     HTTP      57488      8.79        2.19        7.68       24.00   119.71
Singleton  Write     RPC       132045     30.21       6.39        27.59      70.14   86.11
ThreeNodes ReadIndex HTTP      43053      11.72       2.83        7.43       58.92   1125.58
ThreeNodes ReadIndex RPC       267685     14.65       4.36        13.47      31.59   44.72
ThreeNodes Write     HTTP      35241      14.42       4.21        10.41      73.28   114.84
ThreeNodes Write     RPC       103035     38.82       8.91        38.90      76.74    88.05
```

vim benchmark.sh
```sh
#!/bin/sh

nohup ./raftdb -h=localhost -p=7001 -c=8001 -f=9001 -d=6061 -m=32 -peers="" -path=./tmp/default.raftdb.1  >> ./tmp/default.out.1.log 2>&1 &
sleep 30s
curl -XPOST http://localhost:7001/db/foo -d 'bar'
sleep 3s
curl http://localhost:7001/db/foo
sleep 3s
./http_read_index -p=7001 -parallel=1 -total=1000000 -clients=512 -bar=false
sleep 10s
./rpc_read_index -network=tcp -codec=pb -compress=gzip -h=127.0.0.1 -p=8001 -parallel=512 -total=1000000 -multiplexing=true -batch=true -batch_async=true -clients=8 -bar=false
sleep 10s
./http_write -p=7001 -parallel=1 -total=1000000 -clients=512 -bar=false
sleep 10s
./rpc_write -network=tcp -codec=pb -compress=gzip -h=127.0.0.1 -p=8001 -parallel=512 -total=1000000 -multiplexing=true -batch=true -batch_async=true -clients=8 -bar=false
sleep 10s
killall raftdb
sleep 3s
nohup ./raftdb -h=localhost -p=7001 -c=8001 -f=9001 -d=6061 -m=32 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./tmp/raftdb.1  >> ./tmp/out.1.log 2>&1 &
sleep 3s
nohup ./raftdb -h=localhost -p=7002 -c=8002 -f=9002 -d=6062 -m=32 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./tmp/raftdb.2  >> ./tmp/out.2.log 2>&1 &
nohup ./raftdb -h=localhost -p=7003 -c=8003 -f=9003 -d=6063 -m=32 -peers=localhost:9001,localhost:9002,localhost:9003 -path=./tmp/raftdb.3  >> ./tmp/out.3.log 2>&1 &
sleep 30s
curl -XPOST http://localhost:7001/db/foo -d 'bar'
sleep 3s
curl http://localhost:7001/db/foo
sleep 3s
./http_read_index -p=7001 -parallel=1 -total=1000000 -clients=512 -bar=false
sleep 10s
./rpc_read_index -network=tcp -codec=pb -compress=gzip -h=127.0.0.1 -p=8001 -parallel=512 -total=1000000 -multiplexing=true -batch=true -batch_async=true -clients=8 -bar=false
sleep 10s
./http_write -p=7001 -parallel=1 -total=1000000 -clients=512 -bar=false
sleep 10s
./rpc_write -network=tcp -codec=pb -compress=gzip -h=127.0.0.1 -p=8001 -parallel=512 -total=1000000 -multiplexing=true -batch=true -batch_async=true -clients=8 -bar=false
sleep 10s
killall raftdb
```

#### HTTP READINDEX SINGLETON BENCHMARK
```
Summary:
	Clients:	512
	Parallel calls per client:	1
	Total calls:	100000
	Total time:	1.34s
	Requests per second:	74456.78
	Fastest time for request:	2.62ms
	Average time per request:	6.63ms
	Slowest time for request:	110.90ms

Time:
	0.1%	time for request:	3.53ms
	1%	time for request:	4.05ms
	5%	time for request:	4.44ms
	10%	time for request:	4.69ms
	25%	time for request:	5.27ms
	50%	time for request:	6.23ms
	75%	time for request:	7.13ms
	90%	time for request:	8.12ms
	95%	time for request:	9.27ms
	99%	time for request:	12.12ms
	99.9%	time for request:	81.47ms

Result:
	Response ok:	100000 (100.00%)
	Errors:	0 (0.00%)
```
#### RPC READINDEX SINGLETON BENCHMARK
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	100000
	Total time:	0.34s
	Requests per second:	293865.27
	Fastest time for request:	4.09ms
	Average time per request:	13.14ms
	Slowest time for request:	35.09ms

Time:
	0.1%	time for request:	5.12ms
	1%	time for request:	5.63ms
	5%	time for request:	6.82ms
	10%	time for request:	7.24ms
	25%	time for request:	8.30ms
	50%	time for request:	12.14ms
	75%	time for request:	16.22ms
	90%	time for request:	20.89ms
	95%	time for request:	26.36ms
	99%	time for request:	31.22ms
	99.9%	time for request:	34.33ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
#### HTTP WRITE SINGLETON BENCHMARK
```
Summary:
	Clients:	512
	Parallel calls per client:	1
	Total calls:	100000
	Total time:	1.74s
	Requests per second:	57488.54
	Fastest time for request:	2.19ms
	Average time per request:	8.79ms
	Slowest time for request:	119.71ms

Time:
	0.1%	time for request:	3.31ms
	1%	time for request:	3.70ms
	5%	time for request:	4.39ms
	10%	time for request:	5.05ms
	25%	time for request:	6.16ms
	50%	time for request:	7.68ms
	75%	time for request:	9.84ms
	90%	time for request:	12.58ms
	95%	time for request:	15.76ms
	99%	time for request:	24.00ms
	99.9%	time for request:	85.99ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
#### RPC WRITE SINGLETON BENCHMARK
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	100000
	Total time:	0.76s
	Requests per second:	132045.80
	Fastest time for request:	6.39ms
	Average time per request:	30.21ms
	Slowest time for request:	86.11ms

Time:
	0.1%	time for request:	8.53ms
	1%	time for request:	11.02ms
	5%	time for request:	14.28ms
	10%	time for request:	16.39ms
	25%	time for request:	20.35ms
	50%	time for request:	27.59ms
	75%	time for request:	36.17ms
	90%	time for request:	49.05ms
	95%	time for request:	57.28ms
	99%	time for request:	70.14ms
	99.9%	time for request:	76.14ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
#### HTTP READINDEX THREE NODES BENCHMARK
```
Summary:
    Clients:	512
	Parallel calls per client:	1
	Total calls:	100000
	Total time:	2.32s
	Requests per second:	43053.90
	Fastest time for request:	2.83ms
	Average time per request:	11.72ms
	Slowest time for request:	1125.58ms

Time:
	0.1%	time for request:	3.54ms
	1%	time for request:	4.09ms
	5%	time for request:	4.76ms
	10%	time for request:	5.24ms
	25%	time for request:	6.28ms
	50%	time for request:	7.43ms
	75%	time for request:	8.65ms
	90%	time for request:	12.77ms
	95%	time for request:	39.88ms
	99%	time for request:	58.92ms
	99.9%	time for request:	1068.66ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

#### RPC READINDEX THREE NODES BENCHMARK
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	100000
	Total time:	0.37s
	Requests per second:	267685.30
	Fastest time for request:	4.36ms
	Average time per request:	14.65ms
	Slowest time for request:	44.72ms

Time:
	0.1%	time for request:	5.34ms
	1%	time for request:	6.50ms
	5%	time for request:	7.75ms
	10%	time for request:	8.26ms
	25%	time for request:	9.69ms
	50%	time for request:	13.47ms
	75%	time for request:	18.37ms
	90%	time for request:	22.27ms
	95%	time for request:	25.10ms
	99%	time for request:	31.59ms
	99.9%	time for request:	44.32ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

#### HTTP WRITE THREE NODES BENCHMARK
```
Summary:
	Clients:	512
	Parallel calls per client:	1
	Total calls:	100000
	Total time:	2.84s
	Requests per second:	35241.64
	Fastest time for request:	4.21ms
	Average time per request:	14.42ms
	Slowest time for request:	114.84ms

Time:
	0.1%	time for request:	5.59ms
	1%	time for request:	6.48ms
	5%	time for request:	7.32ms
	10%	time for request:	7.84ms
	25%	time for request:	8.85ms
	50%	time for request:	10.41ms
	75%	time for request:	13.18ms
	90%	time for request:	24.42ms
	95%	time for request:	42.78ms
	99%	time for request:	73.28ms
	99.9%	time for request:	108.08ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

#### RPC WRITE THREE NODES BENCHMARK
```
Summary:
	Clients:	8
	Parallel calls per client:	512
	Total calls:	100000
	Total time:	0.97s
	Requests per second:	103035.21
	Fastest time for request:	8.91ms
	Average time per request:	38.82ms
	Slowest time for request:	88.05ms

Time:
	0.1%	time for request:	9.53ms
	1%	time for request:	14.58ms
	5%	time for request:	18.81ms
	10%	time for request:	21.52ms
	25%	time for request:	28.00ms
	50%	time for request:	38.90ms
	75%	time for request:	46.57ms
	90%	time for request:	54.88ms
	95%	time for request:	62.80ms
	99%	time for request:	76.74ms
	99.9%	time for request:	85.04ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

#### ETCD WRITE THREE NODES BENCHMARK
```
Summary:
	Conns:	8
	Clients:	512
	Total calls:	100000
	Total time:	2.54s
	Requests per second:	39357.45
	Fastest time for request:	0.90ms
	Average time per request:	12.90ms
	Slowest time for request:	71.50ms

Time:
	10%	time for request:	7.10ms
	50%	time for request:	10.50ms
	90%	time for request:	17.50ms
	99%	time for request:	52.10ms

Result:
	Response ok:	100000 (100.00%)
	Errors:	0 (0.00%)
```

## Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)

## Authors
raft was written by Mort Huang.