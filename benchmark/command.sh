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