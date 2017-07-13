package gui

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
	"strconv"
)

// Posts represents the posts endpoint group.
type Posts struct {
	*Gateway
	Votes PostVotes
}

// Get obtains posts of specified board and thread.
func (g *Posts) Get(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	posts, e := g.get(bpk, tRef)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, posts, http.StatusOK)
}

func (g *Posts) get(bpk cipher.PubKey, tRef skyobject.Reference) ([]*typ.Post, error) {
	_, has := g.boardSaver.Get(bpk)
	if has == false {
		return nil, errors.New("not subscribed to board")
	}
	return g.container.GetPosts(bpk, tRef)
}

// Add adds a new post on specified board and thread.
func (g *Posts) Add(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get request values.
	post := &typ.Post{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	}
	if e := g.add(bpk, tRef, post); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, PostView{Post: post, Votes: &VotesView{}}, http.StatusOK)
}

func (g *Posts) add(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) (e error) {
	// Check post.
	uc := g.userSaver.GetCurrent()
	if e = post.Sign(uc.GetPK(), uc.GetSK()); e != nil {
		return
	}
	post.Touch()
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via CXO.
		return g.container.NewPost(bpk, bi.Config.GetSK(), tRef, post)
	} else {
		// Via RPC Client.
		return g.queueSaver.AddNewPostReq(bpk, tRef, post)
	}
}

// Remove removes a post on specified board and thread.
func (g *Posts) Remove(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// // Get request reference.
	pRef, e := misc.GetReference(r.FormValue("post"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := g.remove(bpk, tRef, pRef); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *Posts) remove(bpk cipher.PubKey, tRef, pRef skyobject.Reference) (e error) {
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		if e = g.container.RemovePost(bpk, bi.Config.GetSK(), tRef, pRef); e != nil {
			fmt.Println(e)
			return e
		}
	} else {
		// threads and posts are only to be deleted from master.
		return errors.New("not master of board")
	}
	return nil
}

// PostVotes represents the post votes endpoint group.
type PostVotes struct {
	*Gateway
}

// Get gets votes for post of specified board and post.
func (g *PostVotes) Get(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get post reference.
	pRef, e := misc.GetReference(r.FormValue("post"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get posts.
	vv, e := g.get(bpk, pRef)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, vv, http.StatusOK)
}

func (g *PostVotes) get(bpk cipher.PubKey, pRef skyobject.Reference) (*VotesView, error) {
	// Get current user.
	cu := g.userSaver.GetCurrent()
	upk := cu.GetPK()
	// Get votes.
	votes := g.container.GetVotesForPost(bpk, pRef)
	vv := &VotesView{}
	for _, vote := range votes {
		switch vote.Mode {
		case +1:
			vv.UpVotes += 1
		case -1:
			vv.DownVotes += 1
		}
		if vote.User == upk {
			vv.CurrentUserVoted = true
			vv.CurrentUserVoteMode = int(vote.Mode)
		}
	}
	return vv, nil
}

// Add adds a vote for post of specified board and post.
func (g *PostVotes) Add(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get post reference.
	pRef, e := misc.GetReference(r.FormValue("post"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get vote mode (up/down vote).
	mode, e := strconv.Atoi(r.FormValue("mode"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Prepare vote.
	vote := &typ.Vote{Mode: int8(mode), Tag: []byte(r.FormValue("tag"))}
	if e := g.add(bpk, pRef, vote); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *PostVotes) add(bpk cipher.PubKey, pRef skyobject.Reference, vote *typ.Vote) error {
	// Get current user.
	uc := g.userSaver.GetCurrent()
	// Check vote.
	if e := vote.Sign(uc.GetPK(), uc.GetSK()); e != nil {
		return errors.Wrap(e, "vote signing failed")
	}
	// Check board.
	bi, got := g.boardSaver.Get(bpk)
	if !got {
		return errors.Errorf("not subscribed to board '%s'", bpk.Hex())
	}
	g.container.GetStateSaver().AddPostVote(bpk, pRef, vote)
	// Check if this node owns the board.
	if bi.Config.Master {
		// Via CXO.
		switch vote.Mode {
		case 0:
			return g.container.RemoveVoteForPost(uc.GetPK(), bpk, bi.Config.GetSK(), pRef)
		case -1, +1:
			return g.container.AddVoteForPost(bpk, bi.Config.GetSK(), pRef, vote)
		}
	} else {
		return g.queueSaver.AddVotePostReq(bpk, pRef, vote)
	}
	return nil
}
