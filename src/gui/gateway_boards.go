package gui

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
	"strings"
)

// Boards represents the boards endpoint group.
type Boards struct {
	*Gateway
	Meta BoardMeta
	Page BoardPage
}

// GetAll gets all boards.
func (g *Boards) GetAll(w http.ResponseWriter, r *http.Request) {
	send(w, g.getAll(), http.StatusOK)
}

func (g *Boards) getAll() []*typ.Board {
	return g.container.GetBoards(g.boardSaver.ListKeys()...)
}

// Get obtains a single board.
func (g *Boards) Get(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	board, e := g.get(bpk)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
	}
	send(w, board, http.StatusOK)
}

func (g *Boards) get(bpk cipher.PubKey) (*typ.Board, error) {
	board, e := g.container.GetBoard(bpk)
	return board, errors.Wrap(e, "unable to obtain board")
}

// Add creates a new board.
func (g *Boards) Add(w http.ResponseWriter, r *http.Request) {
	// Obtain Board Meta.
	meta := new(typ.BoardMeta)
	meta.SubmissionAddresses =
		strings.Split(r.FormValue("submission_addresses"), ",")
	meta.Trim()

	// Generate board.
	board := &typ.Board{
		Name: r.FormValue("name"),
		Desc: r.FormValue("description"),
	}
	board.SetMeta(meta)

	// Create board in cxo.
	bi, e := g.add(board, r.FormValue("seed"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
	} else {
		send(w, bi, http.StatusOK)
	}
}

func (g *Boards) add(board *typ.Board, seed string) (bi store.BoardInfo, e error) {
	bm, e := board.GetMeta()
	if e != nil {
		e = errors.Wrap(e, "failed to obtain board meta")
		return
	}
	if len(bm.SubmissionAddresses) == 0 {
		bm.AddSubmissionAddress(g.config.RPCRemoteAddr)
		if e = board.SetMeta(bm); e != nil {
			e = errors.Wrap(e, "failed to re-set board meta")
			return
		}
	}
	bpk, bsk := board.TouchWithSeed([]byte(seed))
	if _, e = board.Check(); e != nil {
		e = errors.Wrap(e, "failed to create board")
		return
	}
	if e = g.boardSaver.MasterAdd(bpk, bsk); e != nil {
		return
	}
	if e = g.container.NewBoard(board, bpk, bsk); e != nil {
		return
	}
	bi, _ = g.boardSaver.Get(bpk)
	return
}

// Remove removes a board.
func (g *Boards) Remove(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := g.remove(bpk); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
	} else {
		send(w, true, http.StatusOK)
	}
}

func (g *Boards) remove(bpk cipher.PubKey) error {
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to the board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via CXO.
		return g.container.RemoveBoard(bpk, bi.Config.GetSK())
	} else {
		// threads and posts are only to be deleted from master.
		return errors.New("not master of board")
	}
}

// BoardMeta represents the boards meta endpoint group.
type BoardMeta struct {
	*Gateway
	SubmissionAddresses BoardsMetaSubmissionAddresses
}

// Get obtains the board meta.
func (g *BoardMeta) Get(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	meta, e := g.get(bpk)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, meta, http.StatusOK)
}

func (g *BoardMeta) get(bpk cipher.PubKey) (*typ.BoardMeta, error) {
	board, e := g.container.GetBoard(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain board")
	}
	meta, e := board.GetMeta()
	return meta, errors.Wrap(e, "unable to obtain board meta")
}

// BoardsMetaSubmissionAddresses represents the boards meta submission addresses endpoint group.
type BoardsMetaSubmissionAddresses struct {
	*Gateway
}

// GetAll gets all submission addresses of board from board meta.
func (g *BoardsMetaSubmissionAddresses) GetAll(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	addresses, e := g.getAll(bpk)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, addresses, http.StatusOK)
}

func (g *BoardsMetaSubmissionAddresses) getAll(bpk cipher.PubKey) ([]string, error) {
	list, e := g.container.GetSubmissionAddresses(bpk)
	return list, errors.Wrap(e, "failed to obtain submission addresses of board")
}

// Add adds a submission address to board via board meta.
func (g *BoardsMetaSubmissionAddresses) Add(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	address := r.FormValue("address")
	if address == "" {
		send(w, "no submission address provided", http.StatusBadRequest)
		return
	}
	if e := g.add(bpk, address); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *BoardsMetaSubmissionAddresses) add(bpk cipher.PubKey, address string) error {
	bi, has := g.boardSaver.Get(bpk)
	if !has {
		return errors.New("board not found")
	}
	if !bi.Config.Master {
		return errors.New("not master of board")
	}
	return g.container.AddSubmissionAddress(bpk, bi.Config.GetSK(), address)
}

// Remove removes a submission address from board via board meta.
func (g *BoardsMetaSubmissionAddresses) Remove(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	address := r.FormValue("address")
	if address == "" {
		send(w, "no submission address provided", http.StatusBadRequest)
		return
	}
	if e := g.remove(bpk, address); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *BoardsMetaSubmissionAddresses) remove(bpk cipher.PubKey, address string) error {
	bi, has := g.boardSaver.Get(bpk)
	if !has {
		return errors.New("board not found")
	}
	if !bi.Config.Master {
		return errors.New("not master of board")
	}
	return g.container.RemoveSubmissionAddress(bpk, bi.Config.GetSK(), address)
}

// BoardPage represents the board page endpoint group.
type BoardPage struct {
	*Gateway
}

// Get gets the board page.
func (g *BoardPage) Get(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	bpv, e := g.get(bpk)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, bpv, http.StatusOK)
}

func (g *BoardPage) get(bpk cipher.PubKey) (*BoardPageView, error) {
	board, e := g.container.GetBoard(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain board")
	}
	threads, e := g.container.GetThreads(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain threads")
	}
	threadViews := make([]*ThreadView, len(threads))
	for i := range threadViews {
		votesView, e := g.Threads.Votes.get(bpk, threads[i].GetRef())
		if e != nil {
			return nil, errors.Wrap(e, "unable to obtain votes for thread")
		}
		threadViews[i] = &ThreadView{
			Thread: threads[i],
			Votes:  votesView,
		}
	}
	return &BoardPageView{board, threadViews}, nil
}

// TODO: Implement.
