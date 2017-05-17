package extern

import (
	"errors"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/store"
)

type RPCGateway struct {
	config     *cmd.Config
	container  *store.Container
	boardSaver *store.BoardSaver
	userSaver  *store.UserSaver
}

func NewRPCGateway(
	config *cmd.Config,
	container *store.Container,
	boardSaver *store.BoardSaver,
	userSaver *store.UserSaver,
) *RPCGateway {
	return &RPCGateway{
		config:     config,
		container:  container,
		boardSaver: boardSaver,
		userSaver:  userSaver,
	}
}

func (g *RPCGateway) NewPost(req *store.ReqNewPost, ok *bool) (e error) {
	if req == nil || req.Post == nil || ok == nil {
		return errors.New("nil error")
	}
	// Check post.
	if e := req.Post.Verify(); e != nil {
		*ok = false
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		*ok = false
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.BoardConfig.Master == false {
		*ok = false
		return errors.New("not master of board")
	}
	*ok = true
	return g.container.NewPost(req.BoardPubKey, req.ThreadRef, req.Post)
}

func (g *RPCGateway) NewThread(req *store.ReqNewThread, ok *bool) (e error) {
	if req == nil || req.Thread == nil || ok == nil {
		return errors.New("nil error")
	}
	// Check thread.
	if e := req.Thread.Verify(req.Creator, req.Signature); e != nil {
		*ok = false
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		*ok = false
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.BoardConfig.Master == false {
		*ok = false
		return errors.New("not master of board")
	}
	// Modify thread.
	req.Thread.MasterBoard = req.BoardPubKey.Hex()
	return g.container.NewThread(req.BoardPubKey, req.Thread)
}
