#!/usr/bin/env bash

a_cxod=8998
a_cxorpc=8997
a_cxodir=bbs_a_server
a_bbsrpc=6491
a_bbsgui=6490

# Start BBS Node 'A'.

echo "[ STARTING BBS NODE 'A' ]"

go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    -master=true \
    -memory-mode=true \
    -cxo-port=$a_cxod \
    -cxo-rpc-port=$a_cxorpc \
    -rpc-port=$a_bbsrpc \
    -rpc-remote-address=127.0.0.1:$a_bbsrpc \
    -web-gui-port=$a_bbsgui \
    -web-gui-open-browser=false \
    -web-gui-dir=""

echo "Goodbye!"