#!/usr/bin/env bash

GOPATH=${HOME}/go

inMac() {
    if ! type "curl" > /dev/null; then
        echo "curl is not installed."
        exit 1
    fi
}

inLinux() {
    # Check if curl is installed.
    if ! dpkg -s curl > /dev/null ; then
        echo "curl is not installed."
        exit 1
    fi
}

sysOS=`uname -s`
if [ $sysOS == "Darwin" ];then
	inMac
elif [ $sysOS == "Linux" ];then
	inLinux
else
	echo "Other OS: $sysOS"
    exit 1
fi

BBS_NODE_PATH=${GOPATH}/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go
BBS_CLI_PATH=${GOPATH}/src/github.com/skycoin/bbs/cmd/bbscli/bbscli.go

# Prints awesome stuff.
pv () {
    echo "[ • ]" $1
}

pv2 () {
    echo "[ • ] --- ((( ${1} ))) ---"
}

RunMS() {
    if [[ $# -ne 1 ]] ; then
        echo "1 argument required"
        exit 1
    fi

    ADDRESS_MS=$1

    pv "START MESSENGER SERVER: PORT_MS ${ADDRESS_MS}..."

    go run ${GOPATH}/src/github.com/skycoin/bbs/cmd/devsd/devsd.go \
        -address="${ADDRESS_MS}" \
        &
}

RunNode() {
    if [[ $# -ne 4 ]] ; then
        echo "4 arguments required"
        exit 1
    fi

    PORT_HTTP=$1 ; PORT_CXO=$2 ; PORT_RPC=$3 ; GUI=$4

    pv "START NODE: PORT_HTTP ${PORT_HTTP}, PORT_CXO ${PORT_CXO}, PORT_RPC ${PORT_RPC}, GUI ${GUI}..."

    go run ${BBS_NODE_PATH} \
        -dev=true \
        -memory=true \
        -enforced-messenger-addresses=127.0.0.1:8080 \
        -rpc-port=${PORT_RPC} \
        -cxo-port=${PORT_CXO} \
        -cxo-rpc=false \
        -http-port=${PORT_HTTP} \
        -http-gui=${GUI} \
        &
}

# <<< CONNECTION >>>

NewConnection() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; ADDRESS=$2

    pv "NODE '${PORT}': NEW CONNECTION '${ADDRESS}'"

    go run ${BBS_CLI_PATH} -p ${PORT} connections new \
        -address="${ADDRESS}"
}

DeleteConnection() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; ADDRESS=$2

    pv "NODE '${PORT}': DELETE CONNECTION '${ADDRESS}'"

    go run ${BBS_CLI_PATH} -p ${PORT} connections del \
        -address="${ADDRESS}"
}

# <<< SUBSCRIPTION >>>

NewSubscription() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2

    pv "NODE '${PORT}': NEW SUBSCRIPTION '${BPK}'"

    go run ${BBS_CLI_PATH} -p ${PORT} subscriptions new \
        -public-key="${BPK}"
}

DeleteSubscription() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2

    pv "NODE '${PORT}': NEW SUBSCRIPTION '${BPK}'"

    go run ${BBS_CLI_PATH} -p ${PORT} subscriptions delete \
        -public-key="${BPK}"
}

# <<< CONTENT >>>

NewBoard() {
    if [[ $# -ne 4 ]] ; then
        echo "4 arguments required"
        exit 1
    fi

    PORT=$1 ; NAME=$2 ; BODY=$3 ; SEED=$4

    pv "NODE '${PORT}': NEW BOARD '${NAME}' WITH SEED '${SEED}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content new_board \
        -name="${NAME}" \
        -body="${BODY}" \
        -seed="${SEED}"
}

ExportBoard() {
    if [[ $# -ne 3 ]] ; then
        echo "3 arguments required"
        exit 1
    fi

    PORT=$1 ; PK=$2 ; LOC=$3

    pv "NODE '${PORT}': EXPORT BOARD '${PK}' TO '${LOC}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content export_board \
        -public-key="${PK}" \
        -file-path="${LOC}"
}

ImportBoard() {
    if [[ $# -ne 3 ]] ; then
        echo "3 arguments required"
        exit 1
    fi

    PORT=$1 ; SK=$2 ; LOC=$3

    pv "NODE '${PORT}': IMPORT BOARD FROM '${LOC}' TO SECRET '${SK}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content import_board \
        -secret-key="${SK}" \
        -file-path="${LOC}"
}

NewThread() {
    if [[ $# -ne 6 ]] ; then
        echo "6 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2 ; NAME=$3 ; BODY=$4 ; CSK=$5 ; TS=$6

    pv "NODE '${PORT}': NEW THREAD '${NAME}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content new_thread \
        -board-public-key="${BPK}" \
        -name="${NAME}" \
        -body="${BODY}" \
        -creator-secret-key="${CSK}" \
        -timestamp=${TS}
}

NewPost() {
    if [[ $# -ne 7 ]] ; then
        echo "7 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2 ; THASH=$3 ; NAME=$4 ; BODY=$5 ; CSK=$6 ; TS=$7

    pv "NODE '${PORT}': NEW POST '${NAME}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content new_post \
        -board-public-key="${BPK}" \
        -thread-hash="${THASH}" \
        -name="${NAME}" \
        -body="${BODY}" \
        -creator-secret-key="${CSK}" \
        -timestamp=${TS}
}