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
  "header": {
    "type": "5,board"
  },
  "body": {
    "name": "Test Board 2",
    "body": "This is a test board.",
    "created": 1507011466251404010,
    "submission_keys": [
      "021f8ef6570c719e408bc450ec07f2b946b3c8e7e8acd1eef529c3a941416bad4b"
    ],
    "tags": ["test", "fun"]
  }
}
```

**For a `"thread"`:**

```json
{
  "header": {
    "type": "5,thread",
    "hash": "37215d922cc682ef81763a40cd51edbc3c1c9e9b9c7bc4c52e23b8115fc30c97",
    "pk": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7",
    "sig": "d4572654617b12e186f4604b60a52e57957dab90ca1bcdb9da60b6dfa50788fb12e3554445b0d1d2f01ecb8419cd80fe407f799e27ba247db30f8142815b254400"
  },
  "body": {
    "of_board": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
    "name": "Test Thread",
    "body": "This is a test thread.",
    "created": 1507011466251404010,
    "creator": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"
  }
}
```

**For a `"post"`:**

```json
{
  "header": {
    "type": "5,post",
    "hash": "22947b6018ea170ec4e19a2df3221989fd0531639d86a2b60fb4e9b4a0771383",
    "pk": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7",
    "sig": "8dd49b8c904ac0acd187549c46dc1adb55489aab7df8644a6d96f27306d37c2447e6af4bcaa9a991315a45e2cf363a1b04022b6c3c3fb0ec9a453f9025b0092200"
  },
  "body": {
    "of_board": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
    "of_thread": "37215d922cc682ef81763a40cd51edbc3c1c9e9b9c7bc4c52e23b8115fc30c97",
    "name": "Test Post",
    "body": "This is a test post.",
    "created": 1507011466251404010,
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
  "of_board": "02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1",
  "name": "Test Thread",
  "body": "This is a test thread.",
  "created": 1507011466251404010,
  "creator": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7"
}
```

Notice that the `creator` field contains the user's public key, and the `created` timestamp also needs to be generated.

It's best to represent the above in as little whitespace as possible. And the order of the fields do not matter.

```text
{"name":"Test Thread","body":"This is a test thread.","created":1507011466251404010,"creator":"035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7","of_board":"02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"}
```

Now let's hash it!

```json
{
  "hash": "37215d922cc682ef81763a40cd51edbc3c1c9e9b9c7bc4c52e23b8115fc30c97"
}
```

And sign it with the secret key `69eb49ceac75c3e0993395ecf25730578499f8091a935a59452f9fef7115dd4d` of the user.

```json
{
  "sig": "bae9535a70979a380748c052ae56657229257d6d4971e6a8da58887eee8eb87c40f349820bcc9fa5a5705229de30f6dc933a547d235a8677e821044827fa181001"
}
```

We can then submit the data using the endpoint represented as "New Submission" in the provided Postman collection.

| Key | Value |
| --- | --- |
| `type` | `5,post` |
| `body` | `{"name":"Test Thread","body":"This is a test thread.","created":1507011466251404010,"creator":"035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7","of_board":"02e5be89fa161bf6b0bc64ec9ec7fe27311fbb78949c3ef9739d4c73a84920d6e1"}` |
| `sig` | `bae9535a70979a380748c052ae56657229257d6d4971e6a8da58887eee8eb87c40f349820bcc9fa5a5705229de30f6dc933a547d235a8677e821044827fa181001` |