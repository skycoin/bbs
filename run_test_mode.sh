#!/usr/bin/env bash

echo "[ STARTING BBS NODE ]"
echo "> CXO DAEMON ..."
go run $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go &
echo "> BBS SERVER ..."
go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    --test-mode=true \
    --test-mode-threads=4 \
    --test-mode-users=50 \
    --test-mode-min=0 \
    --test-mode-max=0 \
    --test-mode-timeout=20 \
    --test-mode-post-cap=200 \
wait
echo "Goodbye!"