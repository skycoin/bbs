# Skycoin BBS Node

This is the main executable for running a node.

For general information regarding Skycoin BBS, visit the [README.md](https://github.com/skycoin/bbs/blob/master/README.md) in the repository's root directory.

## Build and Run

```bash
# Get stuff.
go get github.com/skycoin/bbs/cmd/bbsnode

# Execute.
bbsnode

# Execute with some flags.
bbsnode \
    -save-config=false \
    -cxo-memory-mode=true
```

## Flags

The following is a detailed description of all the available flags for `bbsnode`.

* `-master` (boolean) \
    Determines whether the node is to be started as master or not. Master mode allows the node to create and host boards.
    
* `-save-config` (boolean) \
    Determines
    
***!!! TO BE COMPLETED !!!***
