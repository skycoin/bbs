package typ

import (
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

// Post represents a post.
type Post struct {
	ID        cipher.PubKey `json:"id"`
	Signature []byte        `json:"signature"`
	Name      string        `json:"name"`
	Body      string        `json:"body"`
	Publisher cipher.PubKey `json:"publisher"`
	Created   int64         `json:"created"`
}

// NewPost creates a new post.
func NewPost(name, body string, publisher cipher.PubKey) *Post {
	now := time.Now().UnixNano()
	post := Post{
		ID:        MakeTimeStampedRandomID(128),
		Name:      name,
		Body:      body,
		Publisher: publisher,
		Created:   now,
	}
	return &post
}
