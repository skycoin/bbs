package typ

import (
	"github.com/skycoin/skycoin/src/cipher"
	"testing"
)

func TestPost_Verify(t *testing.T) {
	// Get user public and secret keys.
	upk, usk := cipher.GenerateDeterministicKeyPair([]byte("a"))
	// Make random post.
	post := &Post{
		Title: "A Post",
		Body:  "This is just a post.",
	}
	// Sign post.
	if e := post.Sign(upk, usk); e != nil {
		t.Error(e)
	}
	// Touch post.
	post.Touch()
	// Verify post.
	if e := post.Verify(); e != nil {
		t.Error(e)
	}
}
