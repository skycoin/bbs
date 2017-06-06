#!/usr/bin/env bash

echo "[ STARTING BBS NODE ]"
echo "> CXO DAEMON ..."
go run $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go \
    --mem-db=true \
    &
echo "> BBS SERVER ..."
go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    --master=true \
    --save-config=false \
    --cxo-memory-mode=true \
    --web-gui-open-browser=true
wait
echo "Goodbye!"