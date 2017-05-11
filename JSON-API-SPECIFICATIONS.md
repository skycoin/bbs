# JSON API Specifications
JSON API should be implemented in file `[bbs/gui/api.go](https://github.com/evanlinjin/bbs/blob/master/gui/api.go)`.
The struct  `gui.APIHandler` will be exposed to `cxo.Gateway` for methods.
For now, as `cxo.Gateway` is not completed yet, so return fake data.
Data types are stored in `[bbs/typ](https://github.com/evanlinjin/bbs/tree/master/typ)` submodule. Please use types `Board`, `Thread` and `Post` for the JSON API.

# Overview
| **Method** | **URI**                         | **Operation**                                                                 |
| ---------- | ------------------------------- | ----------------------------------------------------------------------------- |
| GET        | /api/boards                     | Lists all boards we are subscribed to.                                        |
| GET        | /api/boards/BOARD_ID            | Show information of board, identified with BOARD_ID, and list it’s threads.   |
| GET        | /api/boards/PUBLIC_ID/THREAD_ID | Show information of thread, identified with THREAD_ID, and list it’s posts.   |
| PUT        | /api/boards                     | Creates a new board.                                                          |
| PUT        | /api/boards/BOARD_ID            | Creates a new thread in board of BOARD_ID.                                    |
| PUT        | /api/boards/PUBLIC_ID/THREAD_ID | Creates a new post in thread of THREAD_ID (which is under board of BOARD_ID). |

# Details
## GET /api/boards

Lists all boards we are subscribed to.

*Example Reply Body*

    {
      "boards": [
        {
          "public_key": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
          "name": "Board of Cool People",
          "master": true,
          "created": 1494117382,
        },
        {
          "public_key": "134dftabt5534aa33ab4516333688bdfd220faae7f35349112df07e971eb36fdb4",
          "name": "Board of Uncool People",
          "master": false,
          "created": 1493596800,
        }
      ]
    }
## GET /api/boards/BOARD_ID

Show information of board, identified with BOARD_ID, and list it’s threads.

*Example Reply Body*

    {
      "board": {
        "public_key": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
        "name": "Board of Cool People",
        "master": true,
        "created": 1494117382,
      },
      "threads": [
        {
          "name": "Barack Obama",
          "description": "Barack Hussein Obama II is an American politician who served as the 44th President of the United States from 2009 to 2017.",
          "created": 1494119917,
        },
        {
          "name": "Mahatma Gandhi",
          "description": "Mohandas Karamchand Gandhi was the leader of the Indian independence movement in British-ruled India.",
          "created": 1494121089,
        }
      ]
    }
## GET /api/boards/PUBLIC_ID/THREAD_ID

Show information of thread, identified with THREAD_ID, and list it’s posts.

*Example Reply Body*

    {
      "board": {
        "public_key": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
        "name": "Board of Cool People",
        "master": true,
        "created": 1494117382,
      },
      "thread": {
        "name": "Barack Obama",
        "description": "Barack Hussein Obama II is an American politician who served as the 44th President of the United States from 2009 to 2017.",
        "created": 1494119917,
      },
      "posts": [
        {
          "title": "Donald Trump is President",
          "body": "Holy cow, how is this even possible?"
          "author": "917e44b903455421ee0faa6b7f35cd3350ee1ff07e971eb36fdee6760688b2fb9d",
          "created": 1494119929,
        }
      ]
    }
## PUT /api/boards

Creates a new board, or subscribes to a board if `public_key` is provided.
Note that when creating a new board, we need a `seed` (for now).

*Example Request Body (Create new board)*

    {
      "seed": "random_seed",
      "board": {
        "name": "Board of Cool People"
      }
    }

*Example Request Body (Subscribe to a board)*

    {
      "board": {
        "public_key": "134dftabt5534aa33ab4516333688bdfd220faae7f35349112df07e971eb36fdb4"
      }
    }

*Example Reply Body (Success)*

    {
      "put_request": {
        "okay": true
      },
      "boards": [ ...Same as with GET /api/boards ... ]
    }

*Example Reply Body (Failed)*

    {
      "put_request": {
        "okay": false,
        "error": "bbs server is not master"
      },
      "boards": [ ...Same as with GET /api/boards ... ]
    }
## PUT /api/boards/BOARD_ID

Creates a new thread in board of BOARD_ID.

*Example Request Body*

    {
      "thread": {
        "name": "Barack Obama",
        "description": "Barack Hussein Obama II is an American politician who served as the 44th President of the United States from 2009 to 2017."
        }
    }

*Example Reply Body (Success)*

    {
      "put_request": {
        "okay": true
      },
      "board": { ... Same as with GET /api/boards ... },
      "threads": [ ... Same as with GET /api/boards ... ]
    }

*Example Reply Body (Failed)*

    {
      "put_request": {
        "okay": false,
        "error": "thread already exists"
      },
      "board": { ... Same as with GET /api/boards ... },
      "threads": [ ... Same as with GET /api/boards ... ]
    }
## PUT /api/boards/PUBLIC_ID/THREAD_ID

Creates a new post in thread of THREAD_ID (which is under board of BOARD_ID).

*Example Request Body*

    {
      "post": {
          "title": "Donald Trump is President",
          "body": "Holy cow, how is this even possible?"
        }
    }

*Example Reply Body (Success)*

    {
      "put_request": {
        "okay": true
      },
      "board": { ... Same as with GET /api/boards ... },
      "thread": { ... Same as with GET /api/boards ... },
      "posts": [ ... Same as with GET /api/boards ... ]
    }

*Example Reply Body (Failed)*

    {
      "put_request": {
        "okay": false,
        "error": "internal server error"
      },
      "board": { ... Same as with GET /api/boards ... },
      "thread": { ... Same as with GET /api/boards ... },
      "posts": [ ... Same as with GET /api/boards ... ]
    }

