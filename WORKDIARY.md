# Work Diary
Information regarding development progress.
* **TODO** shows what is to be implemented next.
* **DONE** shows what has been implemented (split into days).

## TODO
* Make use of RootWalker functionality.
* Test RPC functionality.

## DONE

### 2017 May 13
* Updated JSON API.
* Made it work with the new CXO.

### 2017 May 9
* Implemented adding new posts and threads via RPC. However this is not tested.

### 2017 May 8
* Fix bug where we are unable to post to thread.
* Posts need to be signed before posting.
* Finish JSON API structure according to [specifications](https://paper.dropbox.com/doc/JSON-API-Specifications-S6BHC351LStxlgySl55M2).

### 2017 May 7
 * Add Thread to board and view board with threads implemented according to specifications. This includes access via JSON API.

### 2017 May 6
* Implement loading/saving `BoardConfig` from/into JSON file.
* Split `Board` into `Board` and `BoardThreads`.
  * Hence, `Board` is for the metadata of the board. Name, etc.
  * `BoardThreads` will be the list of Threads in the board and the number of total Threads.

### 2017 May 5
* Implement `UserManager`.
  * Load user config as JSON file on startup.
  * Create new user if config does not exist, and save to config.
