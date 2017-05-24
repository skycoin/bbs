# TODO
A good-old todo list.


## Urgent
Not really urgent, but things to do now.

* Fix bug where config files aren't saved in the right places. They should really be in `$HOME/.skycoin/bbs`. See how it's done in CXO and replicate that method.
* Fix deleting post/thread bug.
* Add a `test-mode` flag. This creates random posts under the influence of a timer. We will use this to test and benchmark syncing. The `test-mode` flag will be accompanied with:
  * `test-threads` - specifies how many threads to create.
  * `interval` - specifies the time between when posts are created (posts will be dropped on a random thread - even ones we don't own).
  * Make a `Tester` class in `extern/dev/tester.go`. This class will call functions in `extern/gui/gateway.go`.
* When signing/checking the hash sum of a post, remove the timestamp from within signature. Hence, the board's master can set the timestamp on injection to board/thread.
* More tests. Golaong tests, and bash scripts to test RPC.
## Later
We need to talk about these together. Think about it for now.

* Board Permissions. Have way to configure a board to allow/disallow people to post on it.