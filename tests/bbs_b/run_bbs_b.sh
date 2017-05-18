#!/usr/bin/env bash

cxo_daemon_port=8988
cxo_rpc_port=8987
bbs_rpc_port=6481
bbs_gui_port=6480

# Build execs.
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cli/cli.go
go build $GOPATH/src/github.com/evanlinjin/bbs/main.go

./cxod \
    --address=[::]:$cxo_daemon_port \
    --rpc-address=[::]:$cxo_rpc_port \
    &

sleep 5

./cli \
    --a=[::]:$cxo_rpc_port \
    --e='add_feed 032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b'

./cli \
    --a=[::]:$cxo_rpc_port \
    --e='connect [::]:8998'

./main \
    --master=true \
    --rpc-server-port=$bbs_rpc_port \
    --rpc-server-remote-address=127.0.0.1:$bbs_rpc_port \
    --web-gui-port=$bbs_gui_port

wait
