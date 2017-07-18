package state

import "github.com/skycoin/skycoin/src/cipher"

// Compiler compiles board states.
type Compiler struct {
	boards map[cipher.PubKey]struct{}
}
