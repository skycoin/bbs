#!/usr/bin/env bash

# Port for web graphical user interface.
PORT_BBS_GUI=7410

go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    -master=true \
    -memory-mode=true \
    -web-gui-port=$PORT_BBS_GUI \
    -web-gui-dir=""
