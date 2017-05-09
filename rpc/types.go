package rpc

import (
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// NewPostReq represents a NewPost request.
type NewPostReq struct {
	PK   cipher.PubKey       // PK Identifier.
	Hash skyobject.Reference // Hash Identifier.
	Post *typ.Post           // Post to inject.
}

// NewThreadReq represents  a NewThread request.
type NewThreadReq struct {
	PK     cipher.PubKey // PK Identifier.
	Thread *typ.Thread   // Hash to inject.
}
