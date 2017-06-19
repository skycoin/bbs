#!/usr/bin/env bash

echo "[ STARTING BBS NODE ]"
# Set test-mode-timeout=-1 to disable.
# Set test-mode-post-cap=-1 to disable.
go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    -test-mode=true \
    -test-mode-threads=4 \
    -test-mode-users=50 \
    -test-mode-min=0 \
    -test-mode-max=1 \
    -test-mode-timeout=-1 \
    -test-mode-post-cap=200
echo "Goodbye!"