package object

import (
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

/*
	<<< ROOT CHILDREN >>>
*/

type BoardPage struct {
	Board skyobject.Ref `skyobject:"schema=bbs.Board"`
	Meta  []byte
}

type ContentPage struct {
	Threads skyobject.Refs `skyobject:"schema=bbs.Content"`
	Content skyobject.Refs `skyobject:"schema=bbs.Content"`
	Deleted skyobject.Refs `skyobject:"schema=bbs.Content"`
}

type ActivityPage struct {
	Users skyobject.Refs `skyobject:"schema=bbs.UserActivity"`
}

type UserActivity struct {
	UserVotes        skyobject.Refs `skyobject:"schema=bbs.Vote"`
	ContentVotes     skyobject.Refs `skyobject:"schema=bbs.Vote"`
	ContentCreations skyobject.Refs `skyobject:"schema=bbs.Content"`
	ContentDeletions skyobject.Refs `skyobject:"schema=bbs.Content"`
}

type Board struct {
	Name     string
	Desc     string
	SubAddrs []string
	Created  int64
}

type Vote struct {
	OfUser   cipher.PubKey
	OfThread cipher.SHA256
	OfPost   cipher.SHA256

	Mode int8
	Tag  []byte

	Created int64         `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func (v Vote) Verify() error { return tag.Verify(&v) }

/*
	<<< USER >>>
*/

type User struct {
	Alias  string        `json:"alias" trans:"alias"`
	PubKey cipher.PubKey `json:"-" trans:"upk"`
	SecKey cipher.SecKey `json:"-" trans:"usk"`
}

type UserView struct {
	User
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key,omitempty"`
}

/*
	<<< CONNECTION >>>
*/

type Connection struct {
	Address string `json:"address"`
	State   string `json:"state"`
}
