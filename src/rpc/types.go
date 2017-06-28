package rpc

import (
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type ReqNewPost struct {
	BoardPubKey cipher.PubKey       `json:"board_public_key,string"`
	ThreadRef   skyobject.Reference `json:"thread_ref,string"`
	Post        *typ.Post           `json:"post"`
}

type ReqNewThread struct {
	BoardPubKey cipher.PubKey `json:"board_public_key,string"`
	Creator     cipher.PubKey `json:"creator,string"`
	Signature   cipher.Sig    `json:"signature,string"`
	Thread      *typ.Thread   `json:"thread"`
}

type ReqVotePost struct {
	BoardPubKey cipher.PubKey       `json:"board_public_key,string"`
	PostRef     skyobject.Reference `json:"post_ref,string"`
	Vote        *typ.Vote           `json:"vote"`
}

type ReqVoteThread struct {
	BoardPubKey cipher.PubKey       `json:"board_public_key,string"`
	ThreadRef   skyobject.Reference `json:"thread_ref,string"`
	Vote        *typ.Vote           `json:"vote"`
}
