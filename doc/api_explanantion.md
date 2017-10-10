# v5.0 API Explanation

It's absolutely retarded to make people run full BBS nodes to interact with boards. Let's make it so that BBS nodes host publicly available websites where the end-user can view boards and submit content. In order to do this, user session management needs to happen client-side. Woops.

The website which is publicly hosted by BBS nodes should only have the following functionality:

* Display content (list boards, threads, posts, votes, etc).
* Submit content (threads, posts and votes).

## Postman Collection

Yes, the Postman collection is gracefully updated.

The stuff under "Admin" folder contains calls which should **not** be publicly available (however, they are publicly available now for convenience sake - don't use them as they will disappear). The functionally here will be implemented in a command-line interface.

The stuff under "Tools" should all be implemented client side. But you can use these endpoints for convenience sake for right now. These include:

* Seed generation (New Seed).
* Generating deterministic key pairs (New Key Pair).
* SHA256 Hashing (Hash).
* Signing data (Sign).

However, these should really be done via the `skycoin-cipher-web` library: https://github.com/skycoin/skycoin-cipher-web

The endpoints that are to be publicly available are as follows.

* For displaying content: "Get Board", "Get Boards", "Get Board Page", "Get Thread Page", "Get Follow Page", "Get Discovered Boards" (Not working). These endpoints are the same in structure as the `v0.3` endpoints (despite the change in URI). However, 

* For submitting content: "New Submission".

These should be self-explanatory. However, the structure for `"board"`, `"thread"` and `"post"` have changed (further discussed below).

## JSON Structure Changes

The data structures for `"board"`, `"thread"` and `"post"` are now represented with `"header"` and `"body"` fields.

**For a `board`:**

```json
{
  "public_key": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
  "header": {},
  "body": {
    "type": "5,board",
    "ts": 1507604859033404987,
    "name": "Test Board",
    "body": "This is a test board."
  }
}
```

**For a `"thread"`:**

```json
{
  "header": {
    "hash": "a3e3850c1dd3933ec44b9f93c42f6431d4d76933a8fb7e73b2e6d3706f8ee63b",
    "sig": "086cdca50071c75a5dbf3d82c202eb7b8edeef29f98a5ebd410145f0e79cd3a574394bcf4852b04f08d47587fa10f5abca6b3263d0b9ce106e45e53e0826343700"
  },
  "body": {
    "type": "5,thread",
    "ts": 123435,
    "of_board": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
    "name": "Test Thread",
    "body": "This is a test thread.",
    "creator": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"
  }
}
```

**For a `"post"`:**

```json
{
  "header": {
    "hash": "63b347773e1f849066007760e49ff3fc96eec297271160bfd08dc636d5b7f9e0",
    "sig": "65677b4deabee53d00ca00de82d2e9258a39e908a2b31979a0d936b2bfd4fcde65bcf9b7364e9338b3ac32453cd0e617a763a83e54a9d0665ae04ebbdccf79ba00"
  },
  "body": {
    "type": "5,post",
    "ts": 123435,
    "of_board": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
    "of_thread": "a3e3850c1dd3933ec44b9f93c42f6431d4d76933a8fb7e73b2e6d3706f8ee63b",
    "name": "Test Post",
    "body": "This is a test post.",
    "creator": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"
  }
}
```

### Content Header

The content `"header"` is the same structure for boards, threads and posts. Here are the fields:

* `"type"` - the type of content of format `"{version},{type}"`. Current valid values are as follows:
    * `5,thread` - Thread (api v5).
    * `5,post` - Post (api v5).
    * `5,thread_vote` - A thread vote (api v5).
    * `5,post_vote` - A post vote (api v5).
    * `5,user_vote` - A user vote (api v5).
    
* `"hash"` - This is the hash of the content being submitted. It is also used to reference the specified content. It's the hex representation of a SHA256 hash.

* `"pk"` - The public key used to verify the content. This is also the public key of the user that generated the content. It's the hex representation of the public key.

* `"sig"` - This is the signature of the submitted data, signed with the creator's private key. This can be verified with the creator's public key. It's the hex representation of the signature.

For the `board` type, as boards are not user-generated, there is no need for verification, hence the `header` for `board` types only have the `type` field.

### Content Body

This varies with the content type. Examples as above.

### Content Submission

So this is the difficult part.

#### Example 1 - Submitting a thread

First, we need a private, public key pair that represents a user. This is generated with a seed.

Given a seed.

```json
{
  "seed": "c"
}
```

We can generate a key pair as follows.

```json
{
  "public_key": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7",
  "secret_key": "69eb49ceac75c3e0993395ecf25730578499f8091a935a59452f9fef7115dd4d"
}
```

Let's add a thread to board of public key `02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1`. This is the data body of such a thread.

```json
{
  "type": "5,thread",
  "of_board": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
  "ts": 123435,
  "name": "Test Thread",
  "body": "This is a test thread.",
  "creator": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"
}
```

Notice that the `creator` field contains the user's public key, and the `created` timestamp also needs to be generated.

It's best to represent the above in as little whitespace as possible. And the order of the fields do not matter.

```text
{"type":"5,thread","of_board":"02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1","ts":123435,"name":"Test Thread","body":"This is a test thread.","creator":"035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"}
```

Now let's hash it!

```json
{
  "hash": "a3e3850c1dd3933ec44b9f93c42f6431d4d76933a8fb7e73b2e6d3706f8ee63b"
}
```

And sign it with the secret key `69eb49ceac75c3e0993395ecf25730578499f8091a935a59452f9fef7115dd4d` of the user.

```json
{
  "sig": "086cdca50071c75a5dbf3d82c202eb7b8edeef29f98a5ebd410145f0e79cd3a574394bcf4852b04f08d47587fa10f5abca6b3263d0b9ce106e45e53e0826343700"
}
```

We can then submit the data using the endpoint represented as "New Submission" in the provided Postman collection.

| Key | Value |
| --- | --- |
| `body` | `{"type":"5,thread","of_board":"02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1","ts":123435,"name":"Test Thread","body":"This is a test thread.","creator":"035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"}` |
| `sig` | `086cdca50071c75a5dbf3d82c202eb7b8edeef29f98a5ebd410145f0e79cd3a574394bcf4852b04f08d47587fa10f5abca6b3263d0b9ce106e45e53e0826343700` |

#### Example 2 - Submitting a thread vote

We will use the key pair generated with seed `"c"`.

```json
{
  "public_key": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7",
  "secret_key": "69eb49ceac75c3e0993395ecf25730578499f8091a935a59452f9fef7115dd4d"
}
```

And again, we are submitting to board of public key `02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1`. We will use the above submitted thread to cast vote on. The hash of the thread is `a3e3850c1dd3933ec44b9f93c42f6431d4d76933a8fb7e73b2e6d3706f8ee63b`.

Hence, the thread-vote body will be as follows.

```json
{
  "type": "5,thread_vote",
  "ts": 12345,
  "of_board": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
  "of_thread": "a3e3850c1dd3933ec44b9f93c42f6431d4d76933a8fb7e73b2e6d3706f8ee63b",
  "value": -1,
  "tag": "spam",
  "creator": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"
}
```

The thread-vote body represented with minimum whitespace is as follows.

```text
{"type":"5,thread_vote","ts":12345,"of_board":"02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1","of_thread":"a3e3850c1dd3933ec44b9f93c42f6431d4d76933a8fb7e73b2e6d3706f8ee63b","value":-1,"tag":"spam","creator":"035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"}
```

And we can obtain a signature using the user's private key as follows.

```text
eaed0a4ec843b0e871645b9be1863d56b79fa372d60498f2900ac16888f9e7de5fe503d87173e27f8c53250e9858c1429b4143af3b9aa715e0448232994edaf601
```

#### Example 3 - Submitting a post

Body (`body`):

```text
{"type":"5,post","of_board":"02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1","of_thread":"a3e3850c1dd3933ec44b9f93c42f6431d4d76933a8fb7e73b2e6d3706f8ee63b","ts":123435,"name":"Test Post","body":"This is a test post.","creator":"035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"}
```

Signature (`sig`):

```text
65677b4deabee53d00ca00de82d2e9258a39e908a2b31979a0d936b2bfd4fcde65bcf9b7364e9338b3ac32453cd0e617a763a83e54a9d0665ae04ebbdccf79ba00
```