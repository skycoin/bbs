package op

import (
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type Send struct {
	PublicKey string
	Msg       string
}

func init() {
	msg.OP_POOL[msg.OP_SEND] = &sync.Pool{
		New: func() interface{} {
			return new(Send)
		},
	}
}

func (s *Send) Execute(c msg.OPer) error {
	key, err := cipher.PubKeyFromHex(s.PublicKey)
	if err != nil {
		return err
	}
	return c.GetConnection().Send(key, []byte(s.Msg))
}
