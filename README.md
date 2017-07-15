# Skycoin BBS

[![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/bbs)](https://goreportcard.com/report/github.com/skycoin/bbs)

Skycoin BBS is a next generation decentralised social network (BBS stands for [Bulletin Board System](https://en.wikipedia.org/wiki/Bulletin_board_system)).

Skycoin BBS uses the [Skycoin CX Object System](https://github.com/skycoin/cxo) (CXO) to store and synchronise data between nodes.  

[![Skycoin BBS Alpha Showcase - YouTube](https://img.youtube.com/vi/OBaaSB369TI/0.jpg)](https://youtu.be/OBaaSB369TI)

## Building Skycoin BBS

First ensure that `go`, `npm`, `zip` and `ng` are installed and the `GOPATH` environment variable is set.

#### BBS Node

To get the source code, dependencies and build BBS Node:
```bash
go get github.com/skycoin/bbs/cmd/bbsnode
```
The executables will be in `$GOPATH/bin`.

Optionally, you can then package BBS Node by running bash script provided: `pkg/package.sh`.
```bash
cd $GOPATH/src/github.com/skycoin/bbs/pkg

bash package.sh
```
Built binaries and static files will be located in `pkg/build/` folder.

## Running Skycoin BBS

By default (aka running `bbsnode` without specifying flags), a BBS Node will have the following characteristics:

* **Not in master mode -** This means that the node does not have the ability to create and host boards. However, it can subscribe to boards of other nodes it is connected to and interact with those.

* **Saves configuration and cxo files -** Configuration files for boards, users and message queue will be stored in the `$HOME/.skybbs` directory. CXO files will be in `$HOME/.skybbs/cxo`.

* **Uses the following ports -** The CXO Daemon will listen on `8998` and `8997` (RPC). The BBS Node will host it's web interface and JSON API on `7410`.

* **Serves static files in `./static/dist` -** If no static files are found, `bbsnode` will panic. However, you can specify flag `-web-gui-dir=""` and `bbsnode` will serve static files in `$GOPATH/src/github.com/skycoin/bbs/static/dist` (assuming that the `GOPATH` env has been set).

* **Launches the browser -** After the `bbsnode` has successfully started, the browser will be launched to display the web interface.

### Examples

The following examples assume that you have `$GOPATH/bin` in your `$PATH`, and the above executables have been built.

#### Show help dialog and exit
```bash
bbsnode --help
```

#### Run a node with static files in `$GOPATH/src/github.com/bbs/static/dist`

First ensure that the `GOPATH` env is set and you have the Skycoin BBS module.

```bash
bbsnode -web-gui-dir=""
```

#### Run a node as master

Master nodes have the ability to create and host boards.

```bash
bbsnode \
    -master=true \
    -rpc-port=8234 \
    -rpc-remote-address=34.215.131.150:8234 \
    -cxo-port=8456 \
    -web-gui-enable=false
```
Other nodes connect to the master node via it's CXO Port (specified with `cxo-port` flag). If the external IP address of the server of the above example is `34.215.131.150`, then other nodes can connect via the address `34.215.131.150:8456` (assuming that the port `8456` is shared).

Nodes that run as master will host a BBS RPC server. Other nodes add posts, threads and votes, on hosted boards via this RPC connection. Hence, a port needs to be specified of where the RPC connection is hosted, as well as the remote address of the connection. The flags `rpc-port` and `rpc-remote-address` are used to specify these.

If a graphical user interface is not needed for the node, setting `web-gui-enable` to false can disable it.

#### Run a node in memory

This mode is ideal for testing, or if you don't wish to save anything to disk.

```bash
bbsnode -memory-mode
```

The `save-config` flag determines whether or not to save the Skycoin BBS configuration files for boards, users and messages. Here, it is set as false.

When the `cxo-memory-mode` flag is set to true, instead of storing objects and configurations to file, the cxo client and server will store everything in memory. When the application exits, everything stored locally will be deleted.

#### Run a node in test mode

In test mode, the node generates a board, and creates threads, posts and votes as different users. It does this in intervals.

```bash
bbsnode \
    -test-mode \
    -test-mode-threads=4 \
    -test-mode-users=50 \
    -test-mode-min=1 \
    -test-mode-max=5 \
    -test-mode-timeout=10 \
    -test-mode-post-cap=200
```
Test mode forces the node to be run as master. BBS configuration files will not be stored in test mode, and CXO files will be stored in `/tmp` and deleted after the node exits.

* `test-mode-threads` specifies how many threads are to be created on the test board.
* `test-mode-users` specifies how many simulated users are to be created.
* `test-mode-min` is the minimum interval in seconds between simulated activity.
* `test-mode-max` is the maximum interval in seconds between simulated activity.
* `test-mode-timeout` is the time in seconds before all simulated activity stops. This can be set as `-1` to disable the timeout.
* `test-most-post-cap` is the maximum number of posts that are allowed to be created. This can be set as `-1` to disable the cap.

## Using Skycoin BBS

There are currently two ways of interacting with Skycoin BBS.
* **Web interface -** By default, the flags `web-gui-enable` and `web-gui-open-browser` are enabled. Hence, when BBS is launched, the web gui will be opened via the system browser.

* **Restful json api -** This is ideal for controlling nodes without a graphical user interface (in a server), or for building applications or administrator tools. Documentation for the api is provided as a [Postman](https://www.getpostman.com/) Collection located at [docfiles/postman_collection.json](https://raw.githubusercontent.com/skycoin/bbs/master/docfiles/postman_collection.json).

## Participate

* [Telegram](https://t.me/joinchat/B_ax-ERzBNAAQDuqgGgOUQ) (Skycoin BBS Development)