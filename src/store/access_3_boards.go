package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
)

type BoardsOutput struct {
	Boards       []object.BoardView `json:"boards"`
	MasterBoards []object.BoardView `json:"master_boards"`
}

func getBoards(ctx context.Context, cxo *state.CXO, file *state.UserFile) *BoardsOutput {

	masters := make([]object.BoardView, len(file.Masters))
	for i, sub := range file.Masters {
		masters[i].PublicKey = sub.PubKey.Hex()
		result, e := content.GetBoardResult(ctx, cxo, sub.PubKey)
		if e != nil {
			masters[i].Name = "Unavailable Board"
			masters[i].Desc = e.Error()
		} else {
			masters[i].Board = *result.Board
		}

	}

	subs := make([]object.BoardView, len(file.Subscriptions))
	for i, sub := range file.Subscriptions {
		subs[i].PublicKey = sub.PubKey.Hex()
		result, e := content.GetBoardResult(ctx, cxo, sub.PubKey)
		if e != nil {
			subs[i].Name = "Unavailable Board"
			subs[i].Desc = e.Error()
		} else {
			subs[i].Board = *result.Board
		}
	}

	return &BoardsOutput{MasterBoards: masters, Boards: subs}
}

func (a *Access) GetBoards(ctx context.Context) (*BoardsOutput, error) {
	cxo := a.Session.GetCXO()

	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	return getBoards(ctx, cxo, file), nil
}

func (a *Access) NewBoard(ctx context.Context, in *object.NewBoardIO) (*BoardsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()

	file, e := a.Session.NewMaster(ctx, in)
	if e != nil {
		return nil, e
	}

	if e := content.NewBoard(ctx, cxo, in); e != nil {
		in := &object.BoardIO{PubKey: in.GetPK().Hex()}
		in.Process()
		a.Session.DeleteMaster(ctx, in)
		return nil, e
	}

	return getBoards(ctx, cxo, file), nil
}

func (a *Access) DeleteBoard(ctx context.Context, in *object.BoardIO) (*BoardsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	i, e := file.FindMaster(in.GetPK())
	if e != nil {
		return nil, e
	}

	in.SecKey = file.Masters[i].SecKey

	file.Masters = append(file.Masters[:i], file.Masters[i+1:]...)

	if e := content.DeleteBoard(ctx, cxo, in); e != nil {
		return nil, e
	}

	return getBoards(ctx, cxo, file), nil
}

func (a *Access) NewSubmissionAddress(ctx context.Context, in *object.AddressIO) (*BoardsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	i, e := file.FindMaster(in.GetPK())
	if e != nil {
		return nil, e
	}

	in.SecKey = file.Masters[i].SecKey

	if e := content.NewSubmissionAddress(ctx, cxo, in); e != nil {
		return nil, e
	}

	return getBoards(ctx, cxo, file), nil
}

func (a *Access) DeleteSubmissionAddress(ctx context.Context, in *object.AddressIO) (*BoardsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	i, e := file.FindMaster(in.GetPK())
	if e != nil {
		return nil, e
	}

	in.SecKey = file.Masters[i].SecKey

	if e := content.DeleteSubmissionAddress(ctx, cxo, in); e != nil {
		return nil, e
	}

	return getBoards(ctx, cxo, file), nil
}