package rpc

import (
	"github.com/evanlinjin/bbs/intern/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type ReqNewPost struct {
	BoardPubKey cipher.PubKey       `json:"board_public_key,string"`
	ThreadRef   skyobject.Reference `json:"thread_reference,string"`
	Post        *typ.Post           `json:"post"`
}

type ReqNewThread struct {
	BoardPubKey cipher.PubKey `json:"board_public_key,string"`
	Creator     cipher.PubKey `json:"creator,string"`
	Signature   cipher.Sig    `json:"signature,string"`
	Thread      *typ.Thread   `json:"thread"`
}
