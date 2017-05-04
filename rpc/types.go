package rpc

import (
	"github.com/evanlinjin/bbs/datastore"
	"github.com/skycoin/skycoin/src/cipher"
)

// NewPostReq represents a NewPost request.
type NewPostReq struct {
	Board  cipher.PubKey  // Board Identifier.
	Thread []byte         // Thread Identifier.
	Post   datastore.Post // Post to inject.
}

// NewThreadReq represents  a NewThread request.
type NewThreadReq struct {
	Board cipher.PubKey // Board Identifier.
	Thread datastore.Thread // Thread to inject.
}