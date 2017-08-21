package state

import (
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
)

type PackAction func(p *skyobject.Pack, h *pack.Headers) error

type PackInstance struct {
	pack    *skyobject.Pack
	headers *pack.Headers
}

func NewPackInstance(oldPI *PackInstance, p *skyobject.Pack) (*PackInstance, error) {
	newPI := &PackInstance{pack: p}
	var e error
	newPI.headers, e = pack.NewHeaders(oldPI.headers, p)
	if e != nil {
		return nil, e
	}
	return newPI, nil
}

func (pi *PackInstance) Do(action PackAction) error {
	return action(pi.pack, pi.headers)
}
