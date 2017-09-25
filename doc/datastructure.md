# Data Structure

## Post

### Current Implementation

**Submitting to API:**

Format: `form-data`

| Key | Value (Example) | Description |
| --- | --- | --- |
| `board_public_key` | `035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7` | The public key of the board in which to submit the post (hex string representation). |
| `thread_ref` | `fcef784bf16a1c62206b6230f5c1d88137bd3641d63087e9fe80cedf4e536d9f` | The hash of the thread in which to submit the post (hex string representation). |
| `post_ref` | `03fd08e02a700e06fcf6772d06261f88809564d3421bffb5a81119bc5dbfe5aca4` | ***(Optional)*** The hash of the post that this post is to be a reply to. Only needed if this post is a reply.
| `name` | `Test Post` | Name of the post to create. |
| `body` | `This is a test post.` | Body of the post to create. |

**Format in CXO:**

Format: Serialized golang object.

```go
type Post struct {
	R   cipher.SHA256 `enc:"-" json:"-"` // Ignored field. Used to store hash of post.
	Raw []byte                           // Serialized raw JSON data of post.
	Sig cipher.Sig `verify:"sig"`        // Signature of post.
}
```

The following object is what is serialized for the `Raw` field of `Post`.

```go
type PostData struct {
	OfBoard  string       `json:"of_board"`          // Public key of board in which post belongs.
	OfThread string       `json:"of_thread"`         // SHA256 of thread in which post belongs.
	OfPost   string       `json:"of_post,omitempty"` // SHA256 of post this post is replying to (optional).
	Name     string       `json:"name"`              // Name of post.
	Body     string       `json:"body"`              // Body of post.
	Images   []*ImageData `json:"images,omitempty"`  // Images of post (optional).
	Created  int64        `json:"created"`           // Time of creation (unix time in ns).
	Creator  string       `json:"creator"`           // Public key of creator (in hex).
}
```

The following is the representation of `ImageData`.

```go
type ImageData struct {
	Name   string       `json:"name"`
	Hash   string       `json:"hash"`
	URL    string       `json:"url,omitempty"`
	Size   int          `json:"size,omitempty"`
	Height int          `json:"height,omitempty"`
	Width  int          `json:"width,omitempty"`
	Thumbs []*ImageData `json:"thumbs,omitempty"`
}
```

**Format in compiled view:**

```go
type PostRepView struct {
	Seq     int             `json:"seq"`              // Sequence of post in the thread.
	Ref     string          `json:"ref"`              // Hash of the post.
	Type    string          `json:"type"`             // Type of post (image/text).
	Name    string          `json:"name"`             // Name of the post.
	Body    string          `json:"body"`             // Body of the post.
	Images  []*r0.ImageData `json:"images,omitempty"` // Images of the post.
	Created int64           `json:"created"`          // Timestamp of when the post is created.
	Creator string          `json:"creator"`          // Hex string representation of the creator's public key.
	Votes   *VoteRepView    `json:"votes,omitempty"`  // A representation of votes that can easily be shown via UI.
}
```

### Future Implementation

In the future, what is submitted, stored and shown will be more similar.

**Submitting to API:**

Format: `form-data`



| Key | Value (Example) | Description |
| --- | --- | --- |
| `post` | JSON format shown below. | The JSON representation of a post to submit. |
| `sig` | `fcef784bf16a1c62206b6230f5c1d88137bd3641d63087e9fe80cedf4e536d9f` | The signature of the above post, signed with the creator's public key. |

Post JSON format:

```json
{
  "of_board": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7",
  "of_thread": "fcef784bf16a1c62206b6230f5c1d88137bd3641d63087e9fe80cedf4e536d9f",
  "name": "Test Post",
  "body": "This is a test post.",
  "created": 1506298296996706914,
  "creator": "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b"
}
```

Format in CXO will be the same.

**Format in compiled view:**

```json
{
  "post": {
    "of_board": "035a630a621aa3483f87cb288438982d7ba8524302ed6f293f667e6d8c9fa369a7",
      "of_thread": "fcef784bf16a1c62206b6230f5c1d88137bd3641d63087e9fe80cedf4e536d9f",
      "name": "Test Post",
      "body": "This is a test post.",
      "created": 1506298296996706914,
      "creator": "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b"
  },
  "sig": "02c9d0d1faca3c852c307b4391af5f353e63a296cded08c1a819f03b7ae768530b",
  "votes": {}
}
```

### Parsing code (Current)

https://github.com/skycoin/bbs/blob/master/src/store/object/revisions/r0/cxo_types_content.go#L93

https://github.com/skycoin/bbs/blob/master/src/store/object/revisions/r0/cxo_types_content_data.go#L59