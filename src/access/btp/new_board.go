package btp

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/bbs/src/store/obj/view"
	"github.com/skycoin/cxo/node"
	"time"
)

// NewBoardInput is configuration struct used when creating a new board.
type NewBoardInput struct {
	Name string `json:"name"`
	Desc string `json:"description"`
	Seed string `json:"seed"`
}

// NewBoard creates a new board.
func (a *BoardAccessor) NewBoard(in *NewBoardInput) (*view.BoardView, error) {
	if !a.cxo.IsMaster() {
		return nil, boo.New(boo.NotMaster)
	}

	board := &obj.Board{
		Name:    in.Name,
		Desc:    in.Desc,
		Created: time.Now().UnixNano(),
	}

	pk, sk, e := a.cxo.NewRoot([]byte(in.Seed), func(r *node.Root) error {
		boardPageDyn, e := r.Dynamic("BoardPage", obj.BoardPage{Board: r.Save(*board)})
		if e != nil {
			return boo.WrapType(e, boo.Internal, "failed to save board page to root")
		}
		if _, e := r.Append(boardPageDyn); e != nil {
			return boo.WrapType(e, boo.Internal, "failed to append to root")
		}
		a.stateSaver.Update(r)
		return nil
	})

	if e != nil {
		return nil, boo.WrapType(e, boo.Internal, "failed to create root")
	}

	a.bFile.AddMaster(pk, sk)

	boardView := &view.BoardView{
		Board:     *board,
		PublicKey: pk.Hex(),
	}

	return boardView, nil
}
