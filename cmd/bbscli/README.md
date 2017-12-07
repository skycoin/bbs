# `bbscli`

This is the Skycoin BBS command-line interface executable. To have access to a BBS node, the node needs to have the `--rpc` flag set as true (default). Note that the `--rpc-port` sets the port in which this interface is served.

## Installation

To install, ensure that the `GOPATH` environmental variable is set, and that the `$GOPATH/bin` directory is in your `PATH`.

Run the following to download and compile the `bbscli` binary:

```bash
$ go get github.com/skycoin/bbs/cmd/bbscli
```

### Enable command-line autocomplete

For bash, run the following command:

```bash
$ PROG=bbscli source $GOPATH/src/github.com/skycoin/bbs/cmd/bbscli/autocomplete/bash_autocomplete
```

If you are using zsh, please replace the `bash_autocomplete` with `zsh_autocomplete` in the previous command.

To avoid running the command every time you start a new terminal session, copy the script into the `~/.bashrc` or `~/.zshrc` file.

## Environmental Settings

### `BBS_RPC_PORT`

`bbscli` will connect to port `8996` by default. You can change the port by setting the `BBS_RPC_PORT` env variable with the following command:

```bash
$ export BBS_RPC_PORT=8996
```

## Usage

This is the help menu for `bbscli`:

```text
$ bbscli -h


NAME:
   bbscli - a command-line interface to interact with a Skycoin BBS node

USAGE:
   bbscli [global options] command [command options] [arguments...]

VERSION:
   5.0

COMMANDS:
     tools          cryptography tools
     messengers     manages messenger connections of the node
     connections    manages connections of the node
     subscriptions  manages subscriptions of the node
     content        manages boards and their content
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value  rpc port of the bbs node (default: 8996) [$BBS_RPC_PORT]
   --help, -h              show help
   --version, -v           print the version

```