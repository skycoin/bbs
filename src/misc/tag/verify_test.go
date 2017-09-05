package tag

import (
	"encoding/json"
	"github.com/skycoin/skycoin/src/cipher"
	"testing"
	"time"
)

type Post struct {
	Title   string        `json:"title"`
	Body    string        `json:"body"`
	Created int64         `json:"created"`
	User    cipher.PubKey `json:"user" verify:"pk"`
	Sig     cipher.Sig    `json:"sig" verify:"sig"`
}

func print(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func TestSign(t *testing.T) {
	post := &Post{
		Title:   "Test title",
		Body:    "Test body",
		Created: time.Now().UnixNano(),
	}
	t.Log("Post:", *post)

	pk, sk := cipher.GenerateKeyPair()
	Sign(post, pk, sk)
	t.Log("Post:", *post)

	tempPost := *post
	if e := Verify(&tempPost); e != nil {
		t.Error(e)
	}
	t.Log("Post:", *post)
}
