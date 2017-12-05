# Content Submission Process

The content submission process has two steps (or two endpoint calls).

1. **Preparation Step** - to prepare the content to be signed. The input is of fields that is to be controlled by the user. What is returned is the object and it's hash, where the hash is to be signed. The preparation step endpoints are as follows:
    * `/api/submission/prepare_thread`
    * `/api/submission/prepare_post`
    * `/api/submission/prepare_thread_vote`
    * `/api/submission/prepare_post_vote`
    * `/api/submission/prepare_user_vote`

2. **Finalization Step** - User signs the prepared content's hash provided by the preparation step's output. The user submits the hash and the signature and is presented with data, dependent of what is submitted. There is only one endpoint:
    * `/api/submission/finalize`

## Examples

For the following examples, we will assume the following:

1. All content (threads, posts and votes) are submitted to board of public key `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` (generated with seed `a`).

2. All submissions are created by a user that's represented by a public/private key pair that's generated with seed `user`:

| Name | Value |
| --- | --- |
| Public Key | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` |
| Private Key | `8705518acec973239f704aa1bdbf7f5300f006682d8f6b435976e49c8b62aab0` |

### Thread Submission Example

The first step is to submit to the thread preparation endpoint (`/api/submission/prepare_thread`):

| Key | Value | Description |
| --- | --- | --- |
| `of_board` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | public key of board to submit thread to |
| `name` | `Test Thread` | name of the thread to submit |
| `body` | `This is a thread.` | body of the thread to submit |
| `creator` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | public key of the user that "creates" this thread |

What is returned, is a json object of this nature:

```json
{
    "okay": true,
    "data": {
        "hash": "dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3",
        "raw": "{\"type\":\"5,thread\",\"ts\":1512450011135448190,\"of_board\":\"032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b\",\"name\":\"Test Thread\",\"body\":\"This is a thread.\",\"creator\":\"0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718\"}"
    }
}
```

At this stage, the front end needs to perform two actions.
1. Ensure that the `"raw"` data does generate a hash of that provided under `"hash"`.
2. Ensure that the data within `"raw"` is expected (that the node did not randomly change fields unexpectedly).

After it is confirmed that everything is okay, we can sign the hash with the creator's private key (`8705518acec973239f704aa1bdbf7f5300f006682d8f6b435976e49c8b62aab0`) to obtain:

```json
{
    "input": {
        "hash": "dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3",
        "secret_key": "8705518acec973239f704aa1bdbf7f5300f006682d8f6b435976e49c8b62aab0"
    },
    "sig": "6265825cf1bdb6e06cec837e8f4c4c2051e8bff9ad04893a0ffb22eee0fdc37600de5e0995011832e447e8ea3d0ffb415fe28e2f42076aa1c8f17f70a4e59b7300"
}
```

After obtaining the signature, we can finalize the submission via the `/api/submission/finalize` endpoint:

| Key | Value | Description |
| --- | --- | --- |
| `hash` | `dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3` | hash of content that needs submission finalization |
| `sig` | `6265825cf1bdb6e06cec837e8f4c4c2051e8bff9ad04893a0ffb22eee0fdc37600de5e0995011832e447e8ea3d0ffb415fe28e2f42076aa1c8f17f70a4e59b7300` | signature of the hash, generated with the creator's private key |

What is returned here would be something of this appearance (if successful):

```json
{
    "okay": true,
    "data": {
        "board": {
            "public_key": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
            "header": {},
            "body": {
                "type": "5,board",
                "ts": 1512351592165056559,
                "name": "Board A",
                "body": "Board generated with 'a'.",
                "submission_keys": [
                    "127.0.0.1:8080,0302cc8bca0b42a4e084dca0dc2c8c774e2dac1062d182103d7623ad63165e1eeb"
                ]
            }
        },
        "threads": [
            {
                "header": {
                    "hash": "dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3",
                    "sig": "6265825cf1bdb6e06cec837e8f4c4c2051e8bff9ad04893a0ffb22eee0fdc37600de5e0995011832e447e8ea3d0ffb415fe28e2f42076aa1c8f17f70a4e59b7300"
                },
                "body": {
                    "type": "5,thread",
                    "ts": 1512450011135448190,
                    "of_board": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
                    "name": "Test Thread",
                    "body": "This is a thread.",
                    "creator": "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718"
                }
            }
        ]
    }
}
```

### Post Submission Example

First the post preparation endpoint (`/api/submission/prepare_post`) is used with the following values:

| Key | Value | Description |
| --- | --- | --- |
| `of_board` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | public key of board to submit post under |
| `of_thread` |  `dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3` | hash of thread to submit post under |
| `name` | `Test Post` | name of this post |
| `body` | `This is a post.` | body of post |
| `creator` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | public key of the creator |

Note that we have left the `of_post` field out, as this post is not a reply of another post.

The endpoint returns the following json object:

```json
{
    "okay": true,
    "data": {
        "hash": "d37a9e42002555aea6023fcf37c875e976dfc4aa6869c4b1b6fee7935ca06e9f",
        "raw": "{\"type\":\"5,post\",\"ts\":1512450789320762587,\"of_board\":\"032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b\",\"of_thread\":\"dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3\",\"name\":\"Test Post\",\"body\":\"This is a test post.\",\"creator\":\"0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718\"}"
    }
}
```

The front-end confirms and verifies this data (like in the previous example) and signs the given hash with the user's secret key:

```json
{
    "input": {
        "hash": "d37a9e42002555aea6023fcf37c875e976dfc4aa6869c4b1b6fee7935ca06e9f",
        "secret_key": "8705518acec973239f704aa1bdbf7f5300f006682d8f6b435976e49c8b62aab0"
    },
    "sig": "a64c0e6f7eb7e97cf85eeeb1775080f1ef7f0419bfa1160159e5f6f511cd27ae174c3b8717b1ecb103e7f61ae3573b9df9157b790c32ed3129c4e080501764ee00"
}
```

Then we can finalize the submission via `/api/submission/finalize`:

| Key | Value | Description |
| --- | --- | --- |
| `hash` | `d37a9e42002555aea6023fcf37c875e976dfc4aa6869c4b1b6fee7935ca06e9f` | hash of content that needs submission finalization |
| `sig` | `a64c0e6f7eb7e97cf85eeeb1775080f1ef7f0419bfa1160159e5f6f511cd27ae174c3b8717b1ecb103e7f61ae3573b9df9157b790c32ed3129c4e080501764ee00` | signature of the hash, generated with the creator's private key |

On success, this returns something of this structure:

```json
{
    "okay": true,
    "data": {
        "board": {
            "public_key": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
            "header": {},
            "body": {
                "type": "5,board",
                "ts": 1512448628817443192,
                "name": "Board A",
                "body": "Board generated with 'a'.",
                "submission_keys": [
                    "127.0.0.1:8080,0372b7d014be6ae181378b27d8ab681a66f3e0a2659313d67ae3804f07f460296b"
                ]
            }
        },
        "thread": {
            "header": {
                "hash": "dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3",
                "sig": "6265825cf1bdb6e06cec837e8f4c4c2051e8bff9ad04893a0ffb22eee0fdc37600de5e0995011832e447e8ea3d0ffb415fe28e2f42076aa1c8f17f70a4e59b7300"
            },
            "body": {
                "type": "5,thread",
                "ts": 1512450011135448190,
                "of_board": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
                "name": "Test Thread",
                "body": "This is a thread.",
                "creator": "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718"
            }
        },
        "posts": [
            {
                "header": {
                    "hash": "d37a9e42002555aea6023fcf37c875e976dfc4aa6869c4b1b6fee7935ca06e9f",
                    "sig": "a64c0e6f7eb7e97cf85eeeb1775080f1ef7f0419bfa1160159e5f6f511cd27ae174c3b8717b1ecb103e7f61ae3573b9df9157b790c32ed3129c4e080501764ee00"
                },
                "body": {
                    "type": "5,post",
                    "ts": 1512450789320762587,
                    "of_board": "032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b",
                    "of_thread": "dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3",
                    "name": "Test Post",
                    "body": "This is a test post.",
                    "creator": "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718"
                }
            }
        ]
    }
}
```

### Thread Vote Submission Example

We first use the thread-vote preparation endpoint (`/api/submission/prepare_thread_vote`) with the following values:

| Key | Value | Description |
| --- | --- | --- |
| `of_board` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | public key of board in which to submit thread vote |
| `of_thread` | `dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3` | hash of thread to cast vote on |
| `value` | `+1` | vote value (-1, 0, +1) |
| `creator` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | public key of the creator of the thread vote |

Thus returning:

```json
{
    "okay": true,
    "data": {
        "hash": "e9889e73b80853fd0a907037a84bdc11a5aa6892c1b6a09e2406c8dca7ade97a",
        "raw": "{\"type\":\"5,thread_vote\",\"ts\":1512457522020970909,\"of_board\":\"032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b\",\"of_thread\":\"dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3\",\"value\":1,\"creator\":\"0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718\"}"
    }
}
```

Which we can sign and submit to `/api/submission/finalize`:

| Key | Value | Description |
| --- | --- | --- |
| `hash` | `e9889e73b80853fd0a907037a84bdc11a5aa6892c1b6a09e2406c8dca7ade97a` | hash of content that needs submission finalization |
| `sig` | `d6b472b711f996cb85597b11c988cfe48fc3081c5ad4fe8d0be347bd4646c6911f61edb53615320da239ab167fc7a95b1f56a32e97302330fb44f7ab52b1715b01` | signature of the hash, generated with the creator's private key |

Thus returning:

```json
{
    "okay": true,
    "data": {
        "votes": {
            "ref": "dd3314649f162aeafbc6034f61d9ef70526543b8a398a72abb9065bef1f89fe3",
            "up_votes": {
                "voted": true,
                "count": 1
            },
            "down_votes": {
                "voted": false,
                "count": 0
            }
        }
    }
}
```

