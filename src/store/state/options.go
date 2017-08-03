package state

import (
	"github.com/skycoin/bbs/src/store/state/states"
	"github.com/skycoin/bbs/src/store/state/states/v1"
)

type Option func(c *Compiler) error

func SetV1State() Option {
	return func(c *Compiler) error {
		c.newBState = states.NewState(v1.NewBoardState)
		return nil
	}
}
