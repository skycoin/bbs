package btp

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/obj/view"
)

type GetBoardsInput struct {
	SortBy string `json:"sort_by"`
}

type GetBoardsOutput struct {
	MasterBoards []view.BoardView `json:"master_boards,omitempty"`
	Boards       []view.BoardView `json:"boards"`
}

func (a *BoardAccessor) GetBoards(in *GetBoardsInput) (*GetBoardsOutput, error) {
	defer a.lock()()

	mbOut := []view.BoardView{}
	for pkStr := range a.bFile.MasterBoards {
		pk, e := misc.GetPubKey(pkStr)
		if e != nil {
			fmt.Println(e)
			continue
		}
		view, e := a.stateSaver.Get(pk)
		if e != nil {
			fmt.Println(e)
			continue
		}
		mbOut = append(mbOut, view)
	}

	bOut := []view.BoardView{}
	for pkStr := range a.bFile.Boards {
		pk, e := misc.GetPubKey(pkStr)
		if e != nil {
			fmt.Println(e)
			continue
		}
		view, e := a.stateSaver.Get(pk)
		if e != nil {
			fmt.Println(e)
			continue
		}
		bOut = append(bOut, view)
	}

	out := &GetBoardsOutput{
		MasterBoards: mbOut,
		Boards:       bOut,
	}

	return out, nil
}
