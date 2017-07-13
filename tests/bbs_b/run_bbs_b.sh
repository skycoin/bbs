#!/usr/bin/env bash

b_cxod=8988
b_cxorpc=8987
b_cxodir=bbs_b_server
b_bbsrpc=6481
b_bbsgui=6480

# Start BBS Node 'B'.

echo "[ STARTING BBS NODE 'B' ]"

go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    -master=true \
    -memory-mode=true \
    -cxo-port=$b_cxod \
    -cxo-rpc-port=$b_cxorpc \
    -rpc-port=$b_bbsrpc \
    -rpc-remote-address=127.0.0.1:$b_bbsrpc \
    -web-gui-port=$b_bbsgui \
    -web-gui-open-browser=false \
    -web-gui-dir=""

echo "Goodbye!"