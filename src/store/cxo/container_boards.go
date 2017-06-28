package cxo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
)

// ChangeBoardURL changes the board's URL of given public key.
func (c *Container) ChangeBoardURL(bpk cipher.PubKey, bsk cipher.SecKey, url string) error {
	c.Lock(c.ChangeBoardURL)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	b := &typ.Board{}
	if e := w.AdvanceFromRefField("Board", b); e != nil {
		return e
	}
	//b.URL = url // .......
	//w.Retreat()
	//_, e := w.ReplaceInRefField("Board", *b)
	//return e
	if e := w.ReplaceCurrent(b); e != nil {
		return e
	}
	return nil
}

// ChangeBoardMeta changes the board meta data.
func (c *Container) ChangeBoardMeta(bpk cipher.PubKey, bsk cipher.SecKey, bm *typ.BoardMeta) error {
	c.Lock(c.ChangeBoardMeta)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := new(typ.BoardContainer)
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	b := new(typ.Board)
	if e := w.AdvanceFromRefField("Board", b); e != nil {
		return e
	}
	if e := b.SetMeta(bm); e != nil {
		return e
	}
	if e := w.ReplaceCurrent(b); e != nil {
		return e
	}
	return nil
}

// GetSubmissionAddresses attempts to obtain all submission addresses from board meta.
func (c *Container) GetSubmissionAddresses(bpk cipher.PubKey) ([]string, error) {
	c.Lock(c.GetSubmissionAddresses)
	defer c.Unlock()

	w := c.c.LastFullRoot(bpk).Walker()
	bc := new(typ.BoardContainer)
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, errors.Wrap(e, "failed to obtain board container")
	}
	b := new(typ.Board)
	if e := w.AdvanceFromRefField("Board", b); e != nil {
		return nil, errors.Wrap(e, "failed to obtain board")
	}
	bm, e := b.GetMeta()
	if e != nil {
		return nil, errors.Wrap(e, "failed to obtain board meta")
	}
	return bm.SubmissionAddresses, nil
}

// AddSubmissionAddress attempts to add a submission address to board meta.
func (c *Container) AddSubmissionAddress(bpk cipher.PubKey, bsk cipher.SecKey, address string) error {
	c.Lock(c.AddSubmissionAddress)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := new(typ.BoardContainer)
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "failed to obtain board container")
	}
	b := new(typ.Board)
	if e := w.AdvanceFromRefField("Board", b); e != nil {
		return errors.Wrap(e, "failed to obtain board")
	}
	bm, e := b.GetMeta()
	if e != nil {
		return errors.Wrap(e, "failed to obtain board meta")
	}
	if e := bm.AddSubmissionAddress(address); e != nil {
		return errors.Wrap(e, "failed to add submission address")
	}
	if e := b.SetMeta(bm); e != nil {
		return errors.Wrap(e, "failed to replace board meta")
	}
	if e := w.ReplaceCurrent(b); e != nil {
		return errors.Wrap(e, "failed to save changes to board")
	}
	return nil
}

// RemoveSubmissionAddress attempts to remove a submission address from board meta.
func (c *Container) RemoveSubmissionAddress(bpk cipher.PubKey, bsk cipher.SecKey, address string) error {
	c.Lock(c.RemoveSubmissionAddress)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := new(typ.BoardContainer)
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "failed to obtain board container")
	}
	b := new(typ.Board)
	if e := w.AdvanceFromRefField("Board", b); e != nil {
		return errors.Wrap(e, "failed to obtain board")
	}
	bm, e := b.GetMeta()
	if e != nil {
		return errors.Wrap(e, "failed to obtain board meta")
	}
	// Remove submission address.
	bm.RemoveSubmissionAddress(address)

	if e := b.SetMeta(bm); e != nil {
		return errors.Wrap(e, "failed to replace board meta")
	}
	if e := w.ReplaceCurrent(b); e != nil {
		return errors.Wrap(e, "failed to save changes to board")
	}
	return nil
}

// GetBoard attempts to obtain the board of a given public key.
func (c *Container) GetBoard(bpk cipher.PubKey) (*typ.Board, error) {
	c.Lock(c.GetBoard)
	defer c.Unlock()

	w := c.c.LastFullRoot(bpk).Walker()
	bc := typ.BoardContainer{}
	if e := w.AdvanceFromRoot(&bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	b := &typ.Board{}
	e := w.AdvanceFromRefField("Board", b)
	return b, e
}

// GetBoards attempts to obtain a list of boards from the given public keys.
func (c *Container) GetBoards(bpks ...cipher.PubKey) []*typ.Board {
	c.Lock(c.GetBoards)
	defer c.Unlock()

	boards := []*typ.Board{}
	for _, bpk := range bpks {
		w := c.c.LastRoot(bpk).Walker()
		bc, b := typ.BoardContainer{}, typ.Board{}
		if e := w.AdvanceFromRoot(&bc, makeBoardContainerFinder(w.Root())); e != nil {
			continue
		}
		w.AdvanceFromRefField("Board", &b)
		boards = append(boards, &b)
	}
	return boards
}

// NewBoard attempts to create a new board from a given board and seed.
func (c *Container) NewBoard(board *typ.Board, pk cipher.PubKey, sk cipher.SecKey) error {
	c.Lock(c.NewBoard)
	defer c.Unlock()

	r, e := c.c.NewRoot(pk, sk)
	if e != nil {
		return e
	}
	bRef := r.Save(*board)
	// Prepare board container.
	bCont := typ.BoardContainer{Board: bRef}
	if _, _, e = r.Inject("BoardContainer", bCont); e != nil {
		return e
	}
	// Prepare thread vote container.
	tvCont := typ.ThreadVotesContainer{}
	if _, _, e := r.Inject("ThreadVotesContainer", tvCont); e != nil {
		return e
	}
	// Prepare post vote container.
	pvCont := typ.PostVotesContainer{}
	if _, _, e := r.Inject("PostVotesContainer", pvCont); e != nil {
		return e
	}
	return nil
}

// RemoveBoard attempts to remove a board by a given public key.
func (c *Container) RemoveBoard(bpk cipher.PubKey, bsk cipher.SecKey) error {
	c.Lock(c.RemoveBoard)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	fmt.Println("Removing board:", bpk.Hex())
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	return w.RemoveCurrent()
}

// GetThread obtains a single thread via reference.
func (c *Container) GetThread(tRef skyobject.Reference) (*typ.Thread, error) {
	c.Lock(c.GetThread)
	defer c.Unlock()

	tData, has := c.c.Get(tRef)
	if !has {
		return nil, errors.New("thread not found")
	}
	thread := &typ.Thread{}
	if e := encoder.DeserializeRaw(tData, thread); e != nil {
		return nil, e
	}
	thread.Ref = tRef.String()
	return thread, nil
}

// GetThreads attempts to obtain a list of threads from a board of public key.
func (c *Container) GetThreads(bpk cipher.PubKey) ([]*typ.Thread, error) {
	c.Lock(c.GetThreads)
	defer c.Unlock()

	w := c.c.LastFullRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	threads := make([]*typ.Thread, len(bc.Threads))
	for i, tRef := range bc.Threads {
		tData, has := c.c.Get(tRef)
		if has == false {
			continue
		}
		threads[i] = new(typ.Thread)
		if e := encoder.DeserializeRaw(tData, threads[i]); e != nil {
			return nil, e
		}
		threads[i].Ref = cipher.SHA256(tRef).Hex()
	}
	return threads, nil
}

// NewThread attempts to create a new thread from a board of given public key.
func (c *Container) NewThread(bpk cipher.PubKey, bsk cipher.SecKey, thread *typ.Thread) error {
	c.Lock(c.NewThread)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	thread.MasterBoard = bpk.Hex()
	tRef, e := w.AppendToRefsField("Threads", *thread)
	if e != nil {
		return e
	}
	tp := typ.ThreadPage{Thread: tRef}
	if _, e := w.AppendToRefsField("ThreadPages", tp); e != nil {
		return e
	}
	thread.Ref = cipher.SHA256(tRef).Hex()
	// Prepare thread vote container.
	w.Clear()
	tvc := &typ.ThreadVotesContainer{}
	if e := w.AdvanceFromRoot(tvc, makeThreadVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	tvc.AddThread(tRef)
	if e := w.ReplaceCurrent(*tvc); e != nil {
		return e
	}
	return nil
}

// RemoveThread attempts to remove a thread from a board of given public key.
func (c *Container) RemoveThread(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference) error {
	c.Lock(c.RemoveThread)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	if e := w.RemoveInRefsByRef("Threads", tRef); e != nil {
		return errors.Wrap(e, "remove thread failed")
	}
	if e := w.RemoveInRefsField("ThreadPages", makeThreadPageFinder(w, tRef)); e != nil {
		return errors.Wrap(e, "remove thread page failed")
	}

	// remove thread votes.
	w.Clear()
	tvc := &typ.ThreadVotesContainer{}
	if e := w.AdvanceFromRoot(tvc, makeThreadVotesContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "obtaining thread vote container failed")
	}
	tvc.RemoveThread(tRef)
	if e := w.ReplaceCurrent(*tvc); e != nil {
		return errors.Wrap(e, "swapping thread vote container failed")
	}
	return nil
}

// GetThreadPage requests a page from a thread
func (c *Container) GetThreadPage(bpk cipher.PubKey, tRef skyobject.Reference) (*typ.Thread, []*typ.Post, error) {
	c.Lock(c.GetThreadPage)
	defer c.Unlock()

	w := c.c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, nil, e
	}
	// Get thread.
	tData, has := c.c.Get(tRef)
	if has == false {
		return nil, nil, errors.New("unable to obtain thread")
	}
	thread := new(typ.Thread)
	if e := encoder.DeserializeRaw(tData, thread); e != nil {
		return nil, nil, e
	}
	thread.Ref = cipher.SHA256(tRef).Hex()
	// Get posts.
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return nil, nil, e
	}
	posts := make([]*typ.Post, len(tp.Posts))
	for i, pRef := range tp.Posts {
		pData, has := c.c.Get(pRef)
		if has == false {
			continue
		}
		posts[i] = new(typ.Post)
		if e := encoder.DeserializeRaw(pData, posts[i]); e != nil {
			return nil, nil, e
		}
		posts[i].Ref = cipher.SHA256(pRef).Hex()
	}
	return thread, posts, nil
}

// GetPosts attempts to obtain posts from a specified board and thread.
func (c *Container) GetPosts(bpk cipher.PubKey, tRef skyobject.Reference) ([]*typ.Post, error) {
	c.Lock(c.GetPosts)
	defer c.Unlock()

	w := c.c.LastFullRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return nil, e
	}
	posts := make([]*typ.Post, len(tp.Posts))
	for i, pRef := range tp.Posts {
		pData, has := c.c.Get(pRef)
		if has == false {
			continue
		}
		posts[i] = new(typ.Post)
		if e := encoder.DeserializeRaw(pData, posts[i]); e != nil {
			return nil, e
		}
		posts[i].Ref = cipher.SHA256(pRef).Hex()
	}
	return posts, nil
}

// NewPost attempts to create a new post in a given board and thread.
func (c *Container) NewPost(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference, post *typ.Post) error {
	c.Lock(c.NewPost)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return e
	}
	t := &typ.Thread{}
	if e := w.GetFromRefField("Thread", t); e != nil {
		return e
	}
	if t.MasterBoard != bpk.Hex() {
		return errors.New("this board is not master of this thread")
	}
	var pRef skyobject.Reference
	var e error
	if pRef, e = w.AppendToRefsField("Posts", *post); e != nil {
		return e
	}
	post.Ref = cipher.SHA256(pRef).Hex()

	// Prepare post vote container.
	w.Clear()
	pvc := &typ.PostVotesContainer{}
	if e := w.AdvanceFromRoot(pvc, makePostVotesContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "unable to obtain post vote container")
	}
	pvc.AddPost(pRef)
	if e := w.ReplaceCurrent(*pvc); e != nil {
		return e
	}
	return nil
}

// RemovePost attempts to remove a post in a given board and thread.
func (c *Container) RemovePost(bpk cipher.PubKey, bsk cipher.SecKey, tRef, pRef skyobject.Reference) error {
	c.Lock(c.RemovePost)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return e
	}
	if e := w.RemoveInRefsByRef("Posts", pRef); e != nil {
		return errors.Wrap(e, "post removal failed")
	}

	// Remove post votes.
	w.Clear()
	pvc := &typ.PostVotesContainer{}
	if e := w.AdvanceFromRoot(pvc, makePostVotesContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "unable to obtain post vote container")
	}
	pvc.RemovePost(pRef)
	if e := w.ReplaceCurrent(*pvc); e != nil {
		return errors.Wrap(e, "unable to replace post vote container")
	}
	return nil
}

// ImportThread imports a thread from a board to another board (which this node owns). If already imported replaces it.
func (c *Container) ImportThread(fromBpk, toBpk cipher.PubKey, toBsk cipher.SecKey, tRef skyobject.Reference) error {
	c.Lock(c.ImportThread)
	defer c.Unlock()

	log.Printf("[CONTAINER] Syncing thread '%s' from board '%s' to board '%s'.",
		tRef.String(), fromBpk.Hex(), toBpk.Hex())

	// Get from 'from' Board.
	w := c.c.LastRoot(fromBpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain board "+fromBpk.Hex())
	}

	// Obtain thread and thread page.
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain thread page for board "+fromBpk.Hex())
	}
	t := &typ.Thread{}
	if e := w.GetFromRefField("Thread", t); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain thread for board "+fromBpk.Hex())
	}

	// Get from 'to' Board.
	w = c.c.LastRootSk(toBpk, toBsk).Walker()
	bc = &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain board "+toBpk.Hex())
	}
	if e := w.ReplaceInRefsField("ThreadPages", *tp, makeThreadPageFinder(w, tRef)); e != nil {
		/* THREAD DOES NOT EXIST */
		// Append thread and thread page.
		if _, e := w.AppendToRefsField("Threads", *t); e != nil {
			return e
		}
		if _, e := w.AppendToRefsField("ThreadPages", *tp); e != nil {
			return e
		}
		return nil
	}
	return nil
}
