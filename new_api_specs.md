# JSON API Specifications
Version 3 candidate (Not even implemented)

## Subscriptions

#### List subscriptions

`GET /api/subscriptions`

#### Subscribe to board

`POST /api/subscriptions/[:board_id]`

#### Unsubscribe to board

`DELETE /api/subscriptions/[:board_id]`

#### Check subscription

`GET /api/subscriptions/[:board_id]`

## Boards

#### List all boards

`GET /api/boards`

#### Create a new board (master only)

`POST /api/boards`

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

`DELETE /api/boards/[:board_id]`

## Threads

#### List all threads of all subscribed boards

`GET /api/threads`

#### List all threads of a board

`GET /api/boards/[:board_id]/threads`

#### Create a new thread in a board

`POST /api/boards/[:board_id]/threads`

```json
{
  "thread": {
    "name": "Why?",
    "description": "An important question."
  }
}
```

#### Remove a thread from a board (master only)

`REMOVE /api/boards/[:board_id]/threads/[:thread_id]`

## Posts

#### List all posts

`GET /api/posts`

#### List all posts of specified thread

`GET /api/threads/[:thread_id]/posts`

#### Create a new post in specified thread

`POST /api/threads/[:thread_id]/posts`

```json
{
  "post": {
    "title": "Because.",
    "body": "It is only because I'm so cool."
  }
}
```