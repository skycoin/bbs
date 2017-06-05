#!/usr/bin/env bash

b_cxod=8988
b_cxorpc=8987
b_cxodir=bbs_b_server
b_bbsrpc=6481
b_bbsgui=6480

# Build executables.

echo "[ BUILDING EXECUTABLES ]"
echo "> cxod ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
echo "> cli ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cli/cli.go
echo "> bbsnode ..."
go build $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go

# Start BBS Node 'B'.

echo "[ STARTING BBS NODE 'B' ]"
echo "> CXO DAEMON ..."
./cxod \
    --address=[::]:$b_cxod \
    --rpc-address=[::]:$b_cxorpc \
    --mem-db=true \
    --data-dir=$b_cxodir \
    &
sleep 5
echo "> CONNECTING TO DAEMON A ..."
./cli \
    --a=[::]:$b_cxorpc \
    --e='connect [::]:8998'
echo "> BBS SERVER ..."
./bbsnode \
    --master=true \
    --save-config=false \
    --cxo-port=$b_cxod \
    --cxo-rpc-port=$b_cxorpc \
    --cxo-memory-mode=true \
    --cxo-dir=bbs_b \
    --rpc-server-port=$b_bbsrpc \
    --rpc-server-remote-address=127.0.0.1:$b_bbsrpc \
    --web-gui-port=$b_bbsgui \
    --web-gui-open-browser=false

# Cleanup.

wait
echo "[ CLEANING UP ]"
rm cli cxod bbsnode
echo "Goodbye!"