package obj

import (
	"github.com/skycoin/bbs/src/misc/verify"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type BoardPage struct {
	Board       skyobject.Reference  `skyobject:"schema=Board"`
	ThreadPages skyobject.References `skyobject:"schema=ThreadPage"`
	Deleted     cipher.SHA256
}

type Board struct {
	Name                string         `json:"name"`
	Desc                string         `json:"description"`
	Created             int64          `json:"created"`
	SubmissionAddresses []string       `json:"submission_addresses"`
	ExternalRoots       []ExternalRoot `json:"-"`
	Meta                []byte         `json:"-"`
}

type ExternalRoot struct {
	ID        string        `json:"id"`
	PublicKey cipher.PubKey `json:"-"`
}

type ThreadPage struct {
	Thread  skyobject.Reference  `skyobject:"schema=Thread"`
	Posts   skyobject.References `skyobject:"schema=Post"`
	Deleted cipher.SHA256
}

type Thread struct {
	Post
	MasterBoardRef skyobject.Reference `json:"-" skyobject:"schema=Board"`
}

type Post struct {
	Title   string        `json:"title"`
	Body    string        `json:"body"`
	Created int64         `json:"created"`
	User    cipher.PubKey `json:"-" verify:"pk"`
	Sig     cipher.Sig    `json:"-" verify:"sig"`
	Meta    []byte        `json:"-"`
}

// Verify verifies the post.
func (p Post) Verify() error { return verify.Check(&p) }

type ThreadVotesPage struct {
	Store []VotesPage
}

type PostVotesPage struct {
	Store []VotesPage
}

type VotesPage struct {
	Ref   cipher.SHA256
	Votes skyobject.References `skyobject:"schema=Vote"`
}

type User struct {
	Alias     string        `json:"alias"`
	PublicKey cipher.PubKey `json:"-"`
	SecretKey cipher.SecKey `json:"-"`
}

type Subscription struct {
	PubKey      cipher.PubKey `json:"pk"`
	SecKey      cipher.SecKey `json:"sk,omitempty"`
	Connections []string      `json:"conns,omitempty"`
}
