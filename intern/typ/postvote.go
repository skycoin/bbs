package typ

import "github.com/skycoin/cxo/skyobject"

// PostVotePage is an element of PostVoteContainer.
type PostVotePage struct {
	Post  skyobject.Reference  `skyobject:"schema=Post"`
	Votes skyobject.References `skyobject:"schema=Vote"`
}

// PostVoteContainer contains the votes of posts.
type PostVoteContainer struct {
	Posts []PostVotePage
}
