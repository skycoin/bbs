#!/usr/bin/env bash

a_cxod=8998
a_cxorpc=8997
a_bbsrpc=6491
a_bbsgui=6490

# Build executables.

echo "[ BUILDING EXECUTABLES ]"
echo "> cxod ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
echo "> cli ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cli/cli.go
echo "> main ..."
go build $GOPATH/src/github.com/evanlinjin/bbs/main.go

# Start BBS Node 'A'.

echo "[ STARTING BBS NODE 'A' ]"
echo "> CXO DAEMON ..."
./cxod \
    --address=[::]:$a_cxod \
    --rpc-address=[::]:$a_cxorpc \
    &
sleep 5
echo "> ADDING FEEDS ..."
./cli \
    --a=[::]:$a_cxorpc \
    --e='add_feed 032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b'
./cli \
    --a=[::]:$a_cxorpc \
    --e='add_feed 02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b'
echo "> BBS SERVER ..."
./main \
    --master=true \
    --cxo-port=$a_cxod \
    --rpc-server-port=$a_bbsrpc \
    --rpc-server-remote-address=127.0.0.1:$a_bbsrpc \
    --web-gui-port=$a_bbsgui

# Cleanup.

wait
echo "[ CLEANING UP ]"
rm cli cxod main *.bak *.json