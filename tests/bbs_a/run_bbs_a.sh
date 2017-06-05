#!/usr/bin/env bash

a_cxod=8998
a_cxorpc=8997
a_cxodir=bbs_a_server
a_bbsrpc=6491
a_bbsgui=6490

# Build executables.

echo "[ BUILDING EXECUTABLES ]"
echo "> cxod ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
echo "> bbsnode ..."
go build $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go

# Start BBS Node 'A'.

echo "[ STARTING BBS NODE 'A' ]"
echo "> CXO DAEMON ..."
./cxod \
    --address=[::]:$a_cxod \
    --rpc-address=[::]:$a_cxorpc \
    --mem-db=true \
    --data-dir=$a_cxodir \
    &
echo "> BBS SERVER ..."
./bbsnode \
    --master=true \
    --save-config=false \
    --cxo-port=$a_cxod \
    --cxo-rpc-port=$a_cxorpc \
    --cxo-memory-mode=true \
    --cxo-dir=bbs_a \
    --rpc-server-port=$a_bbsrpc \
    --rpc-server-remote-address=127.0.0.1:$a_bbsrpc \
    --web-gui-port=$a_bbsgui \
    --web-gui-open-browser=false

# Cleanup.

wait
echo "[ CLEANING UP ]"
rm cxod bbsnode
echo "Goodbye!"