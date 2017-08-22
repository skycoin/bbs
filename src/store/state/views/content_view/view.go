package content_view

import (
	"fmt"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"sync"
)

type ContentView struct {
	sync.Mutex
	board *BoardRep
	tMap  map[cipher.SHA256]*ThreadRep
	pMap  map[cipher.SHA256]*PostRep
}

func (v *ContentView) Init(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error {
	v.Lock()
	defer v.Unlock()

	pages, e := object.GetPages(pack, mux, true)
	if e != nil {
		return e
	}

	// Set board.
	board, e := pages.BoardPage.GetBoard(mux)
	if e != nil {
		return e
	}
	v.board = new(BoardRep).Fill(pages.PK, board)
	v.board.Threads = make([]cipher.SHA256, pages.BoardPage.GetThreadCount())

	v.tMap = make(map[cipher.SHA256]*ThreadRep)
	v.pMap = make(map[cipher.SHA256]*PostRep)

	// Fill threads and posts.
	e = pages.BoardPage.RangeThreadPages(func(i int, tp *object.ThreadPage) error {
		v.board.Threads[i] = tp.Thread.Hash

		threadRep := new(ThreadRep).FillThreadPage(tp, nil)

		v.tMap[tp.Thread.Hash] = threadRep

		// Fill posts.
		threadRep.Posts = make([]cipher.SHA256, tp.GetPostCount())

		e = tp.RangePosts(func(i int, post *object.Post) error {
			threadRep.Posts[i] = post.R
			v.pMap[post.R] = new(PostRep).Fill(post, nil)
			return nil
		}, nil)
		if e != nil {
			return e
		}

		return nil
	}, mux)

	if e != nil {
		return e
	}

	return nil
}

func (v *ContentView) Update(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error {
	v.Lock()
	defer v.Unlock()

	pages, e := object.GetPages(pack, mux, true)
	if e != nil {
		return e
	}
	oldBoard := v.board
	bRaw, e := pages.BoardPage.GetBoard(mux)
	if e != nil {
		return e
	}
	newBoard := new(BoardRep).Fill(pages.PK, bRaw)
	newBoard.Threads = oldBoard.Threads

	changes := headers.GetChanges()

	for _, thread := range changes.NewThreads {
		fmt.Printf("NEW THREAD: %s\n", thread.R.Hex())
		newBoard.Threads = append(v.board.Threads, thread.R)
		v.tMap[thread.R] = new(ThreadRep).FillThread(thread, mux)
	}
	v.board = newBoard

	for _, post := range changes.NewPosts {
		if ofThread, ok := v.tMap[post.OfThread]; !ok {
			log.Println("thread not found")
			continue
		} else {
			ofThread.Posts = append(ofThread.Posts, post.R)
		}
		v.pMap[post.R] = new(PostRep).Fill(post, mux)
	}
	// TODO: Votes.
	return nil
}
