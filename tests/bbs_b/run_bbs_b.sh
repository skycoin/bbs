#!/usr/bin/env bash

b_cxod=8988
b_cxorpc=8987
b_bbsrpc=6481
b_bbsgui=6480

# Build executables.

echo "[ BUILDING EXECUTABLES ]"
echo "> cxod ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
echo "> cli ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cli/cli.go
echo "> main ..."
go build $GOPATH/src/github.com/evanlinjin/bbs/main.go

# Start BBS Node 'B'.

echo "[ STARTING BBS NODE 'B' ]"
echo "> CXO DAEMON ..."
./cxod \
    --address=[::]:$b_cxod \
    --rpc-address=[::]:$b_cxorpc \
    &
sleep 5
echo "> ADDING FEEDS ..."
./cli \
    --a=[::]:$b_cxorpc \
    --e='add_feed 032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b'
./cli \
    --a=[::]:$b_cxorpc \
    --e='add_feed 02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b'
echo "> CONNECTING TO DAEMON A ..."
./cli \
    --a=[::]:$b_cxorpc \
    --e='connect [::]:8998'
echo "> BBS SERVER ..."
./main \
    --master=true \
    --cxo-port=$b_cxod \
    --rpc-server-port=$b_bbsrpc \
    --rpc-server-remote-address=127.0.0.1:$b_bbsrpc \
    --web-gui-port=$b_bbsgui

# Cleanup.

wait
echo "[ CLEANING UP ]"
rm cli cxod main *.bak *.json