#!/usr/bin/env bash

echo "[ STARTING BBS NODE ]"
go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    --test-mode=true \
    --test-mode-threads=4 \
    --test-mode-users=50 \
    --test-mode-min=0 \
    --test-mode-max=1 \
    --test-mode-timeout=5 \
    --test-mode-post-cap=200 \
    --cxo-use-internal=true
echo "Goodbye!"