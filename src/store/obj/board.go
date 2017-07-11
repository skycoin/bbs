package obj

import (
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type BoardPage struct {
	Board       skyobject.Reference  `skyobject:"schema=Board"`
	ThreadPages skyobject.References `skyobject:"schema=ThreadPage"`
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
