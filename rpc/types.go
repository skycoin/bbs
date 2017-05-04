package rpc

import (
	"github.com/evanlinjin/bbs/cxo"
	"github.com/skycoin/skycoin/src/cipher"
)

// NewPostReq represents a NewPost request.
type NewPostReq struct {
	Board  cipher.PubKey // Board Identifier.
	Thread []byte        // Thread Identifier.
	Post   cxo.Post      // Post to inject.
}

// NewThreadReq represents  a NewThread request.
type NewThreadReq struct {
	Board  cipher.PubKey // Board Identifier.
	Thread cxo.Thread    // Thread to inject.
}
