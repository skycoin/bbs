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
BBS_DIS_PATH=${GOPATH}/src/github.com/skycoin/bbs/cmd/discoverynode/discoverynode.go

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

    go run ${BBS_DIS_PATH} \
        -address ${ADDRESS_MS} \
        &
}

RunNode() {
    if [[ $# -ne 5 ]] ; then
        echo "4 arguments required"
        exit 1
    fi

    ADDRESS_MS=$1 ; PORT_HTTP=$2 ; PORT_CXO=$3 ; PORT_RPC=$4 ; GUI=$5

    pv "START NODE: PORT_HTTP ${PORT_HTTP}, PORT_CXO ${PORT_CXO}, PORT_RPC ${PORT_RPC}, GUI ${GUI}..."

    go run ${BBS_NODE_PATH} \
        -memory=true \
        -enforced-messenger-addresses=${ADDRESS_MS} \
        -rpc-port=${PORT_RPC} \
        -cxo-port=${PORT_CXO} \
        -cxo-rpc=false \
        -web-port=${PORT_HTTP} \
        -web-gui=${GUI} \
        -web-gui-dir=${GOPATH}/src/github.com/skycoin/bbs/static/dist \
        -open-browser=${GUI} \
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
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; LOC=$2

    pv "NODE '${PORT}': IMPORT BOARD FROM '${LOC}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content import_board \
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

VoteThread() {
    if [[ $# -ne 6 ]] ; then
        echo "6 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2 ; THASH=$3 ; VALUE=$4 ; CSK=$5 ; TS=$6

    pv "NODE '${PORT}': VOTE THREAD '${THASH}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content vote_thread \
        -board-public-key="${BPK}" \
        -thread-hash="${THASH}" \
        -value="${VALUE}" \
        -creator-secret-key="${CSK}" \
        -timestamp=${TS}
}

VotePost() {
    if [[ $# -ne 6 ]] ; then
        echo "6 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2 ; PHASH=$3 ; VALUE=$4 ; CSK=$5 ; TS=$6

    pv "NODE '${PORT}': VOTE POST '${PHASH}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content vote_post \
        -board-public-key="${BPK}" \
        -post-hash="${PHASH}" \
        -value="${VALUE}" \
        -creator-secret-key="${CSK}" \
        -timestamp=${TS}
}

VoteUser() {
    if [[ $# -ne 7 ]] ; then
        echo "7 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2 ; UPK=$3 ; VALUE=$4 ; TAGS=$5 ; CSK=$6 ; TS=$7

    pv "NODE '${PORT}': VOTE USER '${UPK}'"

    go run ${BBS_CLI_PATH} -p ${PORT} content vote_user \
        -board-public-key="${BPK}" \
        -user-public-key="${UPK}" \
        -value="${VALUE}" \
        -tags=${TAGS} \
        -creator-secret-key="${CSK}" \
        -timestamp=${TS}
}