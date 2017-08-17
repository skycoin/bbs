package state

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type Compiler struct {
	boards map[cipher.PubKey]*BoardInstance
}
