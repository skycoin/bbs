package extern

import (
	"errors"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/store"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// Gateway represents the intermediate between External calls and internal processing.
// It can be seen as a security layer.
type Gateway struct {
	config     *cmd.Config
	container  *store.Container
	boardSaver *store.BoardSaver
	userSaver  *store.UserSaver
}

// NewGateway creates a new Gateway.
func NewGateway(
	config *cmd.Config,
	container *store.Container,
	boardSaver *store.BoardSaver,
	userSaver *store.UserSaver,
) *Gateway {
	return &Gateway{
		config:     config,
		container:  container,
		boardSaver: boardSaver,
		userSaver:  userSaver,
	}
}

/*
	<<< FOR SUBSCRIPTIONS >>>
*/

// GetSubscriptions lists all subscriptions.
func (g *Gateway) GetSubscriptions() []store.BoardInfo {
	return g.boardSaver.List()
}

// GetSubscription gets a subscription.
func (g *Gateway) GetSubscription(bpk cipher.PubKey) (store.BoardInfo, bool) {
	return g.boardSaver.Get(bpk)
}

// Subscribe subscribes to a board.
func (g *Gateway) Subscribe(bpk cipher.PubKey) {
	g.boardSaver.Add(bpk)
}

// Unsubscribe unsubscribes from a board.
func (g *Gateway) Unsubscribe(bpk cipher.PubKey) {
	g.boardSaver.Remove(bpk)
}

/*
	<<< FOR USERS >>>
*/

// GetCurrentUser returns the currently active user.
func (g *Gateway) GetCurrentUser() store.UserConfig {
	return g.userSaver.GetCurrent()
}

// SetCurrent sets the currently active user.
func (g *Gateway) SetCurrent(upk cipher.PubKey) error {
	return g.userSaver.SetCurrent(upk)
}

// NewMasterUser creates a new user configuration of a master user.
// It will replace if user already exists.
func (g *Gateway) NewMasterUser(alias, seed string) store.UserConfig {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	g.userSaver.MasterAdd(alias, pk, sk)
	uc, _ := g.userSaver.Get(pk)
	return uc
}

// GetMasterUsers lists all users this node is master of.
func (g *Gateway) GetMasterUsers() []store.UserConfig {
	return g.userSaver.ListMasters()
}

// GetUsers lists all users, master or not.
func (g *Gateway) GetUsers() []store.UserConfig {
	return g.userSaver.List()
}

// NewUser creates a new user configuration for a user we are not master of.
// It will replace if user already exists.
func (g *Gateway) NewUser(alias string, upk cipher.PubKey) store.UserConfig {
	g.userSaver.Add(alias, upk)
	uc, _ := g.userSaver.Get(upk)
	return uc
}

// RemoveUser removes a user configuration, master or not.
func (g *Gateway) RemoveUser(upk cipher.PubKey) error {
	return g.userSaver.Remove(upk)
}

/*
	<<< FOR BOARDS, THREADS & POSTS >>>
*/

// GetBoards lists all boards.
func (g *Gateway) GetBoards() []*typ.Board {
	return g.container.GetBoards(g.boardSaver.ListKeys()...)
}

// NewBoard creates a new board.
func (g *Gateway) NewBoard(board *typ.Board, seed string) (bi store.BoardInfo, e error) {
	bpk, bsk, e := g.container.NewBoard(board, seed)
	if e != nil {
		return
	}
	if e = g.boardSaver.MasterAdd(bpk, bsk); e != nil {
		return
	}
	bi, _ = g.boardSaver.Get(bpk)
	return
}

// GetThreads obtains threads of boards we are subscribed to.
// Input `bpks` acts as a filter.
// If no `bpks` are specified, threads of all boards will be obtained.
// If one or more `bpks` are specified, only threads under those boards will be returned.
func (g *Gateway) GetThreads(bpks ...cipher.PubKey) []*typ.Thread {
	tMap := make(map[string]*typ.Thread)
	switch len(bpks) {
	case 0:
		for _, bpk := range g.boardSaver.ListKeys() {
			ts, e := g.container.GetThreads(bpk)
			if e != nil {
				continue
			}
			for _, t := range ts {
				tMap[t.Hash] = t
			}
		}
	default:
		for _, bpk := range bpks {
			if _, has := g.boardSaver.Get(bpk); has == false {
				return nil
			}
			ts, e := g.container.GetThreads(bpk)
			if e != nil {
				continue
			}
			for _, t := range ts {
				tMap[t.Hash] = t
			}
		}
	}
	threads, i := make([]*typ.Thread, len(tMap)), 0
	for _, t := range tMap {
		threads[i] = t
		i += 1
	}
	return threads
}

// NewThread creates a new thread and sets the board of public key as it's master.
func (g *Gateway) NewThread(bpk cipher.PubKey, thread *typ.Thread) error {
	// Check thread.
	if e := thread.Check(); e != nil {
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.BoardConfig.Master == true {
		// Via Container.
		if e := g.container.NewThread(bpk, thread); e != nil {
			return e
		}
	} else {
		// Via RPC Client.
		// TODO: Implement.
		return errors.New("not implemented")
	}
	return nil
}

// GetPosts obtains posts of specified board and thread.
// TODO: In the future, as a single thread can exist across different boards, we will only need to specify the thread.
func (g *Gateway) GetPosts(bpk cipher.PubKey, tRef skyobject.Reference) ([]*typ.Post, error) {
	_, has := g.boardSaver.Get(bpk)
	if has == false {
		return nil, errors.New("not subscribed to board")
	}
	return g.container.GetPosts(bpk, tRef)
}

// NewPost creates a new post in specified board and thread.
// TODO: In the future, as a single thread can exist across different boards, we will only need to specify the thread.
func (g *Gateway) NewPost(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) error {
	// Check post.
	uc := g.userSaver.GetCurrent()
	if e := post.Sign(uc.GetPK(), uc.GetSK()); e != nil {
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.BoardConfig.Master == true {
		// Via Container.
		return g.container.NewPost(bpk, tRef, post)
	} else {
		// Via RPC Client.
		// TODO: Implement.
		return errors.New("not implemented")
	}
	return nil
}
