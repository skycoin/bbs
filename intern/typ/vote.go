package typ

import (
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type ThreadVoteContainer struct {
	Threads []struct {
		Thread skyobject.Reference  `skyobject:"schema=Thread"`
		Votes  skyobject.References `skyobject:"schema=Vote"`
	}
}

type PostVoteContainer struct {
	Posts []struct {
		Post  skyobject.Reference  `skyobject:"schema=Post"`
		Votes skyobject.References `skyobject:"schema=Vote"`
	}
}

type Vote struct {
	User cipher.PubKey // User who voted.
	Mode int8          // +1 is up, -1 is down.
	Tag  []byte        // What's this?
	Sig  cipher.Sig    // Signature.
}
