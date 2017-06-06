#!/usr/bin/env bash

# Build executables.

echo "[ BUILDING EXECUTABLES ]"
echo "> cxod ..."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
echo "> bbsnode ..."
go build $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go

echo "[ STARTING BBS NODE ]"
echo "> CXO DAEMON ..."
./cxod &
sleep 5
echo "> BBS SERVER ..."
./bbsnode \
    --test-mode=true \
    --test-mode-threads=3 \
    --test-mode-users=50 \
    --test-mode-min=0 \
    --test-mode-max=1

# Clean up.

wait
echo "[ CLEANING UP ]"
rm cxod bbsnode
echo "Goodbye!"