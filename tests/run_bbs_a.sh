#!/usr/bin/env bash

# Build execs.
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cxod/cxod.go
go build $GOPATH/src/github.com/skycoin/cxo/cmd/cli/cli.go
go build $GOPATH/src/github.com/evanlinjin/bbs/main.go

cmd_cxod="./cxod"

# TODO: Complete.