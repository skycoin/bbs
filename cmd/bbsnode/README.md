# `bbsnode`

This is the main Skycoin BBS executable. It runs a node that can host/subscribe to boards and connect to other nodes. It can host a thin client where users can submit content and interact with the board.

This is the help menu for `bbsnode`:

```
$ bbsnode -h


NAME:
   bbsnode - Runs a Skycoin BBS Node

USAGE:
   bbsnode [global options] command [command options] [arguments...]

VERSION:
   5.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --memory                              avoid storing BBS data on disk and use memory instead
   --config-dir value                    the name of the directory to store and access BBS configuration and associated cxo data (if left blank, $HOME/.skybbs will be used)
   --rpc                                 whether to enable RPC interface to interact with BBS node (used for bbscli)
   --rpc-port value                      port to serve BBS RPC interface (default: 8996)
   --cxo-port value                      port to listen for CXO connections (default: 8998)
   --cxo-rpc                             whether to enable RPC interface to interact with CXO (used for cxocli) (default: true)
   --cxo-rpc-port value                  port to serve CXO RPC interface (default: 8997)
   --enforced-messenger-addresses value  list of addresses to messenger servers to enforce connections with
   --enforced-subscriptions value        list of public keys of boards to enforce subscriptions with
   --web-port value                      port to serve http api (default: 8080)
   --web-gui                             whether to enable web interface thin client
   --web-gui-dir value                   directory where web interface static files are located
   --web-tls                             whether to enable https for web interface thin client and api
   --web-tls-cert-file value             path of the tls certificate file
   --web-tls-key-file value              path of the tls key file
   --open-browser                        whether to open a browser window
   --help, -h                            show help
   --version, -v                         print the version
   
```

The following command runs a BBS node with an enforced connect with the main messenger server and subscribed to the BBS Community board.

```bash
$ bbsnode -enforced-messenger-addresses="messenger.skycoin.net:8080" \
          -enforced-subscriptions="03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742"
```