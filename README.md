![skycoin bbs logo](https://user-images.githubusercontent.com/26845312/32426755-274b72b0-c282-11e7-989f-dc8870f4635e.png)

# Skycoin BBS

[![GoReportCard](https://goreportcard.com/badge/skycoin/bbs)](https://goreportcard.com/report/skycoin/bbs)
[![Telegram group link](telegram-group.svg)](https://t.me/skycoinbbs)

Skycoin BBS is a next generation decentralised social network (BBS stands for [Bulletin Board System](https://en.wikipedia.org/wiki/Bulletin_board_system)).

Skycoin BBS uses the [Skycoin CX Object System](https://github.com/skycoin/cxo) (CXO) to store and synchronise data between nodes and the [Skycoin Messenger](https://github.com/skycoin/net) (net) for inter-node content submission.

[![Skycoin BBS Showcase 4 - YouTube](https://i.ytimg.com/vi/6ZqwgefYauU/0.jpg)](https://youtu.be/6ZqwgefYauU)

## Installation

### Go 1.9+ Installation and Setup

[Golang 1.9+ Installation/Setup](https://github.com/skycoin/skycoin/blob/develop/Installation.md)

After installation, ensure that the `GOPATH` environmental variable is set and that `$GOPATH/bin` is added to the `PATH` environment variable.

### Dependencies

Dependencies are managed with [dep](https://github.com/golang/dep).

To install `dep`:

```sh
go get -u github.com/golang/dep
```

`dep` vendors all dependencies into the repo.

### Download and Compile BBS Executables

```sh
go get https://github.com/skycoin/bbs/...
```

This will download `github.com/skycoin/bbs` to `$GOPATH/src/github.com/skycoin/bbs`.

You can also clone the repo directly with `git clone https://github.com/skycoin/bbs`,
but it must be cloned to this path: `$GOPATH/src/github.com/skycoin/bbs`.

### Static Files For The Web Thin Client

Building instructions for static files can be found in [static/README.md](./static/README.md).

## Running BBS Node

```bash
$ bbsnode
```

For more detailed instructions:
* [cmd/bbsnode/README.md](./cmd/bbsnode/README.md)
* [Wiki: Setting up a Skycoin BBS Node](./wiki/Setting-up-a-Skycoin-BBS-Node)

The script [`run.sh`](./run.sh) is provided as a convenient to run BBS, serving static files in `static/dist`.

```bash
$ ./run.sh
```

## Using Skycoin BBS

There are currently two ways of interacting with Skycoin BBS.
* **Web interface thin client -** By default, the flag `-web-gui` is enabled. Hence, when BBS is launched, the web gui will be served at a port specified by `-web-port`. One can only submit and view content via the thin client.

* **Restful json api -** This is ideal for viewing/submitting content without a graphical user interface. Documentation for the api is provided as a [Postman](https://www.getpostman.com/) Collection located at [doc/bbs_postman_collection.json](doc/postman_collection.json) which can be viewed online at: [https://documenter.getpostman.com/view/985347/skycoin-bbs-v05/719YYTS](https://documenter.getpostman.com/view/985347/skycoin-bbs-v05/719YYTS). A brief written documentation is provided at [doc/api_explnation.md](./doc/api_explanantion.md).

* **Command-line interface -** This is for administration tools. Detailed instructions are located at [cmd/bbscli/README.md](./cmd/bbscli/README.md).

## Documentation

Please make use of the [Skycoin BBS Wiki](./wiki)!

## Participate

#### Telegram

* [Community Chat](https://t.me/skycoinbbs) - Get up to date with development and talk to the developers.
* [Board Hosting Channel](https://t.me/skycoinbbshosting) - Get a list of nodes to connect to and boards to subscribe to.
