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

type ReqVoteUser struct {
	BoardPubKey cipher.PubKey `json:"board_public_key,string"`
	UserPubKey  cipher.PubKey `json:"user_public_key,string"`
	Vote        *typ.Vote     `json:"vote"`
}
