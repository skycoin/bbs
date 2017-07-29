package object

import (
	"github.com/skycoin/bbs/src/misc/verify"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type BoardPage struct {
	R           cipher.SHA256        `json:"-" enc:"-"`
	Board       skyobject.Reference  `skyobject:"schema=Board"`
	ThreadPages skyobject.References `skyobject:"schema=ThreadPage"`
}

type Board struct {
	R                   cipher.SHA256  `json:"-" enc:"-"`
	Name                string         `json:"name"`
	Desc                string         `json:"description"`
	Created             int64          `json:"created"`
	SubmissionAddresses []string       `json:"submission_addresses"`
	ExternalRoots       []ExternalRoot `json:"-"`
	Meta                []byte         `json:"-"`
}

type ExternalRoot struct {
	R         cipher.SHA256 `json:"-" enc:"-"`
	ID        string        `json:"id"`
	PublicKey cipher.PubKey `json:"-"`
}

type ThreadPage struct {
	R      cipher.SHA256        `json:"-" enc:"-"`
	Thread skyobject.Reference  `skyobject:"schema=Thread"`
	Posts  skyobject.References `skyobject:"schema=Post"`
}

type Thread struct {
	R cipher.SHA256 `json:"-" enc:"-"`
	Post
}

type Post struct {
	R       cipher.SHA256 `json:"-" enc:"-"`
	Title   string        `json:"title"`
	Body    string        `json:"body"`
	Created int64         `json:"created"`
	User    cipher.PubKey `json:"-" verify:"pk"`
	Sig     cipher.Sig    `json:"-" verify:"sig"`
	Meta    []byte        `json:"-"`
}

// Verify verifies the post.
func (p Post) Verify() error { return verify.Check(&p) }

// Vote represents a post by a user.
type Vote struct {
	R       cipher.SHA256 `json:"-" enc:"-"`
	User    cipher.PubKey `json:"-" verify:"pk"` // User who voted.
	Mode    int8          `json:"-"`             // +1 is up, -1 is down.
	Tag     []byte        `json:"-"`             // What's this?
	Created int64         `json:"created"`
	Sig     cipher.Sig    `json:"-" verify:"sig"` // Signature.
}

func (v Vote) Verify() error { return verify.Check(&v) }

type ThreadVotesPage struct {
	R       cipher.SHA256 `json:"-" enc:"-"`
	Store   []VotesPage
	Deleted []cipher.SHA256
}

type PostVotesPage struct {
	R       cipher.SHA256 `json:"-" enc:"-"`
	Store   []VotesPage
	Deleted []cipher.SHA256
}

type VotesPage struct {
	Ref   cipher.SHA256
	Votes skyobject.References `skyobject:"schema=Vote"`
}

type User struct {
	R         cipher.SHA256 `json:"-" enc:"-"`
	Alias     string        `json:"alias"`
	PublicKey cipher.PubKey `json:"-"`
	SecretKey cipher.SecKey `json:"-"`
}

type Subscription struct {
	R      cipher.SHA256 `json:"-" enc:"-"`
	PubKey cipher.PubKey `json:"pk"`
	SecKey cipher.SecKey `json:"sk,omitempty"`
}
