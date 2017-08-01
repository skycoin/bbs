package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/session"
)

// Access allows access to store.
type Access struct {
	Session *session.Manager
}

/*
	<<< SESSION >>>
*/

// GetUsers gets a list of all available users.
func (a *Access) GetUsers(ctx context.Context) (*UsersOutput, error) {
	aliases, e := a.Session.GetUsers(ctx)
	if e != nil {
		return nil, e
	}
	return getUsers(ctx, aliases), nil
}

// NewUser creates a new user.
func (a *Access) NewUser(ctx context.Context, in *object.NewUserIO) (*UsersOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	if _, e := a.Session.NewUser(ctx, in); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

// DeleteUser deletes a user.
func (a *Access) DeleteUser(ctx context.Context, alias string) (*UsersOutput, error) {
	if e := a.Session.DeleteUser(ctx, alias); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

// Login logs a user in.
func (a *Access) Login(ctx context.Context, in *object.LoginIO) (*object.UserView, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.Login(ctx, in)
	if e != nil {
		return nil, e
	}
	out := &object.UserView{
		User:      object.User{Alias: file.User.Alias},
		PublicKey: file.User.PublicKey.Hex(),
		SecretKey: file.User.SecretKey.Hex(),
	}
	return out, nil
}

// Logout logs a user out.
func (a *Access) Logout(ctx context.Context) error {
	return a.Session.Logout(ctx)
}

// GetSession obtains the current session information.
func (a *Access) GetSession(ctx context.Context) (*session.UserFileView, error) {
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	view := file.GenerateView(a.Session.GetCXO())
	return view, nil
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
	if e := object.Process(in); e != nil {
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
	if e := object.Process(in); e != nil {
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
	if e := object.Process(in); e != nil {
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
	if e := object.Process(in); e != nil {
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
	if e := object.Process(in); e != nil {
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
		object.Process(in)
		a.Session.DeleteMaster(ctx, in)
		return nil, e
	}
	a.Session.GetCompiler().Trigger(root)
	return getBoards(ctx, cxo, file), nil
}

// DeleteBoard removes a board that this user owns.
func (a *Access) DeleteBoard(ctx context.Context, in *object.BoardIO) (*BoardsOutput, error) {
	if e := object.Process(in); e != nil {
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
	if e := object.Process(in); e != nil {
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
	if e := object.Process(in); e != nil {
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
	if e := object.Process(in); e != nil {
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
	return getBoardPage(ctx, compiler, result), nil
}

// NewThread creates a new thread on a board.
func (a *Access) NewThread(ctx context.Context, in *object.NewThreadIO) (*BoardPageOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	file.FillUser(in)
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
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
	compiler.Trigger(root)
	return getBoardPage(ctx, compiler, result), nil
}

// DeleteThread removes a thread from a board.
func (a *Access) DeleteThread(ctx context.Context, in *object.ThreadIO) (*BoardPageOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
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
	compiler.Trigger(root)
	return getBoardPage(ctx, compiler, result), nil
}

func (a *Access) VoteThread(ctx context.Context, in *object.VoteThreadIO) (*VotesOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	file.FillUser(in)
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
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
	compiler.Trigger(root)
	return getThreadVotes(ctx, compiler, result, in.ThreadRef), nil
}

/*
	<<< POSTS >>>
*/

// GetThreadPage obtains a page that displays thread information and lists the thread's posts.
func (a *Access) GetThreadPage(ctx context.Context, in *object.ThreadIO) (*ThreadPageOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
	root, e := cxo.GetRoot(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	result, e := content.GetThreadPageResult(ctx, root, in)
	if e != nil {
		return nil, e
	}
	return getThreadPage(ctx, compiler, result), nil
}

// NewPost creates a new post on specified thead and board.
func (a *Access) NewPost(ctx context.Context, in *object.NewPostIO) (*ThreadPageOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	file.FillUser(in)
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
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
	compiler.Trigger(root)
	return getThreadPage(ctx, compiler, result), nil
}

// DeletePost removes a post from specified thread and board.
func (a *Access) DeletePost(ctx context.Context, in *object.PostIO) (*ThreadPageOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
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
	compiler.Trigger(root)
	return getThreadPage(ctx, compiler, result), nil
}

func (a *Access) VotePost(ctx context.Context, in *object.VotePostIO) (*VotesOutput, error) {
	if e := object.Process(in); e != nil {
		return nil, e
	}
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	file.FillUser(in)
	cxo := a.Session.GetCXO()
	defer cxo.Lock()()
	compiler := a.Session.GetCompiler()
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
	compiler.Trigger(root)
	return getPostVotes(ctx, compiler, result, in.PostRef), nil
}
