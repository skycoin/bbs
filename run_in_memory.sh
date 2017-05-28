#!/usr/bin/env bash

# Build executables.

echo "[ BUILDING EXECUTABLES ]"
echo "> cxod ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
echo "> cli ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cli/cli.go
echo "> bbsnode ..."
go build $GOPATH/src/github.com/evanlinjin/bbs/cmd/bbsnode/bbsnode.go

echo "[ STARTING BBS NODE ]"
echo "> CXO DAEMON ..."
./cxod \
    --mem-db=true \
    &
sleep 5
echo "> ADDING FEEDS ..."
./cli \
    --e='add_feed 032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b'
echo "> BBS SERVER ..."
./bbsnode \
    --master=true \
    --save-config=false \
    --cxo-memory-mode=true \
    --web-gui-open-browser=true

# Clean up.

wait
echo "[ CLEANING UP ]"
rm cli cxod bbsnode
echo "Goodbye!"