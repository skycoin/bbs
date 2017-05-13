# JSON API Specifications
Version 3 candidate (Not even implemented)

## Boards

#### Subscribe to a board
Request:

`POST /api/boards/[:board_id]/subscribe`

#### Unsubscribe from a board
Request:

`POST /api/boards/[:board_id]/unsubscribe`

#### List all boards
Request:

`GET /api/boards`

#### Create a new board (master only)
Request:

`POST /api/boards/new`

```json
{
  "seed": "random_seed",
  "board": {
    "name": "Coolies",
    "description": "Board of cool people."
  }
}
```

#### Remove a board (master only)
Request:

`DELETE /api/boards/remove?board=[:board_id]`

## Threads

#### List all threads of all subscribed boards
Request:

`GET /api/threads`

#### List all threads of a board
Request:

`GET /api/boards/[:board_id]/threads`

#### Create a new thread in a board
Request:

`POST /api/boards/[:board_id]/threads/new`

```json
{
  "thread": {
    "name": "Why?",
    "description": "An important question."
  }
}
```

#### Remove a thread from a board (master only)
Request:

`REMOVE /api/boards/[:board_id]/threads/remove?thread=[:thread_id]`

## Posts

#### List all posts
Request:

`GET /api/posts`

#### List all posts of specified thread
Request:

`GET /api/threads/[:thread_id]/posts`

#### Create a new post in specified thread
Request:

`POST /api/threads/[:thread_id]/posts/new`

```json
{
  "post": {
    "title": "Because.",
    "body": "It is only because I'm so cool."
  }
}
```