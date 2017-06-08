# Skycoin BBS
Next generation decentralised social network.
BBS stands for [Bulletin Board System](https://en.wikipedia.org/wiki/Bulletin_board_system).

Skycoin BBS uses the [Skycoin CX Object System](https://github.com/skycoin/cxo) (CXO) to store and synchronise data between nodes.  

## Run Skycoin BBS
There are many configurations for running a Skycoin BBS Node.
* By default, a node cannot host/create new boards. But a node can be set as master to enable such an ability.
* The node can be set to have it's own internal CXO Daemon, or use and external CXO Daemon that is running as a separate process.
* A node can be configured to run purely in memory (RAM) so nothing is stored on disk.

### Build and run a node with default configurations
Get source, dependencies and build:
```bash
$ go get github.com/skycoin/bbs/cmd/bbsnode
```
Run:
```bash
$ cd $GOPATH/bin
$ ./bbsnode
```
_If `$GOPATH/bin` is already a part of your path, running `bbsnode` will work too._

### Run a node in memory

Don't want to make any changes to disk?
```bash
$ ./run_in_memory.sh
```
