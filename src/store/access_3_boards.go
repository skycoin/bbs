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
			masters[i].Board = &object.Board{
				Name: "Unavailable Board",
				Desc: e.Error(),
			}
		} else {
			masters[i].Board = result.Board
		}

	}

	subs := make([]object.BoardView, len(file.Subscriptions))
	for i, sub := range file.Subscriptions {
		subs[i].PublicKey = sub.PubKey.Hex()
		result, e := content.GetBoardResult(ctx, cxo, sub.PubKey)
		if e != nil {
			subs[i].Board = &object.Board{
				Name: "Unavailable Board",
				Desc: e.Error(),
			}
		} else {
			subs[i].Board = result.Board
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

	file, e := a.Session.DeleteMaster(ctx, in)
	if e != nil {
		return nil, e
	}

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

type BoardPageOutput struct {
	Board   object.BoardView
	Threads []object.ThreadView
}

func getBoardPage(_ context.Context, result *content.Result) *BoardPageOutput {

	out := &BoardPageOutput{
		Board: object.BoardView{
			Board:     result.Board,
			PublicKey: result.GetPK().Hex(),
		},
		Threads: make([]object.ThreadView, len(result.Threads)),
	}

	for i, thread := range result.Threads {
		out.Threads[i] = object.ThreadView{
			Thread:      thread,
			Ref:         thread.R.Hex(),
			AuthorRef:   thread.User.Hex(),
			AuthorAlias: "-", // TODO: Implement.
			Votes:       nil, // TODO: Implement.
		}
	}

	return out
}

func (a *Access) GetBoardPage(ctx context.Context, in *object.BoardIO) (*BoardPageOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	result, e := content.GetThreadsResult(ctx, cxo, in)
	if e != nil {
		return nil, e
	}

	return getBoardPage(ctx, result), nil
}

func (a *Access) NewThread(ctx context.Context, in *object.NewThreadIO) (*BoardPageOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	i, e := file.FindMaster(in.GetBoardPK())
	if e != nil {
		// TODO: RPC.
		return nil, e
	}

	in.BoardSecKey = file.Masters[i].SecKey
	in.UserPubKey = file.User.PublicKey
	in.UserSecKey = file.User.SecretKey

	result, e := content.NewThread(ctx, cxo, in)
	if e != nil {
		return nil, e
	}

	return getBoardPage(ctx, result), nil
}

func (a *Access) DeleteThread(ctx context.Context, in *object.ThreadIO) (*BoardPageOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	i, e := file.FindMaster(in.GetBoardPK())
	if e != nil {
		// TODO: RPC.
		return nil, e
	}

	in.BoardSecKey = file.Masters[i].SecKey

	result, e := content.DeleteThread(ctx, cxo, in)
	if e != nil {
		return nil, e
	}

	return getBoardPage(ctx, result), nil
}
