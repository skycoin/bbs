#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "skycoin bbs binary dir:" "$DIR"

pushd "$DIR" >/dev/null

go run cmd/bbsnode/bbsnode.go \
    -web-gui-dir="${DIR}/static/dist" \
    -enforced-messenger-addresses="35.185.110.6:8080" \
    -enforced-subscriptions="03cfd850ed4df2e43c8e6359c00de574df4a233b710ee44252166a5097f1c6056d" \
    -open-browser=true
    $@

popd >/dev/null