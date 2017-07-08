package obj

import "github.com/skycoin/skycoin/src/cipher"

type Board struct {
	Name            string        `json:"name"`
	Desc            string        `json:"description"`
	Created         int64         `json:"created"`
	Ref             cipher.PubKey `json:"-"`
	ContentVotesRef cipher.PubKey `json:"-"`
	UserVotesRef    cipher.PubKey `json:"-"`
	Meta            []byte        `json:"-"`
}
