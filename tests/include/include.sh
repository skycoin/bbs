#!/usr/bin/env bash

inMac() {
    if ! type "curl" > /dev/null; then
        echo "curl is not installed."
        exit 1
    fi
    if ! type "jq" > /dev/null; then
        echo "jq is not installed."
        exit 1
    fi
}

inLinux() {
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

    pv "START MESSENGER SERVER: PORT_MS ${PORT_MS}..."

    go run ${GOPATH}/src/github.com/skycoin/bbs/cmd/devsd/devsd.go \
        -address=${ADDRESS_MS} \
        &
}

RunNode() {
    if [[ $# -ne 3 ]] ; then
        echo "3 arguments required"
        exit 1
    fi

    PORT_HTTP=$1 ; PORT_CXO=$2 ; GUI=$3

    pv "START NODE: PORT_HTTP ${PORT_HTTP}, PORT_SUB ${PORT_SUB}, PORT_CXO ${PORT_CXO}..."

    go run ${GOPATH}/src/github.com/skycoin/bbs/cmd/bbsnode/bbsnode.go \
        -dev=true \
        -master=true \
        -memory=true \
        -defaults=false \
        -cxo-port=${PORT_CXO} \
        -cxo-rpc=false \
        -http-port=${PORT_HTTP} \
        -http-gui=${GUI} \
        &
}

# <<< SESSION >>>

NewUser() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; SEED=$2

    pv "NODE '${PORT}': NEW USER WITH SEED '${SEED}'"

    curl \
        --noproxy "*" \
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
        --noproxy "*" \
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
        --noproxy "*" \
        -X POST \
        -sS "http://127.0.0.1:${PORT}/api/session/logout" | jq
}

# <<< CONNECTION >>>

NewConnection() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; ADDRESS=$2

    pv "NODE '${PORT}': NEW CONNECTION '${ADDRESS}'"

    curl \
        --noproxy "*" \
        -X POST \
        -F "address=${ADDRESS}" \
        -sS "http://127.0.0.1:${PORT}/api/connections/new" | jq
}

DeleteConnection() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; ADDRESS=$2

    pv "NODE '${PORT}': DELETE CONNECTION '${ADDRESS}'"

    curl \
        --noproxy "*" \
        -X POST \
        -F "address=${ADDRESS}" \
        -sS "http://127.0.0.1:${PORT}/api/connections/delete" | jq
}

# <<< SUBSCRIPTION >>>

NewSubscription() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2

    pv "NODE '${PORT}': NEW SUBSCRIPTION '${BPK}'"

    curl \
        --noproxy "*" \
        -X POST \
        -F "public_key=${BPK}" \
        -sS "http://127.0.0.1:${PORT}/api/subscriptions/new" | jq
}

DeleteSubscription() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2

    pv "NODE '${PORT}': NEW SUBSCRIPTION '${BPK}'"

    curl \
        --noproxy "*" \
        -X POST \
        -F "public_key=${BPK}" \
        -sS "http://127.0.0.1:${PORT}/api/subscriptions/new" | jq
}

# <<< CONTENT >>>

NewBoard() {
    if [[ $# -ne 2 ]] ; then
        echo "2 arguments required"
        exit 1
    fi

    PORT=$1 ; NAME=$2
    
    pv "NODE '${PORT}': NEW BOARD '${NAME}'"

    curl \
        --noproxy "*" \
        -X POST \
        -F "seed=${NAME}" \
        -F "name=Board ${NAME}" \
        -F "body=A board generated with seed '${NAME}'." \
        -sS "http://127.0.0.1:${PORT}/api/content/new_board" | jq
}

NewThread() {
    if [[ $# -ne 4 ]] ; then
        echo "4 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2 ; NAME=$3 ; BODY=$4

    pv "NODE '${PORT}': NEW THREAD '${NAME}'"

    curl \
        --noproxy "*" \
        -X POST \
        -F "board_public_key=${BPK}" \
        -F "name=${NAME}" \
        -F "body=${BODY}" \
        -sS "http://127.0.0.1:${PORT}/api/content/new_thread" | jq
}

NewTestThread() {
    PORT=$1 ; BPK=$2 ; NUMBER=$3
    NewThread ${PORT} ${BPK} "Test Thread ${NUMBER}" "This is test thread of index ${NUMBER}."
}

NewPost() {
    if [[ $# -ne 5 ]] ; then
        echo "5 arguments required"
        exit 1
    fi

    PORT=$1 ; BPK=$2 ; TREF=$3 ; NAME=$4 ; BODY=$5

    pv "NODE '${PORT}': NEW POST '${NAME}'"

    curl \
        --noproxy "*" \
        -X POST \
        -F "board_public_key=${BPK}" \
        -F "thread_ref=${TREF}" \
        -F "name=${NAME}" \
        -F "body=${BODY}" \
        -sS "http://127.0.0.1:${PORT}/api/content/new_post" | jq
}