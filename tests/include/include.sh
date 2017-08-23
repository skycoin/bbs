#!/usr/bin/env bash

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

RunNode() {
    if [[ $# -ne 3 ]] ; then
        echo "3 arguments required"
        exit 1
    fi

    PORT_HTTP=$1 ; PORT_SUB=$2 ; PORT_CXO=$3

    pv "START NODE: PORT_HTTP ${PORT_HTTP}, PORT_SUB ${PORT_SUB}, PORT_CXO ${PORT_CXO}..."

    go run ${GOPATH}/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
        -master=true \
        -memory=true \
        -cxo-port=${PORT_CXO} \
        -cxo-rpc=false \
        -sub-port=${PORT_SUB} \
        -sub-addr="[::]:${PORT_SUB}" \
        -http-port=${PORT_HTTP} \
        -http-gui=false \
        &
}

NewUser() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; SEED=$2

    pv "NODE '${PORT}': NEW USER WITH SEED '${SEED}'"

    curl \
        -X POST \
        -F "seed=${SEED}" \
        -F "alias=${SEED}" \
        -sS "http://127.0.0.1:${PORT}/api/session/users/new" | jq
}

Login() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; SEED=$2

    pv "NODE '${PORT}': LOGIN '${SEED}'"

    curl \
        -X POST \
        -F "alias=${SEED}" \
        -sS "http://127.0.0.1:${PORT}/api/session/login" | jq
}

Logout() {
    if [[ $# -ne 1 ]] ; then
        echo "1 arguments required"
        exit 1
    fi

    PORT=$1

    pv "NODE '${PORT}': LOGOUT"

    curl \
        -X POST \
        -sS "http://127.0.0.1:${PORT}/api/session/logout" | jq
}

NewBoard() {
    if [[ $# -ne 5 ]] ; then
        echo "5 arguments required"
        exit 1
    fi

    # TODO: FINISH
}