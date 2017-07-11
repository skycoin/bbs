#!/usr/bin/env bash

# Usage: bash connect_node_to_bbs_community_server.sh [PORT OF YOUR WEB GUI]

if [[ $# -ne 1 ]] ; then
    echo "1 argument required. Please specify web gui port of your local BBS node."
    exit 1
fi

PORT=$1

echo "> SUBSCRIBING TO BOARD..."
curl \
    -X POST \
    -F "address=34.204.161.180:8210" \
    -F "board=03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742" \
    -sS "http://127.0.0.1:${PORT}/api/subscriptions/add" | jq

echo "> FINISHED."