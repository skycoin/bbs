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
		*ok = false
		return e
	}
	// Modify thread.
	req.Thread.MasterBoard = req.BoardPubKey.Hex()
	return nil
}

func (g *Gateway) RemoveBoard(req *ReqRemoveBoard, ok *bool) error {
	log.Println("[RPCGATEWAY] Received RemoveBoard Request.")
	if req == nil || req.BoardPubKey.Hex() == "" || ok == nil {
		*ok = false
		return errors.New("nil error")
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		*ok = false
		return errors.New("not subscribed to the board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == false {
		*ok = false
		return errors.New("not master of the board")
	}
	// Remove board.
	if e := g.container.RemoveBoard(req.BoardPubKey, bi.Config.GetSK()); e != nil {
		return e
	}
	*ok = true
	return nil
}

func (g *Gateway) RemoveThread(req *ReqRemoveThread, ok *bool) error {
	log.Println("[RPCGATEWAY] Received RemoveThread Request.")
	if req == nil || req.ThreadRef.IsBlank() || ok == nil {
		*ok = false
		return errors.New("nil error")
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		*ok = false
		return errors.New("not subscribed to the board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == false {
		*ok = false
		return errors.New("not master of the board")
	}
	// Remove thread.
	if e := g.container.RemoveThread(req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef); e != nil {
		return e
	}
	*ok = true
	return nil
}

func (g *Gateway) RemovePost(req *ReqRemovePost, ok *bool) (e error) {
	log.Println("[RPCGATEWAY] Received RemovePost Request.")
	if req == nil || req.PostRef.IsBlank() || ok == nil {
		return errors.New("nil error")
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if has == false {
		*ok = false
		return errors.New("not subscribed to the board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == false {
		*ok = false
		return errors.New("not master of the board")
	}
	// Remove post
	if e := g.container.RemovePost(req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef, req.PostRef); e != nil {
		*ok = false
		return e
	}
	*ok = true
	return nil
}
