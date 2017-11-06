![skycoin bbs logo](https://user-images.githubusercontent.com/26845312/32426755-274b72b0-c282-11e7-989f-dc8870f4635e.png)

# Skycoin BBS

[![GoReportCard](https://goreportcard.com/badge/skycoin/bbs)](https://goreportcard.com/report/skycoin/bbs)
[![Telegram group link](telegram-group.svg)](https://t.me/skycoinbbs)

Skycoin BBS is a next generation decentralised social network (BBS stands for [Bulletin Board System](https://en.wikipedia.org/wiki/Bulletin_board_system)).

Skycoin BBS uses the [Skycoin CX Object System](https://github.com/skycoin/cxo) (CXO) to store and synchronise data between nodes.

[![Skycoin BBS Showcase 4 - YouTube](https://i.ytimg.com/vi/6ZqwgefYauU/0.jpg)](https://youtu.be/6ZqwgefYauU)

## Building Static Files

Detailed building instructions for static files can be found in [static/README.md](https://github.com/skycoin/bbs/blob/master/static/README.md).


## Running BBS Node

#### Dependencies

**Go Programming Language** (https://golang.org/doc/install)

After installation, ensure that the `GOPATH` environmental variable is set.

#### Via bash script

The script `run.sh` is provided to run BBS, serving static files in `static.dist`.

```bash
bash run.sh
```

The following flags can be used:

* `-dev` (default: `false`) Serves GUI static files from Skycoin BBS location in `$GOPATH`.

* `-master` (default: `false`) Enables the node to host a port for submitting content to boards.

* `-memory` (default: `false`) Disables the node from saving to disk. By default; user files, cxo database, node connections and subscriptions are saved to disk.

* `-config-dir` (default: `""`) Sets the directory used to store configuration files of Skycoin BBS. Leave blank to use default location of `$HOME/.skybbs`.

* `-cxo-port` (default: `8998`) Port that CXO listens on (self-replication database).

* `-cxo-rpc` (default: `false`) Whether to enable CXO RPC Port (for admin control).

* `-cxo-rpc-port` (default: `8997`) Port used for CXO RPC (if enabled).

* `-http-port` (default: `7410`) Port to serve JSON API and GUI.

* `-http-gui` (default: `true`) Enables serving GUI.

* `-http-gui-dir` (default: `""`) Set's directory where static files are to be served from. Leave blank to use `./static/dist` (unless if `-dev` flag is set).


## Using Skycoin BBS

There are currently two ways of interacting with Skycoin BBS.
* **Web interface -** By default, the flag `-http-gui` is enabled. Hence, when BBS is launched, the web gui will be opened via the system browser.

* **Restful json api -** This is ideal for controlling nodes without a graphical user interface (in a server), or for building applications or administrator tools. Documentation for the api is provided as a [Postman](https://www.getpostman.com/) Collection located at [docfiles/postman_collection.json](https://raw.githubusercontent.com/skycoin/bbs/master/docfiles/postman_collection.json).

## Participate

#### Telegram

* [Community Chat](https://t.me/skycoinbbs) - Get up to date with development and talk to the developers.
* [Board Hosting Channel](https://t.me/skycoinbbshosting) - Get a list of nodes to connect to and boards to subscribe to.
