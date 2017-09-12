package factory

import (
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
)

func init() {
	ops[OP_REG] = &sync.Pool{
		New: func() interface{} {
			return new(reg)
		},
	}
	resps[OP_REG] = &sync.Pool{
		New: func() interface{} {
			return new(regResp)
		},
	}
}

type reg struct {
}

func (reg *reg) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	if conn.IsKeySet() {
		conn.GetContextLogger().Infof("reg %s already", conn.key.Hex())
		return
	}
	key, _ := cipher.GenerateKeyPair()
	conn.SetKey(key)
	conn.SetContextLogger(conn.GetContextLogger().WithField("pubkey", key.Hex()))
	f.register(key, conn)
	r = &regResp{PubKey: key}
	return
}

type regResp struct {
	PubKey cipher.PubKey
}

func (resp *regResp) Run(conn *Connection) (err error) {
	conn.SetKey(resp.PubKey)
	conn.factory.register(resp.PubKey, conn)
	conn.SetContextLogger(conn.GetContextLogger().WithField("pubkey", resp.PubKey.Hex()))
	return
}
