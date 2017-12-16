#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "skycoin bbs binary dir:" "$DIR"

pushd "$DIR" >/dev/null

go run cmd/bbsnode/bbsnode.go \
    -web-gui-dir="${DIR}/static/dist" \
    -enforced-messenger-addresses="35.227.102.45:8005" \
    -enforced-subscriptions="03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742" \
    -open-browser=true
    $@

popd >/dev/null