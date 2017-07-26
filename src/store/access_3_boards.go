package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
)

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

func (a *Access) GetBoardPage(ctx context.Context, in *object.BoardIO) (*BoardPageOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	result, e := content.GetBoardPageResult(ctx, cxo, in)
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
	in.UserPubKey = file.User.PublicKey
	in.UserSecKey = file.User.SecretKey

	i, e := file.FindMaster(in.GetBoardPK())
	if e != nil {
		// TODO: RPC.
		return nil, e
	}

	in.BoardSecKey = file.Masters[i].SecKey

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

func (a *Access) GetThreadPage(ctx context.Context, in *object.ThreadIO) (*ThreadPageOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	result, e := content.GetThreadPageResult(ctx, cxo, in)
	if e != nil {
		return nil, e
	}

	return getThreadPage(ctx, result), nil
}

func (a *Access) NewPost(ctx context.Context, in *object.NewPostIO) (*ThreadPageOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()

	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	in.UserPubKey = file.User.PublicKey
	in.UserSecKey = file.User.SecretKey

	i, e := file.FindMaster(in.GetBoardPK())
	if e != nil {
		// TODO: Via RPC.
		return nil, e
	}
	in.BoardSecKey = file.Masters[i].SecKey

	result, e := content.NewPost(ctx, cxo, in)
	if e != nil {
		return nil, e
	}

	return getThreadPage(ctx, result), nil
}

func (a *Access) DeletePost(ctx context.Context, in *object.PostIO) (*ThreadPageOutput, error) {
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

	result, e := content.DeletePost(ctx, cxo, in)
	if e != nil {
		return nil, e
	}

	return getThreadPage(ctx, result), nil
}
