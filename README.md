![skycoin bbs logo](https://user-images.githubusercontent.com/26845312/32426755-274b72b0-c282-11e7-989f-dc8870f4635e.png)

# Skycoin BBS

[![GoReportCard](https://goreportcard.com/badge/skycoin/bbs)](https://goreportcard.com/report/skycoin/bbs)
[![Telegram group link](telegram-group.svg)](https://t.me/skycoinbbs)
[![Build Status](https://travis-ci.org/skycoin/bbs.svg?branch=ui_automation)](https://travis-ci.org/skycoin/bbs)

Skycoin BBS uses the [Skycoin CX Object System](https://github.com/skycoin/cxo) (CXO) to store and synchronise data between nodes and the [Skycoin Messenger](https://github.com/skycoin/net) (net) for inter-node content submission.

[![Skycoin BBS Showcase 4 - YouTube](https://i.ytimg.com/vi/6ZqwgefYauU/0.jpg)](https://youtu.be/6ZqwgefYauU)

## Install

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

## Run

```bash
$ bbsnode
```

For more detailed instructions:
* [cmd/bbsnode/README.md](./cmd/bbsnode/README.md)
* [Wiki: Setting up a Skycoin BBS Node](https://github.com/skycoin/bbs/wiki/Setting-up-a-Skycoin-BBS-Node)

The script [`run.sh`](./run.sh) is provided as a convenient to run BBS, serving static files in `static/dist`.

```bash
$ ./run.sh
```

## Docker

Pull docker image.

```bash
$ docker pull skycoin/bbs
```

Create a docker volume.

```bash
$ docker volume create bbs-data
```

Run Skycoin BBS.

```bash
$ docker run -p 8080:8080 -p 8998:8998 -p 8996:8996 -v bbs0:/data skycoin/bbs
```

List network interfaces.

```bash
$ ifconfig
```

Use CLI.

```bash
# help menu
$ docker run skycoin/bbs bbscli -h

# interact with bbs node
$ docker run skycoin/bbs bbscli -a 172.17.0.1:8996 messengers discover
```

## Command-line interface

The Command-line interface is for administration control over the BBS node.

Detailed instructions are located at [cmd/bbscli/README.md](./cmd/bbscli/README.md).

## Documentation

Please make use of the [Skycoin BBS Wiki](https://github.com/skycoin/bbs/wiki)!

## Participate

#### Telegram

* [Community Chat](https://t.me/skycoinbbs) - Get up to date with development and talk to the developers.
* [Board Hosting Channel](https://t.me/skycoinbbshosting) - Get a list of nodes to connect to and boards to subscribe to.
