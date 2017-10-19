#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "skycoin binary dir:" "$DIR"

pushd "$DIR" >/dev/null

go run cmd/bbsnode/bbsnode.go --http-gui-dir="${DIR}/static/dist" $@

popd >/dev/null