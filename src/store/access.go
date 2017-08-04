package store

import (
	"context"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/session"
	"github.com/skycoin/bbs/src/store/users"
)

// Access allows access to store.
type Access struct {
	Session *session.Manager
	Users   *users.Manager
}

/*
	<<< SESSION >>>
*/

// GetUsers gets a list of available users.
func (a *Access) GetUsers(ctx context.Context) (*UsersOutput, error) {
	aliases, e := a.Users.GetUsers()
	if e != nil {
		return nil, e
	}
	return getUsers(ctx, aliases), nil
}

// NewUser creates a new user.
func (a *Access) NewUser(ctx context.Context, in *object.NewUserIO) (*UsersOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	if e := a.Users.NewUser(in); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

// DeleteUser deletes a user.
func (a *Access) DeleteUser(ctx context.Context, alias string) (*UsersOutput, error) {
	if e := tag.CheckAlias(alias); e != nil {
		return nil, e
	}
	if e := a.Users.DeleteUser(alias); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

// Login logs a user in.
func (a *Access) Login(ctx context.Context, in *object.LoginIO) (*users.FileView, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	out, e := a.Users.Login(in)
	if e != nil {
		return nil, e
	}
	return out.GenerateView(), nil
}

// Logout logs a user out.
func (a *Access) Logout(ctx context.Context) error {
	return a.Users.Logout()
}

// GetSession obtains the current session information.
func (a *Access) GetSession(ctx context.Context) (*session.FileView, error) {
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	return file.GenerateView(a.Session.GetCXO()), nil
}

/*
	<<< CONNECTIONS >>>
*/

// GetConnections gets list of external node address that this node is currently connected to.
func (a *Access) GetConnections(ctx context.Context) (*ConnectionsOutput, error) {
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	return getConnections(a.Session.GetCXO(), file)
}

// NewConnection creates a new connection.
func (a *Access) NewConnection(ctx context.Context, in *object.ConnectionIO) (*ConnectionsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.NewConnection(ctx, in)
	if e != nil {
		return nil, e
	}
	return getConnections(a.Session.GetCXO(), file)
}

// DeleteConnection removes a connection.
func (a *Access) DeleteConnection(ctx context.Context, in *object.ConnectionIO) (*ConnectionsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.DeleteConnection(ctx, in)
	if e != nil {
		return nil, e
	}
	return getConnections(a.Session.GetCXO(), file)
}

/*
	<<< SUBSCRIPTIONS >>>
*/

// GetSubs gets list of boards this node/user is currently subscribed to.
func (a *Access) GetSubs(ctx context.Context) (*SubsOutput, error) {
	cxo := a.Session.GetCXO()
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	return getSubs(ctx, cxo, file), nil
}

// NewSub subscribes node/user to a board that this user does not own.
func (a *Access) NewSub(ctx context.Context, in *object.BoardIO) (*SubsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	file, e := a.Session.NewSubscription(ctx, in)
	if e != nil {
		return nil, e
	}
	return getSubs(ctx, cxo, file), nil
}

// DeleteSub removes a subscription to a board that this user does not own.
func (a *Access) DeleteSub(ctx context.Context, in *object.BoardIO) (*SubsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	file, e := a.Session.DeleteSubscription(ctx, in)
	if e != nil {
		return nil, e
	}
	return getSubs(ctx, cxo, file), nil
}

/*
	<<< BOARDS >>>
*/

// GetBoards obtains list of boards.
func (a *Access) GetBoards(ctx context.Context) (*BoardsOutput, error) {
	cxo := a.Session.GetCXO()
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	return getBoards(ctx, cxo, file), nil
}

// NewBoard creates a new board that this user owns.
func (a *Access) NewBoard(ctx context.Context, in *object.NewBoardIO) (*BoardsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.NewMaster(ctx, in)
	if e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	root, e := cxo.NewRoot(in.BoardPubKey, in.BoardSecKey)
	if e != nil {
		return nil, e
	}
	if e := content.NewBoard(ctx, root, in); e != nil {
		in := &object.BoardIO{PubKeyStr: in.BoardPubKey.Hex()}
		tag.Process(in)
		a.Session.DeleteMaster(ctx, in)
		return nil, e
	}
	a.Session.GetCompiler().Trigger(root)
	return getBoards(ctx, cxo, file), nil
}

// DeleteBoard removes a board that this user owns.
func (a *Access) DeleteBoard(ctx context.Context, in *object.BoardIO) (*BoardsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.DeleteMaster(ctx, in)
	if e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	root, e := cxo.GetRoot(in.PubKey, in.SecKey)
	if e != nil {
		return nil, e
	}
	if e := content.DeleteBoard(ctx, root, in); e != nil {
		return nil, e
	}
	a.Session.GetCompiler().DeleteBoard(root.Pub())
	return getBoards(ctx, cxo, file), nil
}

// NewSubmissionAddress adds a submission address to a board which this user owns.
func (a *Access) NewSubmissionAddress(ctx context.Context, in *object.AddressIO) (*BoardsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	if e := file.FillMaster(in); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	root, e := cxo.GetRoot(in.PubKey, in.SecKey)
	if e != nil {
		return nil, e
	}
	if e := content.NewSubmissionAddress(ctx, root, in); e != nil {
		return nil, e
	}
	a.Session.GetCompiler().Trigger(root)
	return getBoards(ctx, cxo, file), nil
}

// DeleteSubmissionAddress removes a submission address from a board that this user owns.
func (a *Access) DeleteSubmissionAddress(ctx context.Context, in *object.AddressIO) (*BoardsOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	if e := file.FillMaster(in); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	root, e := cxo.GetRoot(in.PubKey, in.SecKey)
	if e != nil {
		return nil, e
	}
	if e := content.DeleteSubmissionAddress(ctx, root, in); e != nil {
		return nil, e
	}
	a.Session.GetCompiler().Trigger(root)
	return getBoards(ctx, cxo, file), nil
}

/*
	<<< THREADS >>>
*/

// GetBoardPage obtains a page that displays board information and lists the board's threads.
func (a *Access) GetBoardPage(ctx context.Context, in *object.BoardIO) (*BoardPageOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
	root, e := cxo.GetRoot(in.PubKey)
	if e != nil {
		return nil, e
	}
	result, e := content.GetBoardPageResult(ctx, root, in)
	if e != nil {
		return nil, e
	}
	return getBoardPage(ctx, compiler, result, a.Users.GetUPK()), nil
}

// NewThread creates a new thread on a board.
func (a *Access) NewThread(ctx context.Context, in *object.NewThreadIO) (*BoardPageOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	in.Thread = new(object.Thread)
	tag.Transfer(in, &in.Thread.Post)
	if e := a.Users.Sign(&in.Thread.Post); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	if e := file.FillMaster(in); e != nil {
		//TODO: RPC
		return nil, e
	}
	root, e := cxo.GetRoot(in.BoardPubKey, in.BoardSecKey)
	if e != nil {
		return nil, e
	}
	result, e := content.NewThread(ctx, root, in)
	if e != nil {
		return nil, e
	}
	compiler := a.Session.GetCompiler()
	compiler.Trigger(root)
	return getBoardPage(ctx, compiler, result, a.Users.GetUPK()), nil
}

// DeleteThread removes a thread from a board.
func (a *Access) DeleteThread(ctx context.Context, in *object.ThreadIO) (*BoardPageOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	if e := file.FillMaster(in); e != nil {
		// TODO: RPC
		return nil, e
	}
	root, e := cxo.GetRoot(in.BoardPubKey, in.BoardSecKey)
	if e != nil {
		return nil, e
	}
	result, e := content.DeleteThread(ctx, root, in)
	if e != nil {
		return nil, e
	}
	compiler := a.Session.GetCompiler()
	compiler.Trigger(root)
	return getBoardPage(ctx, compiler, result, a.Users.GetUPK()), nil
}

func (a *Access) VoteThread(ctx context.Context, in *object.VoteThreadIO) (*VotesOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	in.Vote = new(object.Vote)
	tag.Transfer(in, in.Vote)
	if e := a.Users.Sign(in.Vote); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	if e := file.FillMaster(in); e != nil {
		// TODO: RPC
		return nil, e
	}
	root, e := cxo.GetRoot(in.BoardPubKey, in.BoardSecKey)
	if e != nil {
		return nil, e
	}
	result, e := content.VoteThread(ctx, root, in)
	if e != nil {
		return nil, e
	}
	compiler := a.Session.GetCompiler()
	compiler.Trigger(root)
	return getThreadVotes(ctx, compiler, result, a.Users.GetUPK(), in.ThreadRef), nil
}

/*
	<<< POSTS >>>
*/

// GetThreadPage obtains a page that displays thread information and lists the thread's posts.
func (a *Access) GetThreadPage(ctx context.Context, in *object.ThreadIO) (*ThreadPageOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	root, e := cxo.GetRoot(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	result, e := content.GetThreadPageResult(ctx, root, in)
	if e != nil {
		return nil, e
	}
	compiler := a.Session.GetCompiler()
	return getThreadPage(ctx, compiler, result, a.Users.GetUPK()), nil
}

// NewPost creates a new post on specified thead and board.
func (a *Access) NewPost(ctx context.Context, in *object.NewPostIO) (*ThreadPageOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	in.Post = new(object.Post)
	tag.Transfer(in, in.Post)
	if e := a.Users.Sign(in.Post); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	if e := file.FillMaster(in); e != nil {
		// TODO: RPC
		return nil, e
	}
	root, e := cxo.GetRoot(in.BoardPubKey, in.BoardSecKey)
	if e != nil {
		return nil, e
	}
	result, e := content.NewPost(ctx, root, in)
	if e != nil {
		return nil, e
	}
	compiler := a.Session.GetCompiler()
	compiler.Trigger(root)
	return getThreadPage(ctx, compiler, result, a.Users.GetUPK()), nil
}

// DeletePost removes a post from specified thread and board.
func (a *Access) DeletePost(ctx context.Context, in *object.PostIO) (*ThreadPageOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	if e := file.FillMaster(in); e != nil {
		// TODO: RPC
		return nil, e
	}
	root, e := cxo.GetRoot(in.BoardPubKey, in.BoardSecKey)
	if e != nil {
		return nil, e
	}
	result, e := content.DeletePost(ctx, root, in)
	if e != nil {
		return nil, e
	}
	compiler := a.Session.GetCompiler()
	compiler.Trigger(root)
	return getThreadPage(ctx, compiler, result, a.Users.GetUPK()), nil
}

func (a *Access) VotePost(ctx context.Context, in *object.VotePostIO) (*VotesOutput, error) {
	if e := tag.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}

	in.Vote = new(object.Vote)
	tag.Transfer(in, in.Vote)
	if e := a.Users.Sign(in.Vote); e != nil {
		return nil, e
	}

	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	if e := file.FillMaster(in); e != nil {
		// TODO: RPC
		return nil, e
	}
	root, e := cxo.GetRoot(in.BoardPubKey, in.BoardSecKey)
	if e != nil {
		return nil, e
	}
	result, e := content.VotePost(ctx, root, in)
	if e != nil {
		return nil, e
	}
	compiler := a.Session.GetCompiler()
	compiler.Trigger(root)
	return getPostVotes(ctx, compiler, result, a.Users.GetUPK(), in.PostRef), nil
}
