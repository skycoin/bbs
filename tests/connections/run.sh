#!/usr/bin/env bash

# ------------------------------------------------------- #
# This script does the following things:                  #
#  1. Runs 5 BBS nodes in memory.                         #
#  2. 'Connects' the last four nodes with first BBS node. #
#  3. Injects a 'filled' board for each node.             #
#  4. Subscribes first node to all the boards available.  #
# ------------------------------------------------------- #

# Check if curl is installed.
if ! dpkg -s curl > /dev/null ; then
    echo "curl is not installed."
    exit 1
fi

# Check if jq is installed.
if ! dpkg -s jq > /dev/null ; then
    echo "jq is not installed."
    exit 1
fi

# Prints awesome stuff.
pv () {
    echo "[ • ]" $1
}

pv2 () {
    echo "[ • ] --- ((( ${1} ))) ---"
}

# Generates a public key from seed.
GenPK() {
    if [[ $# -ne 1 ]] ; then
        echo "1 argument required"
        exit 1
    fi
    go run $GOPATH/src/github.com/skycoin/bbs/cmd/genpublickey/genpublickey.go $1
}

# Runs a BBS Node with specified ports.
RunNode() {
    if [[ $# -ne 5 ]] ; then
        echo "5 arguments required"
        exit 1
    fi
    pv "STARTING BBS NODE AT: ${1}, ${2}, ${3}, ${4}..."
    PORT_BBS_GUI=$1
    PORT_BBS_RPC=$2
    PORT_CXO_SERVER=$3
    PORT_CXO_RPC=$4
    OPEN_GUI=$5
    go run $GOPATH/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
        --master=true \
        --memory-mode=true \
        --rpc-port=$PORT_BBS_RPC \
        --rpc-remote-address=127.0.0.1:$PORT_BBS_RPC \
        --cxo-port=$PORT_CXO_SERVER \
        --cxo-rpc-port=$PORT_CXO_RPC \
        --web-gui-port=$PORT_BBS_GUI \
        --web-gui-open-browser=$OPEN_GUI \
        --web-gui-dir=$GOPATH/src/github.com/skycoin/bbs/static/dist \
        &
}

# Injects a test board on a specified node with seed.
InjectFilledBoard() {
    if [[ $# -ne 5 ]] ; then
        echo "5 arguments required"
        exit 1
    fi
    pv "INJECTING A BOARD WITH SEED ${2} ON NODE ${1}..."
    curl \
        -X POST \
        -F "seed=${2}" \
        -F "threads=${3}" \
        -F "min_posts=${4}" \
        -F "max_posts=${5}" \
        -sS "http://127.0.0.1:${1}/api/tests/add_filled_board" | jq
#    sleep 1
}

# Subscribes a node to a board.
SubscribeToBoard() {
    if [[ $# -ne 3 ]] ; then
        echo "3 arguments required"
        exit 1
    fi
    pv "SUBSCRIBING TO BOARD ${3} ON NODE ${1}..."
    curl \
        -X POST \
        -F "address=127.0.0.1:${2}" \
        -F "board=${3}" \
        -sS "http://127.0.0.1:${1}/api/subscriptions/add" | jq
    sleep 1
}

# Run 5 nodes.
pv2 "RUNNING SOME NODES"
RunNode 7410 7411 7412 7413 true
RunNode 7420 7421 7422 7423 false
RunNode 7430 7431 7432 7433 false
RunNode 7440 7441 7442 7443 false
RunNode 7450 7451 7452 7453 false

pv2 "SLEEPING 10s"
sleep 10

# Make some filled boards on the nodes.
pv2 "INJECTING FILLED BOARDS ON THE NODES"
InjectFilledBoard 7410 a 5 10 20
InjectFilledBoard 7420 b 5 10 20
InjectFilledBoard 7430 c 5 10 20
InjectFilledBoard 7440 d 5 10 20
InjectFilledBoard 7450 e 5 10 20

# Subscribe first node to all boards of other nodes.
pv2 "SUBSCRIBING FIRST NODE TO OTHER BOARDS"
SubscribeToBoard 7410 7422 $(GenPK b)
SubscribeToBoard 7410 7432 $(GenPK c)
SubscribeToBoard 7410 7442 $(GenPK d)
SubscribeToBoard 7410 7452 $(GenPK e)

pv "ALL DONE!"
wait