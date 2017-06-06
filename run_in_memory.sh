#!/usr/bin/env bash

echo "[ STARTING BBS NODE ]"
go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    --master=true \
    --save-config=false \
    --cxo-use-internal=true \
    --cxo-memory-mode=true \
    --web-gui-open-browser=true
echo "Goodbye!"