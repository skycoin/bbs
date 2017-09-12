package op

import (
	"sync"

	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/net/skycoin-messenger/msg"
)

type Reg struct {
	Address string
}

func init() {
	msg.OP_POOL[msg.OP_REG] = &sync.Pool{
		New: func() interface{} {
			return new(Reg)
		},
	}
}

func (r *Reg) Execute(c msg.OPer) error {
	f := factory.NewMessengerFactory()
	conn, err := f.Connect(r.Address)
	if err != nil {
		return err
	}
	c.SetConnection(conn)
	go c.PushLoop(conn)
	return nil
}
