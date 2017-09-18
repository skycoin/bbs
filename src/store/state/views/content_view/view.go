package content_view

import (
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
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
	vMap  map[cipher.SHA256]*VotesRep
}

func (v *ContentView) Init(pack *skyobject.Pack, headers *pack.Headers) error {
	v.Lock()
	defer v.Unlock()

	pages, e := r0.GetPages(pack, false, true, false, true)
	if e != nil {
		return e
	}

	// Set board.
	board, e := pages.BoardPage.GetBoard()
	if e != nil {
		return e
	}
	v.board = new(BoardRep).Fill(pages.PK, board)
	v.board.Threads = make([]IndexHash, pages.BoardPage.GetThreadCount())

	v.tMap = make(map[cipher.SHA256]*ThreadRep)
	v.pMap = make(map[cipher.SHA256]*PostRep)
	v.vMap = make(map[cipher.SHA256]*VotesRep)

	// Fill threads and posts.
	e = pages.BoardPage.RangeThreadPages(func(i int, tp *r0.ThreadPage) error {
		v.board.Threads[i] = IndexHash{h: tp.Thread.Hash, i: i}

		threadRep := new(ThreadRep).FillThreadPage(tp)

		// Fill posts.
		threadRep.Posts = make([]IndexHash, tp.GetPostCount())

		e = tp.RangePosts(func(i int, post *r0.Post) error {
			threadRep.Posts[i] = IndexHash{h: post.R, i: i}
			v.pMap[post.R] = new(PostRep).Fill(post)
			return nil
		})
		if e != nil {
			return e
		}

		// Add thread rep to map.
		v.tMap[tp.Thread.Hash] = threadRep

		return nil
	})

	if e != nil {
		return e
	}

	return pages.UsersPage.RangeUserActivityPages(func(_ int, uap *r0.UserActivityPage) error {
		return uap.RangeVoteActions(func(_ int, vote *r0.Vote) error {
			return v.processVote(vote)
		})
	})
}

func (v *ContentView) Update(pack *skyobject.Pack, headers *pack.Headers) error {
	v.Lock()
	defer v.Unlock()

	pages, e := r0.GetPages(pack, false, true)
	if e != nil {
		return e
	}
	bRaw, e := pages.BoardPage.GetBoard()
	if e != nil {
		return e
	}
	v.board.Fill(pages.PK, bRaw)

	changes := headers.GetChanges()

	for _, thread := range changes.NewThreads {
		v.board.Threads = append(v.board.Threads, IndexHash{h: thread.R, i: len(v.board.Threads)})
		v.tMap[thread.R] = new(ThreadRep).FillThread(thread)
	}

	for _, post := range changes.NewPosts {
		if ofThread, ok := v.tMap[post.GetData().GetOfThread()]; !ok {
			log.Println("thread not found")
			continue
		} else {
			ofThread.Posts = append(ofThread.Posts, IndexHash{h: post.R, i: len(ofThread.Posts)})
		}
		v.pMap[post.R] = new(PostRep).Fill(post)
	}

	for _, vote := range changes.NewVotes {
		v.processVote(vote)
	}

	return nil
}

func (v *ContentView) processVote(vote *r0.Vote) error {
	var cHash cipher.SHA256

	// Only if vote is for post or thread.
	switch vote.GetType() {
	case r0.UserVote, r0.UnknownVoteType:
		return nil

	case r0.ThreadVote:
		if v.tMap[vote.OfThread] == nil {
			return nil
		}
		cHash = vote.OfThread

	case r0.PostVote:
		if v.pMap[vote.OfPost] == nil {
			return nil
		}
		cHash = vote.OfPost
	}

	// Add to votes map.
	voteRep, has := v.vMap[cHash]
	if !has {
		voteRep = new(VotesRep).Fill(cHash)
		v.vMap[cHash] = voteRep
	}
	voteRep.Add(vote)

	return nil
}
