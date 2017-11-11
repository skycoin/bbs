# `bbscli`

This is the Skycoin BBS command-line interface executable. To have access to a BBS node, the node needs to have the `--rpc` flag set as true (default). Note that the `--rpc-port` sets the port in which this interface is served.

This is the help menu for `bbscli`:

```bash
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