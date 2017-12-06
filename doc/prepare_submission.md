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

### Post Vote Submission Example

We first use the post vote preparation endpoint (`/api/submission/prepare_post_vote`) with the following data:

| Key | Value | Description |
| --- | --- | --- |
| `of_board` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | public key of board in which to submit thread vote |
| `of_post` | `513e46db5d45a3c3ee18b41591535e9844d27ae282b9329a65b990a58db770ae` | hash of thread to cast vote on |
| `value` | `-1` | vote value (-1, 0, +1) |
| `creator` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | public key of the creator of the thread vote |

Something similar to the following will be returned:

```json
{
    "okay": true,
    "data": {
        "hash": "071960062b21f9bb30b8e3ff294d8d7688f0079b062c10886c24f22918f2a790",
        "raw": "{\"type\":\"5,post_vote\",\"ts\":1512476136349103076,\"of_board\":\"032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b\",\"of_post\":\"db5a12a80e10208407fb82ae075705ca967a11424e960d52af7f6056b05384af\",\"value\":-1,\"creator\":\"0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718\"}"
    }
}
```

We sign the returned hash with the user's secret key (like before) and send it to the finalization endpoint (`/api/submission/finalize`):

| Key | Value | Description |
| --- | --- | --- |
| `hash` | `071960062b21f9bb30b8e3ff294d8d7688f0079b062c10886c24f22918f2a790` | hash of content that needs submission finalization |
| `sig` | `9909d9daf139c9b1d5b49d0a4e3a195fd81c44058573ddf6c3e43b289e423237687a2fd546c7d95f948de1a00d0bc600dcae65ade6bb90c372a8a8a425b304a100` | signature of the hash, generated with the creator's private key |

Something similar to the following will be returned:

```json
{
    "okay": true,
    "data": {
        "votes": {
            "ref": "db5a12a80e10208407fb82ae075705ca967a11424e960d52af7f6056b05384af",
            "up_votes": {
                "voted": false,
                "count": 0
            },
            "down_votes": {
                "voted": true,
                "count": 1
            }
        }
    }
}
```

### User Vote Submission Example

Within each user vote, the value of the fields `value` and `tags` will determine the nature of the vote. This is represented in the following table:

| Value (`"value"`) | Tags (`"tags"`)  | Nature                                   |
| ----------------- | ---------------- | ---------------------------------------- |
| `1`               | `["trust"]`      | User A trusts User B.                    |
| `-1`              | `["spam"]`       | User A determines User B as a spammer.   |
| `-1`              | `["block"]`      | User A dislikes User B and blocks User B. |
| `-1`              | `["spam,block"]` | User A determines User B as a spammer and blocks User B. |

Let's generate 3 users to demonstrate user voting.

**User 1 (generated with seed `1`):**

```json
{
  "public_key": "02f46d2461e2c3aba0585efb5b2ddb8acb34f38a56865f8a2a3f10272e6de257c1",
  "secret_key": "12348e8a15fcce27de6c187a5ecace09af622d495474cb3280e5e614f8b789b5"
}
```

**User 2 (generated with seed `2`):**

```json
{
  "public_key": "0284da18e80d5ec08cf54ed9c86bbbce6bbd2838b8c700a373a5886e4de44ce895",
  "secret_key": "9b43b74f9737e15e36921b418ec5c31ebcc92025133240c9061b2846a88f2e0c"
}
```

**User 3 (generated with seed `3`):**

```json
{
  "public_key": "03f5bcfadd87e625bf62900a7d1ed673ce74034dbfc5d5c624cedd4612a8dc6d1c",
  "secret_key": "9b40df8b560259c41af5ca5049d3fcd010e925511325fe0c11e362a9d0cada60"
}
```

Now these three users will vote on user of public key `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` (generated with seed `user`). We will call this user "default user" for this example.

**User 1's actions:**

User 1 trusts the default user and hence, prepares the following vote via the `/api/submission/prepare_user_vote` endpoint:

| Key | Value | Description |
| --- | --- | --- |
| `of_board` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | public key of board in which to submit thread vote |
| `of_user` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | public key of user to cast vote on |
| `value` | `+1` | vote value (-1, 0, +1) |
| `tags` | `trust` | vote tags, separated by commas |
| `creator` | `02f46d2461e2c3aba0585efb5b2ddb8acb34f38a56865f8a2a3f10272e6de257c1` | public key of the creator of the thread vote |

Hence, returning:

```json
{
    "okay": true,
    "data": {
        "hash": "d87bdbd68ea221e8b8b3622f0476e555430e4a4fd5550c2a6bec8ef1685a5551",
        "raw": "{\"type\":\"5,user_vote\",\"ts\":1512525362048762752,\"of_board\":\"032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b\",\"of_user\":\"0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718\",\"value\":1,\"tags\":[\"trust\"],\"creator\":\"02f46d2461e2c3aba0585efb5b2ddb8acb34f38a56865f8a2a3f10272e6de257c1\"}"
    }
}
```

As the returned data is expected and valid, User 1 finalizes the submission via `/api/submission/finalize`, using a signature that is generated with the secret key `12348e8a15fcce27de6c187a5ecace09af622d495474cb3280e5e614f8b789b5`:

| Key | Value | Description |
| --- | --- | --- |
| `hash` | `d87bdbd68ea221e8b8b3622f0476e555430e4a4fd5550c2a6bec8ef1685a5551` | hash of content that needs submission finalization |
| `sig` | `8b5ac5c22a75f003fb68a354ddc72ff9567c4252667de0024f87fca4c3255a27427346d0da44d524af6f20a2cae78c6041e439d3798ad716299422339537665200` | signature of the hash, generated with the creator's private key |

The returned result is:

```json
{
    "okay": true,
    "data": {
        "user_public_key": "02f46d2461e2c3aba0585efb5b2ddb8acb34f38a56865f8a2a3f10272e6de257c1",
        "profile": {
            "trusted_count": 1,
            "trusted": [
                "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718"
            ],
            "marked_as_spam_count": 0,
            "marked_as_spam": [],
            "blocked_count": 0,
            "blocked": [],
            "trusted_by_count": 0,
            "trusted_by": [],
            "marked_as_spam_by_count": 0,
            "marked_as_spam_by": [],
            "blocked_by_count": 0,
            "blocked_by": []
        }
    }
}
```

This is User 1's profile. As shown, User 1 trusts one user, and that user has a public key of `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718`.

**User 2's actions:**

User 2 has a grudge against the default user but knows that the default user is not a spammer/bot. Hence, User 2 prepares the following vote:

| Key | Value | Description |
| --- | --- | --- |
| `of_board` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | public key of board in which to submit thread vote |
| `of_user` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | public key of user to cast vote on |
| `value` | `-1` | vote value (-1, 0, +1) |
| `tags` | `block` | vote tags, separated by commas |
| `creator` | `0284da18e80d5ec08cf54ed9c86bbbce6bbd2838b8c700a373a5886e4de44ce895` | public key of the creator of the thread vote |

Thus, returning:

```json
{
    "okay": true,
    "data": {
        "hash": "c612955ffa7e82d51548f3cf57ad494bf4b1d735eb27a2db0341a04d2b79231a",
        "raw": "{\"type\":\"5,user_vote\",\"ts\":1512530064808609933,\"of_board\":\"032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b\",\"of_user\":\"0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718\",\"value\":-1,\"tags\":[\"block\"],\"creator\":\"0284da18e80d5ec08cf54ed9c86bbbce6bbd2838b8c700a373a5886e4de44ce895\"}"
    }
}
```

User 2 finalises the submission via the endpoint `/api/submission/finalize` using a signature generated with hash `c612955ffa7e82d51548f3cf57ad494bf4b1d735eb27a2db0341a04d2b79231a` and secret key `9b43b74f9737e15e36921b418ec5c31ebcc92025133240c9061b2846a88f2e0c`:

| Key | Value | Description |
| --- | --- | --- |
| `hash` | `c612955ffa7e82d51548f3cf57ad494bf4b1d735eb27a2db0341a04d2b79231a` | hash of content that needs submission finalization |
| `sig` | `9c34485e8b03f193499fe39ab821d9f9a591672406d07ce6555bfc2a68b4a085243cfb8f88914bff94d6397ff547a49c1f25354e53519e376268424dcd5ac7aa00` |

The endpoint returns the following result:

```json
{
    "okay": true,
    "data": {
        "user_public_key": "0284da18e80d5ec08cf54ed9c86bbbce6bbd2838b8c700a373a5886e4de44ce895",
        "profile": {
            "trusted_count": 0,
            "trusted": [],
            "marked_as_spam_count": 0,
            "marked_as_spam": [],
            "blocked_count": 1,
            "blocked": [
                "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718"
            ],
            "trusted_by_count": 0,
            "trusted_by": [],
            "marked_as_spam_by_count": 0,
            "marked_as_spam_by": [],
            "blocked_by_count": 0,
            "blocked_by": []
        }
    }
}
```

**User 3's actions:**

User 3 is highly annoyed at the default user as User 3 deems the default user as complete scam (and possibly a bot). User 3 wishes to mark the default user as spam, and also block the default user.

User 3 prepares the following submission to the endpoint `/api/submission/prepare_user_vote`:

| Key | Value | Description |
| --- | --- | --- |
| `of_board` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | public key of board in which to submit thread vote |
| `of_user` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | public key of user to cast vote on |
| `value` | `-1` | vote value (-1, 0, +1) |
| `tags` | `block,spam` | vote tags, separated by commas |
| `creator` | `03f5bcfadd87e625bf62900a7d1ed673ce74034dbfc5d5c624cedd4612a8dc6d1c` | public key of the creator of the thread vote |

Thus, returning:

```json
{
    "okay": true,
    "data": {
        "hash": "9f3c37be48bfaaa1c0d6a0231110d46b79961a0fb85a889c23efc8005b871e36",
        "raw": "{\"type\":\"5,user_vote\",\"ts\":1512533086034755384,\"of_board\":\"032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b\",\"of_user\":\"0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718\",\"value\":-1,\"tags\":[\"block\",\"spam\"],\"creator\":\"03f5bcfadd87e625bf62900a7d1ed673ce74034dbfc5d5c624cedd4612a8dc6d1c\"}"
    }
}
```

User 3 finalizes the submission via `/api/submission/finalize`:

| Key | Value | Description |
| --- | --- | --- |
| `hash` | `9f3c37be48bfaaa1c0d6a0231110d46b79961a0fb85a889c23efc8005b871e36` | hash of content that needs submission finalization |
| `sig` | `642c47dba5bb06069d19f6ef2fe7041c849788dcc7580d46fadb130131110cb978b0bebe4d70522d7cac83f646b2b929bd2fb285efba96c8cc98af41dab4e6c000` |

Thus, returning:

```json
{
    "okay": true,
    "data": {
        "user_public_key": "03f5bcfadd87e625bf62900a7d1ed673ce74034dbfc5d5c624cedd4612a8dc6d1c",
        "profile": {
            "trusted_count": 0,
            "trusted": [],
            "marked_as_spam_count": 1,
            "marked_as_spam": [
                "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718"
            ],
            "blocked_count": 1,
            "blocked": [
                "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718"
            ],
            "trusted_by_count": 0,
            "trusted_by": [],
            "marked_as_spam_by_count": 0,
            "marked_as_spam_by": [],
            "blocked_by_count": 0,
            "blocked_by": []
        }
    }
}
```

Note that both `"marked_as_spam"` and `"blocked"` now have entries in User 3's profile.

**Extra:**

We can also now check out the default user's profile via the `/api/get_user_profile` endpoint:

| Key | Value | Description |
| --- | --- | --- |
| `board_public_key` | `032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b` | Public key of board to extract user's follow page from. |
| `user_public_key` | `0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718` | User's public key to extract follow page from. |

Which returns the following result:

```json
{
    "okay": true,
    "data": {
        "user_public_key": "0254020da01e33cbaf2ff01e7cf28de4bb6cea43b357153fad3a50a0e7dd728718",
        "profile": {
            "trusted_count": 0,
            "trusted": [],
            "marked_as_spam_count": 0,
            "marked_as_spam": [],
            "blocked_count": 0,
            "blocked": [],
            "trusted_by_count": 1,
            "trusted_by": [
                "02f46d2461e2c3aba0585efb5b2ddb8acb34f38a56865f8a2a3f10272e6de257c1"
            ],
            "marked_as_spam_by_count": 1,
            "marked_as_spam_by": [
                "03f5bcfadd87e625bf62900a7d1ed673ce74034dbfc5d5c624cedd4612a8dc6d1c"
            ],
            "blocked_by_count": 2,
            "blocked_by": [
                "0284da18e80d5ec08cf54ed9c86bbbce6bbd2838b8c700a373a5886e4de44ce895",
                "03f5bcfadd87e625bf62900a7d1ed673ce74034dbfc5d5c624cedd4612a8dc6d1c"
            ]
        }
    }
}
```