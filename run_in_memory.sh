#!/usr/bin/env bash

# Build executables.

echo "[ BUILDING EXECUTABLES ]"
echo "> cxod ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
echo "> bbsnode ..."
go build $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go

echo "[ STARTING BBS NODE ]"
echo "> CXO DAEMON ..."
./cxod \
    --mem-db=true \
    &
echo "> BBS SERVER ..."
./bbsnode \
    --master=true \
    --save-config=false \
    --cxo-memory-mode=true \
    --web-gui-open-browser=true

# Clean up.

wait
echo "[ CLEANING UP ]"
rm cxod bbsnode
echo "Goodbye!"