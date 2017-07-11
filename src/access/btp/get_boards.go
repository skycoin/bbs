package btp

import (
	"github.com/skycoin/bbs/src/store/view"
	"github.com/skycoin/bbs/src/misc"
	"fmt"
)

type GetBoardsInput struct {
	SortBy string `json:"sort_by"`
}

type GetBoardsOutput struct {
	MasterBoards []view.Board `json:"master_boards,omitempty"`
	Boards       []view.Board `json:"boards"`
}

func (a *BoardAccessor) GetBoards(in *GetBoardsInput) (*GetBoardsOutput, error) {
	defer a.lock()()

	mbOut := []view.Board{}
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

	bOut := []view.Board{}
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
