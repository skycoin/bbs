package rpc

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/cxo"
	"log"
)

type Gateway struct {
	config     *args.Config
	container  *cxo.Container
	boardSaver *store.BoardSaver
	userSaver  *store.UserSaver
}

func NewGateway(
	config *args.Config,
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

func (g *Gateway) PingPong(_, ok *bool) error {
	*ok = true
	return nil
}

func (g *Gateway) NewPost(req *ReqNewPost, pRefStr *string) (e error) {
	log.Println("[RPCGATEWAY] NewPost request recieved. Processing...")
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
	*pRefStr = req.Post.Ref
	return g.container.NewPost(req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef, req.Post)
}

func (g *Gateway) NewThread(req *ReqNewThread, tRefStr *string) error {
	log.Println("[RPCGATEWAY] NewThread request recieved. Processing...")
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
	*tRefStr = req.Thread.Ref
	return nil
}

func (g *Gateway) VotePost(req *ReqVotePost, ok *bool) error {
	log.Println("[RPCGATEWAY] VotePost request recieved. Processing...")
	if req == nil || ok == nil {
		return errors.New("nil error")
	}
	// Check vote.
	vote := req.Vote
	if e := vote.Verify(); e != nil {
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if !has {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if !bi.Config.Master {
		return errors.New("not master of board")
	}
	// Do vote.
	switch vote.Mode {
	case 0:
		return g.container.RemoveVoteForPost(
			vote.User, req.BoardPubKey, bi.Config.GetSK(), req.PostRef)
	case -1, +1:
		return g.container.AddVoteForPost(
			req.BoardPubKey, bi.Config.GetSK(), req.PostRef, vote)
	default:
		return errors.Errorf("invalid vote mode '%d'", vote.Mode)
	}
}

func (g *Gateway) VoteThread(req *ReqVoteThread, ok *bool) error {
	log.Println("[RPCGATEWAY] VoteThread request recieved. Processing...")
	if req == nil || ok == nil {
		return errors.New("nil error")
	}
	// Check vote.
	vote := req.Vote
	if e := vote.Verify(); e != nil {
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(req.BoardPubKey)
	if !has {
		return errors.Errorf("not subscribed to board '%s'",
			req.BoardPubKey.Hex())
	}
	// Check if this BBS Node owns the board.
	if !bi.Config.Master {
		return errors.Errorf("not master of board '%s'",
			req.BoardPubKey.Hex())
	}
	// Do vote.
	switch vote.Mode {
	case 0:
		return g.container.RemoveVoteForThread(
			vote.User, req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef)
	case -1, +1:
		return g.container.AddVoteForThread(
			req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef, vote)
	default:
		return errors.Errorf("invalid vote mode '%d'", vote.Mode)
	}
}
