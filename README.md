# Skycoin BBS
Skycoin BBS is a next generation decentralised social network (BBS stands for [Bulletin Board System](https://en.wikipedia.org/wiki/Bulletin_board_system)).

Skycoin BBS uses the [Skycoin CX Object System](https://github.com/skycoin/cxo) (CXO) to store and synchronise data between nodes.  

There are many configurations for running a Skycoin BBS Node.
* By default, a node cannot host/create new boards. But a node can be set as master to enable such an ability.
* A node can be set to have it's own internal CXO Daemon, or use and external CXO Daemon that is running as a separate process.
* A node can also be configured to run purely in memory (RAM) so nothing is stored on disk.
* Additionally, you can run a node in test mode where it simulates the actions of users. This is useful for developers or those setting up master nodes.

## Build Skycoin BBS

Get source code, dependencies and build BBS Node:
```bash
go get github.com/skycoin/bbs/cmd/bbsnode
```

If the CXO Daemon is to run as a separate process:
```bash
go get github.com/skycoin/cxo/cmd/cxod
```

If a command line interface for the CXO Daemon is required:
```bash
go get github.com/skycoin/cxo/cmd/cli
```

For all the commands above, the executables will be in `$GOPATH/bin`.

## Run Skycoin BBS

By default (aka running `bbsnode` without specifying flags), a BBS Node will have the following characteristics:

* **Not in master mode -** This means that the node does not have the ability to create and host boards. However, it can subscribe to boards of other nodes it is connected to and interact with those.

* **Saves configuration and cxo files -** Configuration files for boards, users and message queue will be stored in the `$HOME/.skybbs` directory. CXO files will be in `$HOME/.skybbs/cxo`.

* **Runs an internal cxo daemon -** The cxo daemon will run as a goroutine within the executable.

* **Uses the following ports -** The CXO Daemon will listen on `8998` and `8997` (RPC). The BBS Node will host it's web interface and JSON API on `7410`.

* **Launches the browser -** After the `bbsnode` has successfully started, the browser will be launched to display the web interface.

### Examples

The following examples assume that you have `$GOPATH/bin` in your `$PATH`, and the above executables have been built.

#### Run a node as master

Master nodes have the ability to create and host boards.

```bash
bbsnode \
    -master=true \
    -rpc-port=1234 \
    -rpc-remote-address=34.215.131.150:1234 \
    -web-gui-enable=false
```

Nodes that run as master will host an RPC server. Other nodes add posts, threads and votes, on hosted boards via this RPC connection. Hence, a port needs to be specified of where the RPC connection is hosted, as well as the remote address of the connection. The flags `rpc-port` and `rpc-remote-address` are used to specify these.

If a graphical user interface is not needed for the node, setting `web-gui-enable` to false can disable it.

#### Run a node in memory

This mode is ideal for testing, or if you don't wish to save anything to disk.

```bash
bbsnode \
    -save-config=false \
    -cxo-memory-mode=true
```

The `save-config` flag determines whether or not to save the Skycoin BBS configuration files for boards, users and messages. Here, it is set as false.

When the `cxo-memory-mode` flag is set to true, instead of storing objects and configurations to file, the cxo client and server will store everything in memory. When the application exits, everything stored locally will be deleted.
