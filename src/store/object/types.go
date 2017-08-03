package object

import (
	"github.com/skycoin/bbs/src/misc/tag"
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
	Name                string         `json:"name" transfer:"heading"`
	Desc                string         `json:"description" transfer:"body"`
	Created             int64          `json:"created" transfer:"time"`
	SubmissionAddresses []string       `json:"submission_addresses" transfer:"subAddrs"`
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
	Title   string        `json:"title" transfer:"heading"`
	Body    string        `json:"body" transfer:"body"`
	Created int64         `json:"created" transfer:"time"`
	User    cipher.PubKey `json:"-" verify:"pk" transfer:"upk"`
	Sig     cipher.Sig    `json:"-" verify:"sig"`
	Meta    []byte        `json:"-"`
}

// Verify verifies the post.
func (p Post) Verify() error { return tag.Verify(&p) }

// Vote represents a post by a user.
type Vote struct {
	R       cipher.SHA256 `json:"-" enc:"-"`
	User    cipher.PubKey `json:"-" verify:"pk" transfer:"upk"` // User who voted.
	Mode    int8          `json:"-" transfer:"mode"`            // +1 is up, -1 is down.
	Tag     []byte        `json:"-" transfer:"tag"`             // What's this?
	Created int64         `json:"created" transfer:"time"`
	Sig     cipher.Sig    `json:"-" verify:"sig"` // Signature.
}

func (v Vote) Verify() error { return tag.Verify(&v) }

type ThreadVotesPages struct {
	R           cipher.SHA256 `json:"-" enc:"-"`
	StoreHash   cipher.SHA256 `enc:"-"`
	DeletedHash cipher.SHA256 `enc:"-"`
	Store       []VotesPage
	Deleted     []cipher.SHA256
}

type PostVotesPages struct {
	R           cipher.SHA256 `json:"-" enc:"-"`
	StoreHash   cipher.SHA256 `enc:"-"`
	DeletedHash cipher.SHA256 `enc:"-"`
	Store       []VotesPage
	Deleted     []cipher.SHA256
}

type VotesPage struct {
	Ref   cipher.SHA256
	Votes skyobject.References `skyobject:"schema=Vote"`
}

type User struct {
	R         cipher.SHA256 `json:"-" enc:"-"`
	Alias     string        `json:"alias" transfer:"alias"`
	PublicKey cipher.PubKey `json:"-" transfer:"upk"`
	SecretKey cipher.SecKey `json:"-" transfer:"usk"`
}

type Subscription struct {
	R      cipher.SHA256 `json:"-" enc:"-"`
	PubKey cipher.PubKey `json:"pk" transfer:"bpk"`
	SecKey cipher.SecKey `json:"sk,omitempty" transfer:"bsk"`
}
