#!/usr/bin/env bash

## WARNING: WORK IN PROGRESS.

# Print function.
pv () {
    echo "[ • ]" $1
}

# Ensure GOPATH environment variable is set.
: ${GOPATH:?"Please set GOPATH before building"}
pv "GOPATH=${GOPATH}"

# Make directory environment variables.
ROOT_DIR=$GOPATH/src/github.com/skycoin/bbs
BUILD_DIR=$ROOT_DIR/build

# Make directories.
mkdir $PWD/build

go get github.com/skycoin/bbs/cmd/bbsnode
go get github.com/skycoin/bbs/cmd/bbscli