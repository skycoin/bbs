# Work Diary
Information regarding development progress.
* **TODO** shows what is to be implemented next.
* **DONE** shows what has been implemented (split into days).

## TODO
* Finish JSON API structure according to [specifications](https://paper.dropbox.com/doc/JSON-API-Specifications-S6BHC351LStxlgySl55M2).

## DONE

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
