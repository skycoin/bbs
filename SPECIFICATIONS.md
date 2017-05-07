# Skycoin BBS Specifications

## Data Structure
* Root
  * Board
  * BoardThreads
    * Thread 0
    * Thread 1
    * Thread 2
  * ThreadPosts 0
    * Post 0
  * ThreadPosts 1
    * Post 1
    * Post 2
  * ThreadPosts 2
    * Post 3

## For JSON API
Refer to [JSON API Specifications on Dropbox Paper](https://paper.dropbox.com/doc/JSON-API-Specifications-S6BHC351LStxlgySl55M2) and access the [Postman Collection](https://www.getpostman.com/collections/39a4168cf9be0c746f47).

P.S. Command line arguments for making bbs server master (So you can create boards).
```
--master=true --launch-browser=false
```

## Temporary Specifications

### Overarching Instructions

#### V1
* It's p2p
* Instead of HTTP, use meshnet.
* Data Model: Board, Thread, Post.
* Board needs name and timestamp.
* Expose JSON API (Use fake data for now where database (cxo) should be)
* [PACKAGE] `"datastore"`
  * Implement methods such as: GetBoards, GetPosts, etc. (But return fake data)
  * cxo will be hooked up here in the future.

#### V2
* Use [skycoin/src/gui](https://github.com/skycoin/skycoin/tree/master/src/gui) as example/template for webserver.
* Make submodule that subscribes to pubkey and acts as CXO client. And provides data needed to drive the URLs.
* Commandline arguments - use [skycoin/cmd/skycoin/skycoin.go](https://github.com/skycoin/skycoin/blob/master/cmd/skycoin/skycoin.go) as example.
* Implement master/slave ???

#### V3
1. Start webserver and routes, like in [src/gui](https://github.com/skycoin/skycoin/tree/master/src/gui).

2. Create json api, and routes for adding post, getting posts for thread, getting list of threads on board, getting list of boards. We will add urls/functions as needed.

3. Keep a list of board subscriptions (public keys).

4. Command line argument to point server at CXO daemon (ours will be a client).

5. Module for getting the data from CXO, to serve to app.

6. command line argument, for optionally opening up port for RPC (like the CLI RPC type thing), for remote nodes to submit posts to a master node (who can then add the posts/theads in CXO and write the data)

7. Ask me and I will do the schema.

8. On connect to CXO server, subscribe to the pubkeys for each board in our list of boards we are subscribed to.

9. Have a state struct per board and keep track of whether local server is "master" for that board (can write/update board) and []byte field for storing the private key for writing updates, if the server is actually master for that board.

#### V4
1. User signs post and sends to master server. Then master server adds it to board.

2. Master server needs rpc port and server for commands like add post / create thread.

3. A User subscribes to a board. A board is public key publishing data. A board subscription is a subscription in CXO.

4. Subscription to board is CXO client telling CXO daemon to subscribe to public key data feed (which is board). The public key is owner/moderator of the board.

5. Post's signature is the user who wrote the post. Signing a hash of the post to prove he wrote the post.

6. Client stores a profile: public/private key. Sign posts with private key.

### For RPC

Only two methods needed.
1. Inject post.
2. New thread.