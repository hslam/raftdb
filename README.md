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
pkg  cluster    operation transport requests/s average fastest median  p99       slowest
RAFT Singleton  ReadIndex HTTP      73348      6.77ms  2.59ms  6.29ms  14.23ms   116.32ms
RAFT Singleton  Write     HTTP      60671      8.27ms  2.47ms  7.34ms  23.08ms   134.83ms
RAFT ThreeNodes ReadIndex HTTP      47642      10.49ms 3.04ms  8.14ms  52.66ms   1048.29ms
RAFT ThreeNodes Write     HTTP      37647      13.39ms 4.62ms  9.61ms  77.57ms   142.90ms
RAFT Singleton  ReadIndex RPC       310222     12.51ms 4.36ms  12.07ms 29.79ms   34.28ms
RAFT Singleton  Write     RPC       138411     28.92ms 6.09ms  24.64ms 103.57ms  121.66ms
RAFT ThreeNodes ReadIndex RPC       285650     13.40ms 4.27ms  12.49ms 29.01ms   32.91ms
RAFT ThreeNodes Write     RPC       118325     33.74ms 9.76ms  33.40ms 71.38ms   81.32ms
ETCD Singleton  ReadIndex HTTP      22991      12.04ms -       3.42ms  42.38ms   1416.00ms
ETCD Singleton  Write     HTTP      23189      16.09ms -       4.16ms  1003.71ms 1212.00ms
ETCD ThreeNodes ReadIndex HTTP      29180      11.35ms -       3.14ms  29.41ms   1208.00ms
ETCD ThreeNodes Write     HTTP      15325      26.67ms -       7.31ms  1010.76ms 1218.00ms
ETCD Singleton  ReadIndex -         94665      4.90ms  0.10ms  4.80ms  7.70ms    17.90ms
ETCD Singleton  Write     -         65242      7.70ms  0.70ms  7.30ms  11.00ms   27.90ms
ETCD ThreeNodes ReadIndex -         92092      5.20ms  0.30ms  4.90ms  8.50ms    21.40ms
ETCD ThreeNodes Write     -         38790      13.10ms 1.00ms  10.30ms 18.20ms   71.10ms
```

#### HTTP READINDEX SINGLETON BENCHMARK
```
Summary:
	Clients:	512
	Parallel calls per client:	1
	Total calls:	100000
	Total time:	1.36s
	Requests per second:	73348.95
	Fastest time for request:	2.59ms
	Average time per request:	6.77ms
	Slowest time for request:	116.32ms

Time:
	0.1%	time for request:	3.40ms
	1%	time for request:	4.14ms
	5%	time for request:	4.59ms
	10%	time for request:	4.87ms
	25%	time for request:	5.45ms
	50%	time for request:	6.29ms
	75%	time for request:	7.12ms
	90%	time for request:	8.19ms
	95%	time for request:	9.54ms
	99%	time for request:	14.23ms
	99.9%	time for request:	79.55ms

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
	Total time:	0.32s
	Requests per second:	310222.77
	Fastest time for request:	4.36ms
	Average time per request:	12.51ms
	Slowest time for request:	34.28ms

Time:
	0.1%	time for request:	4.71ms
	1%	time for request:	5.84ms
	5%	time for request:	6.90ms
	10%	time for request:	7.43ms
	25%	time for request:	8.30ms
	50%	time for request:	12.07ms
	75%	time for request:	15.54ms
	90%	time for request:	17.86ms
	95%	time for request:	20.70ms
	99%	time for request:	29.79ms
	99.9%	time for request:	34.10ms

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
	Total time:	1.65s
	Requests per second:	60671.95
	Fastest time for request:	2.47ms
	Average time per request:	8.27ms
	Slowest time for request:	134.83ms

Time:
	0.1%	time for request:	3.22ms
	1%	time for request:	3.71ms
	5%	time for request:	4.41ms
	10%	time for request:	4.95ms
	25%	time for request:	5.99ms
	50%	time for request:	7.34ms
	75%	time for request:	9.34ms
	90%	time for request:	11.29ms
	95%	time for request:	13.00ms
	99%	time for request:	23.08ms
	99.9%	time for request:	94.16ms

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
	Total time:	0.72s
	Requests per second:	138411.37
	Fastest time for request:	6.09ms
	Average time per request:	28.92ms
	Slowest time for request:	121.66ms

Time:
	0.1%	time for request:	8.47ms
	1%	time for request:	11.02ms
	5%	time for request:	13.24ms
	10%	time for request:	15.18ms
	25%	time for request:	18.75ms
	50%	time for request:	24.64ms
	75%	time for request:	32.47ms
	90%	time for request:	44.07ms
	95%	time for request:	56.01ms
	99%	time for request:	103.57ms
	99.9%	time for request:	114.01ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```

#### ETCD READINDEX SINGLETON BENCHMARK
```
Summary:
	Conns:	8
	Clients:	512
	Total calls:	100000
	Total time:	2.54s
	Requests per second:	94665.88
	Fastest time for request:	0.10ms
	Average time per request:	4.90ms
	Slowest time for request:	17.90ms

Time:
	10%	time for request:	2.20ms
	50%	time for request:	4.80ms
	90%	time for request:	7.70ms
	99%	time for request:	11.40ms

Result:
	Response ok:	100000 (100.00%)
	Errors:	0 (0.00%)
```
#### ETCD WRITE SINGLETON BENCHMARK
```
Summary:
	Conns:	8
	Clients:	512
	Total calls:	100000
	Total time:	2.54s
	Requests per second:	65242.66
	Fastest time for request:	0.70ms
	Average time per request:	7.70ms
	Slowest time for request:	27.90ms

Time:
	10%	time for request:	4.50ms
	50%	time for request:	7.30ms
	90%	time for request:	11.00ms
	99%	time for request:	15.80ms

Result:
	Response ok:	100000 (100.00%)
	Errors:	0 (0.00%)
```


#### HTTP READINDEX THREE NODES BENCHMARK
```
Summary:
	Clients:	512
	Parallel calls per client:	1
	Total calls:	100000
	Total time:	2.10s
	Requests per second:	47642.48
	Fastest time for request:	3.04ms
	Average time per request:	10.49ms
	Slowest time for request:	1048.29ms

Time:
	0.1%	time for request:	3.67ms
	1%	time for request:	4.30ms
	5%	time for request:	5.17ms
	10%	time for request:	5.78ms
	25%	time for request:	7.09ms
	50%	time for request:	8.14ms
	75%	time for request:	9.02ms
	90%	time for request:	11.69ms
	95%	time for request:	32.71ms
	99%	time for request:	52.66ms
	99.9%	time for request:	92.87ms

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
	Total time:	0.35s
	Requests per second:	285650.63
	Fastest time for request:	4.27ms
	Average time per request:	13.40ms
	Slowest time for request:	32.91ms

Time:
	0.1%	time for request:	5.08ms
	1%	time for request:	5.91ms
	5%	time for request:	7.07ms
	10%	time for request:	7.60ms
	25%	time for request:	8.87ms
	50%	time for request:	12.49ms
	75%	time for request:	17.17ms
	90%	time for request:	20.14ms
	95%	time for request:	22.96ms
	99%	time for request:	29.01ms
	99.9%	time for request:	32.58ms

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
	Total time:	2.66s
	Requests per second:	37647.50
	Fastest time for request:	4.62ms
	Average time per request:	13.39ms
	Slowest time for request:	142.90ms

Time:
	0.1%	time for request:	5.64ms
	1%	time for request:	6.26ms
	5%	time for request:	7.01ms
	10%	time for request:	7.47ms
	25%	time for request:	8.35ms
	50%	time for request:	9.61ms
	75%	time for request:	11.64ms
	90%	time for request:	18.48ms
	95%	time for request:	39.33ms
	99%	time for request:	77.57ms
	99.9%	time for request:	117.20ms

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
	Total time:	0.85s
	Requests per second:	118325.69
	Fastest time for request:	9.76ms
	Average time per request:	33.74ms
	Slowest time for request:	81.32ms

Time:
	0.1%	time for request:	11.21ms
	1%	time for request:	15.13ms
	5%	time for request:	17.24ms
	10%	time for request:	19.39ms
	25%	time for request:	24.52ms
	50%	time for request:	33.40ms
	75%	time for request:	39.64ms
	90%	time for request:	48.87ms
	95%	time for request:	56.72ms
	99%	time for request:	71.38ms
	99.9%	time for request:	80.04ms

Result:
	Response ok:	1000000 (100.00%)
	Errors:	0 (0.00%)
```
#### ETCD READINDEX THREE NODES BENCHMARK
```
Summary:
	Conns:	8
	Clients:	512
	Total calls:	100000
	Total time:	2.54s
	Requests per second:	92092.45
	Fastest time for request:	0.30ms
	Average time per request:	5.20ms
	Slowest time for request:	21.40ms

Time:
	10%	time for request:	2.50ms
	50%	time for request:	4.90ms
	90%	time for request:	8.50ms
	99%	time for request:	13.00ms

Result:
	Response ok:	100000 (100.00%)
	Errors:	0 (0.00%)
```
#### ETCD WRITE THREE NODES BENCHMARK
```
Summary:
	Conns:	8
	Clients:	512
	Total calls:	100000
	Total time:	2.54s
	Requests per second:	38790.33
	Fastest time for request:	1.00ms
	Average time per request:	13.10ms
	Slowest time for request:	71.10ms

Time:
	10%	time for request:	7.10ms
	50%	time for request:	10.30ms
	90%	time for request:	18.20ms
	99%	time for request:	56.40ms

Result:
	Response ok:	100000 (100.00%)
	Errors:	0 (0.00%)
```

## Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Mort Huang)

## Authors
raftdb was written by Mort Huang.