# Skycoin BBS Node

This is the main executable for running a node.

For general information regarding Skycoin BBS, visit the [README.md](https://github.com/skycoin/bbs/blob/master/README.md) in the repository's root directory.

## Build and Run

```bash
# Get stuff.
go get github.com/skycoin/bbs/cmd/bbsnode

# Execute.
bbsnode

# Execute showing available flags and exit.
bbsnode --help

# Execute with some flags.
bbsnode \
    -save-config=false \
    -cxo-memory-mode=true
```