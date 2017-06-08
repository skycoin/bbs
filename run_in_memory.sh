#!/usr/bin/env bash

# Port for web graphical user interface.
PORT_BBS_GUI=6480
# Port for bbs rpc (For cross-node communication).
PORT_BBS_RPC=6481
# Port for cxo server.
PORT_CXO_SERVER=8988
# Port for cxo rpc (for cross-cxo communication).
PORT_CXO_RPC=8987

go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
    --master=true \
    --save-config=false \
    --rpc-server-port=$PORT_BBS_RPC \
    --rpc-server-remote-address=127.0.0.1:$PORT_BBS_RPC \
    --cxo-use-internal=true \
    --cxo-port=$PORT_CXO_SERVER \
    --cxo-rpc-port=$PORT_CXO_RPC \
    --cxo-memory-mode=true \
    --web-gui-port=$PORT_BBS_GUI \
    --web-gui-open-browser=true \
    --web-gui-dir=$GOPATH/src/github.com/skycoin/bbs/static/dist
