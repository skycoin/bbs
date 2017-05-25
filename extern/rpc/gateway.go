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

func (g *Gateway) NewPost(req *ReqNewPost, ok *bool) (e error) {
	if req == nil || req.Post == nil || ok == nil {
		return errors.New("nil error")
	}
	// Check post.
	if e := req.Post.Verify(); e != nil {
		*ok = false
		return e
	}
	req.Post.Touch()
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		*ok = false
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == false {
		*ok = false
		return errors.New("not master of board")
	}
	*ok = true
	return g.container.NewPost(req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef, req.Post)
}

func (g *Gateway) NewThread(req *ReqNewThread, ok *bool) error {
	log.Println("[RPCGATEWAY] Recieved NewThread Request.")
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
	if bi.Config.Master == false {
		*ok = false
		return errors.New("not master of board")
	}
	// Create new thread.
	if e := g.container.NewThread(req.BoardPubKey, bi.Config.GetSK(), req.Thread); e != nil {
		return e
	}
	// Modify thread.
	req.Thread.MasterBoard = req.BoardPubKey.Hex()
	return nil
}
