package typ

import (
	"errors"
	"github.com/skycoin/cxo/skyobject"
)

// BoardContainer represents the first branch object from root, in cxo.
type BoardContainer struct {
	Board       skyobject.Reference  `skyobject:"schema=Board"`
	Threads     skyobject.References `skyobject:"schema=Thread"`
	ThreadPages skyobject.References `skyobject:"schema=ThreadPage"`
	Seq         uint64
}

func NewBoardContainer(bRef skyobject.Reference) *BoardContainer {
	return &BoardContainer{
		Board: bRef,
		Seq:   0,
	}
}

func (c *BoardContainer) AddThread(ref skyobject.Reference) error {
	for _, r := range c.Threads {
		if r == ref {
			return errors.New("thread already exists")
		}
	}
	c.Threads = append(c.Threads, ref)
	return nil
}

func (c *BoardContainer) AddThreadPage(ref skyobject.Reference) error {
	for _, r := range c.ThreadPages {
		if r == ref {
			return errors.New("thread already exists")
		}
	}
	c.ThreadPages = append(c.ThreadPages, ref)
	return nil
}

func (c *BoardContainer) Touch() {
	c.Seq++
}
