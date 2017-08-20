package content_view

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
	"sync"
)

type BoardView struct {
	PubKey       string   `json:"public_key"`
	Name         string   `json:"name"`
	Desc         string   `json:"description"`
	Created      int64    `json:"created"`
	SubAddresses []string `json:"submission_addresses"`
}

type ContentView struct {

}

func (v *ContentView) Init(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error {
	pages, e := object.GetPages(pack, mux)
	if e != nil {
		return e
	}

	board, e := pages.BoardPage.GetBoard(mux)
	if e != nil {
		return e
	}
	_ = object.GetData(board)

	return nil
}

func (v *ContentView) Update(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error {
	return nil
}

func (v *ContentView) Get(id string, a ...interface{}) (interface{}, error) {
	return nil, nil
}
