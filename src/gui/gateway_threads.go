package gui

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
	"strconv"
)

// Threads represents the threads endpoint group.
type Threads struct {
	*Gateway
	Page  ThreadPage
	Votes ThreadVotes
}

// GetAll obtains all threads of specified board(s).
func (g *Threads) GetAll(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	if bpkStr == "" {
		send(w, g.getAll(), http.StatusOK)
		return
	}
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, g.getAll(bpk), http.StatusOK)
}

func (g *Threads) getAll(bpks ...cipher.PubKey) []*typ.Thread {
	tMap := make(map[string]*typ.Thread)
	switch len(bpks) {
	case 0:
		for _, bpk := range g.boardSaver.ListKeys() {
			ts, e := g.container.GetThreads(bpk)
			if e != nil {
				continue
			}
			for _, t := range ts {
				tMap[t.Ref] = t
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
				tMap[t.Ref] = t
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

// Add adds a new thread to specified board.
func (g *Threads) Add(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread values.
	thread := &typ.Thread{
		Name:        r.FormValue("name"),
		Desc:        r.FormValue("description"),
		MasterBoard: bpk.Hex(),
	}
	if e := g.add(bpk, thread); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, thread, http.StatusOK)
}

func (g *Threads) add(bpk cipher.PubKey, thread *typ.Thread) error {
	// Check thread.
	if e := thread.Check(); e != nil {
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to the board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via Container.
		return g.container.NewThread(bpk, bi.Config.GetSK(), thread)
	} else {
		// Via RPC Client.
		uc := g.userSaver.GetCurrent()
		return g.queueSaver.AddNewThreadReq(bpk, uc.GetPK(), uc.GetSK(), thread)
	}
}

// Remove removes a thread from specified board.
func (g *Threads) Remove(w http.ResponseWriter, r *http.Request) {
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
	if e := g.remove(bpk, tRef); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *Threads) remove(bpk cipher.PubKey, tRef skyobject.Reference) error {
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via Container.

		// Obtain thread.
		thread, e := g.container.GetThread(tRef)
		if e != nil {
			return e
		}
		// Obtain thread master's public key.
		masterPK, e := misc.GetPubKey(thread.MasterBoard)
		if e != nil {
			return e
		}
		// Remove dependency (if has).
		bi.Config.RemoveDep(masterPK, tRef)
		// Remove thread.
		if e := g.container.RemoveThread(bpk, bi.Config.GetSK(), tRef); e != nil {
			return e
		}
	} else {
		// threads and posts are only to be deleted from master.
		return errors.New("not owning the board")
	}
	return nil
}

// Import imports a thread from one board to another.
func (g *Threads) Import(w http.ResponseWriter, r *http.Request) {
	// Get "from" board's public key.
	fromBpk, e := misc.GetPubKey(r.FormValue("from_board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread's reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get "to" board's public key.
	toBpk, e := misc.GetPubKey(r.FormValue("to_board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Import thread.
	if e := g.importThread(fromBpk, toBpk, tRef); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *Threads) importThread(fromBpk, toBpk cipher.PubKey, tRef skyobject.Reference) error {
	// Check "to" board.
	bi, has := g.boardSaver.Get(toBpk)
	if !has {
		return errors.New("not subscribed to board")
	}
	if bi.Config.Master == false {
		return errors.New("to board is not master")
	}
	// Add Dependency to BoardSaver.
	if e := g.boardSaver.AddBoardDep(toBpk, fromBpk, tRef); e != nil {
		return e
	}
	// Import thread.
	return g.container.ImportThread(fromBpk, toBpk, bi.Config.GetSK(), tRef)
}

// ThreadPage represents the thread page endpoint group.
type ThreadPage struct {
	*Gateway
}

// Get gets the thread page of specified board and thread.
func (g *ThreadPage) Get(w http.ResponseWriter, r *http.Request) {
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
	// Get thread page.
	threadPage, e := g.get(bpk, tRef)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, threadPage, http.StatusOK)
}

func (g *ThreadPage) get(bpk cipher.PubKey, tRef skyobject.Reference) (*ThreadPageView, error) {
	b, e := g.container.GetBoard(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain board")
	}
	thread, posts, e := g.container.GetThreadPage(bpk, tRef)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain threadpage")
	}
	return &ThreadPageView{b, thread, posts}, nil
}

// ThreadVotes represents the thread votes endpoint group.
type ThreadVotes struct {
	*Gateway
}

// Get gets votes for thread of specified board and thread.
func (g *ThreadVotes) Get(w http.ResponseWriter, r *http.Request) {
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
	// Get votes.
	vv, e := g.get(bpk, tRef)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, vv, http.StatusOK)
}

func (g *ThreadVotes) get(bpk cipher.PubKey, tRef skyobject.Reference) (*VotesView, error) {
	// Get current user.
	cu := g.userSaver.GetCurrent()
	upk := cu.GetPK()
	// Get votes.
	votes, e := g.container.GetVotesForThread(bpk, tRef)
	if e != nil {
		return nil, e
	}
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

// Add adds a vote for thread of specified board and thread.
func (g *ThreadVotes) Add(w http.ResponseWriter, r *http.Request) {
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
	// Get vote mode (up/down vote).
	mode, e := strconv.Atoi(r.FormValue("mode"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Prepare vote.
	vote := &typ.Vote{Mode: int8(mode), Tag: []byte(r.FormValue("tag"))}
	if e := g.add(bpk, tRef, vote); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *ThreadVotes) add(bpk cipher.PubKey, tRef skyobject.Reference, vote *typ.Vote) error {
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
	// Check if this node owns the board.
	if bi.Config.Master {
		// Via Container.
		switch vote.Mode {
		case 0:
			return g.container.RemoveVoteForThread(uc.GetPK(), bpk, bi.Config.GetSK(), tRef)
		case -1, +1:
			return g.container.AddVoteForThread(bpk, bi.Config.GetSK(), tRef, vote)
		}
	} else {
		return g.queueSaver.AddVoteThreadReq(bpk, tRef, vote)
	}
	return nil
}
