package rpc

import (
	"errors"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/intern/cxo"
	"github.com/evanlinjin/bbs/intern/store"
	"log"
)

type Gateway struct {
	config     *cmd.Config
	container  *cxo.Container
	boardSaver *store.BoardSaver
	userSaver  *store.UserSaver
}

func NewGateway(
	config *cmd.Config,
	container *cxo.Container,
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

func (g *Gateway) NewPost(req *ReqNewPost, pRefStr *string) (e error) {
	if req == nil || req.Post == nil {
		return errors.New("nil error")
	}
	// Check post.
	if e := req.Post.Verify(); e != nil {
		return e
	}
	req.Post.Touch()
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == false {
		return errors.New("not master of board")
	}
	pRefStr = &req.Post.Ref
	return g.container.NewPost(req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef, req.Post)
}

func (g *Gateway) NewThread(req *ReqNewThread, tRefStr *string) error {
	log.Println("[RPCGATEWAY] Recieved NewThread Request.")
	if req == nil || req.Thread == nil {
		return errors.New("nil error")
	}
	// Check thread.
	if e := req.Thread.Verify(req.Creator, req.Signature); e != nil {
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == false {
		return errors.New("not master of board")
	}
	// Create new thread.
	if e := g.container.NewThread(req.BoardPubKey, bi.Config.GetSK(), req.Thread); e != nil {
		return e
	}
	// Modify thread.
	req.Thread.MasterBoard = req.BoardPubKey.Hex()
	tRefStr = &req.Thread.Ref
	return nil
}
