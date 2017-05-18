#!/usr/bin/env bash

cxo_daemon_port=8988
cxo_rpc_port=8987
bbs_rpc_port=6481
bbs_gui_port=6480

echo "> Building Executables."
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cli/cli.go
go build $GOPATH/src/github.com/evanlinjin/bbs/main.go

echo "> Starting CXO Daemon."
./cxod \
    --address=[::]:$cxo_daemon_port \
    --rpc-address=[::]:$cxo_rpc_port \
    &

sleep 5

echo "> Adding a Feed."
./cli \
    --a=[::]:$cxo_rpc_port \
    --e='add_feed 032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b'

echo "> Adding Another Feed."
./cli \
    --a=[::]:$cxo_rpc_port \
    --e='add_feed 02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b'

echo "> Connecting to Other CXO Daemon."
./cli \
    --a=[::]:$cxo_rpc_port \
    --e='connect [::]:8998'

echo "> Starting BBS Server."
./main \
    --master=true \
    --cxo-port=$cxo_daemon_port \
    --rpc-server-port=$bbs_rpc_port \
    --rpc-server-remote-address=127.0.0.1:$bbs_rpc_port \
    --web-gui-port=$bbs_gui_port

wait
echo "> Quitting. Removing Unnecessary Files."
rm cli cxod main *.bak