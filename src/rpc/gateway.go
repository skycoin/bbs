package rpc

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/store"
	"log"
)

type Gateway struct {
	container  *store.CXO
	boardSaver *store.BoardSaver
	userSaver  *store.UserSaver
}

func NewGateway(
	container *store.CXO,
	boardSaver *store.BoardSaver,
	userSaver *store.UserSaver,
) *Gateway {
	return &Gateway{
		container:  container,
		boardSaver: boardSaver,
		userSaver:  userSaver,
	}
}

func (g *Gateway) PingPong(_, ok *bool) error {
	*ok = true
	return nil
}

func (g *Gateway) NewPost(req *ReqNewPost, pRefStr *string) error {
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
	// Create new post.
	if e := g.container.NewPost(req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef, req.Post); e != nil {
		return e
	}
	// Modify post.
	*pRefStr = req.Post.Ref
	return nil
}

func (g *Gateway) NewThread(req *ReqNewThread, tRefStr *string) error {
	log.Println("[RPCGATEWAY] NewThread request recieved. Processing...")
	if req == nil || req.Thread == nil {
		return errors.New("nil error")
	}
	// Check thread.
	if e := req.Thread.Verify(); e != nil {
		return e
	}
	req.Thread.Touch()
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
	var e error
	switch vote.Mode {
	case 0:
		e = g.container.RemoveVoteForPost(
			vote.User, req.BoardPubKey, bi.Config.GetSK(), req.PostRef)
	case -1, +1:
		e = g.container.AddVoteForPost(
			req.BoardPubKey, bi.Config.GetSK(), req.PostRef, vote)
	default:
		e = errors.Errorf("invalid vote mode '%d'", vote.Mode)
	}
	if e != nil {
		return e
	}
	*ok = true
	return nil
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
	var e error
	switch vote.Mode {
	case 0:
		e = g.container.RemoveVoteForThread(
			vote.User, req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef)
	case -1, +1:
		e = g.container.AddVoteForThread(
			req.BoardPubKey, bi.Config.GetSK(), req.ThreadRef, vote)
	default:
		e = errors.Errorf("invalid vote mode '%d'", vote.Mode)
	}
	if e != nil {
		return e
	}
	*ok = true
	return nil
}

func (g *Gateway) VoteUser(req *ReqVoteUser, ok *bool) error {
	log.Println("[RPCGATEWAY] VoteUser request recieved. Processing...")
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
	var e error
	switch vote.Mode {
	case 0:
		e = g.container.RemoveVoteForUser(
			vote.User, req.BoardPubKey, req.UserPubKey, bi.Config.GetSK())
	case -1, +1:
		e = g.container.AddVoteForUser(
			req.BoardPubKey, req.UserPubKey, bi.Config.GetSK(), vote)
	default:
		e = errors.Errorf("invalid vote mode '%d'", vote.Mode)
	}
	if e != nil {
		return e
	}
	*ok = true
	return nil
}
