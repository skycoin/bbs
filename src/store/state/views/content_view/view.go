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

type indexPage struct {
	Board   string
	Threads []string            // Board threads.
	Posts   map[string][]string // Key: thread hashes, Value: post hash array.
}

func newIndexPage() *indexPage {
	return &indexPage{
		Posts: make(map[string][]string),
	}
}

type ContentView struct {
	sync.Mutex
	pk cipher.PubKey
	i  *indexPage
	c  map[string]*object.ContentRep
	v  map[string]*VotesRep
}

func (v *ContentView) Init(pack *skyobject.Pack, headers *pack.Headers) error {
	v.Lock()
	defer v.Unlock()

	pages, e := object.GetPages(pack, &object.GetPagesIn{
		RootPage:  false,
		BoardPage: true,
		DiffPage:  false,
		UsersPage: true,
	})
	if e != nil {
		return e
	}

	v.pk = pack.Root().Pub
	v.i = newIndexPage()
	v.c = make(map[string]*object.ContentRep)
	v.v = make(map[string]*VotesRep)

	// Set board.
	board, e := pages.BoardPage.GetBoard()
	if e != nil {
		return e
	}
	v.i.Board = board.GetHeader().Hash
	boardRep := board.ToRep()
	boardRep.PubKey = v.pk.Hex()
	v.c[v.i.Board] = boardRep

	log.Printf("INITIATING THREADS : count(%d)", pages.BoardPage.GetThreadCount())
	v.i.Threads = make([]string, pages.BoardPage.GetThreadCount())

	// Fill threads and posts.
	e = pages.BoardPage.RangeThreadPages(func(i int, tp *object.ThreadPage) error {

		thread, e := tp.GetThread()
		if e != nil {
			return e
		}
		threadHash := thread.GetHeader().Hash

		v.i.Threads[i] = threadHash
		v.c[threadHash] = thread.ToRep()

		log.Printf("\t- [%d] THREAD : hash(%s) post_count(%d)",
			i, threadHash, tp.GetPostCount())

		// Fill posts.
		postHashes := make([]string, tp.GetPostCount())
		e = tp.RangePosts(func(i int, post *object.Content) error {

			log.Printf("\t\t- [%d] POST : hash(%s)",
				i, post.GetHeader().Hash)

			postHashes[i] = post.GetHeader().Hash
			v.c[postHashes[i]] = post.ToRep()
			return nil
		})
		if e != nil {
			log.Println("\t\t- Range post error:", e)
			return e
		}
		v.i.Posts[threadHash] = postHashes
		return nil
	})

	if e != nil {
		return e
	}

	log.Printf("INITIATING VOTES FROM USER PROFILES : user_count(%d)",
		pages.UsersPage.GetUsersLen())

	return pages.UsersPage.RangeUserProfiles(func(i int, uap *object.UserProfile) error {

		log.Printf("\t- [%d] USER : pk(%s) submission_count(%d)",
			i, uap.PubKey, uap.GetSubmissionsLen())

		return uap.RangeSubmissions(func(i int, c *object.Content) error {

			log.Printf("\t\t- [%d] SUBMISSION : type(%s) hash(%s)",
				i, c.GetBody().Type, c.GetHeader().Hash)

			return v.processVote(c)
		})
	})
}

func (v *ContentView) Update(pack *skyobject.Pack, headers *pack.Headers) error {
	v.Lock()
	defer v.Unlock()

	pages, e := object.GetPages(pack, &object.GetPagesIn{
		RootPage:  false,
		BoardPage: true,
		DiffPage:  false,
		UsersPage: false,
	})
	if e != nil {
		return e
	}

	board, e := pages.BoardPage.GetBoard()
	if e != nil {
		return e
	}

	delete(v.c, v.i.Board)
	v.i.Board = board.GetHeader().Hash
	boardRep := board.ToRep()
	boardRep.PubKey = v.pk.Hex()
	v.c[v.i.Board] = boardRep

	changes := headers.GetChanges()

	for _, content := range changes.New {
		var (
			header = content.GetHeader()
			body   = content.GetBody()
		)
		switch body.Type {
		case object.V5ThreadType:
			v.i.Threads = append(v.i.Threads, header.Hash)
			v.c[header.Hash] = content.ToRep()

		case object.V5PostType:
			posts, _ := v.i.Posts[body.OfThread]
			v.i.Posts[body.OfThread] = append(posts, header.Hash)
			v.c[header.Hash] = content.ToRep()

		case object.V5ThreadVoteType, object.V5PostVoteType:
			v.processVote(content)
		}
	}
	return nil
}

func (v *ContentView) processVote(c *object.Content) error {
	var cHash string
	var cType object.ContentType

	// Only if vote is for post or thread.
	switch c.GetBody().Type {
	case object.V5ThreadVoteType:
		cHash = c.GetBody().OfThread
		cType = object.V5ThreadVoteType

	case object.V5PostVoteType:
		cHash = c.GetBody().OfPost
		cType = object.V5PostVoteType

	default:
		return nil
	}

	if v.c[cHash] == nil {
		return nil
	}

	// Add to votes map.
	voteRep, has := v.v[cHash]
	if !has {
		voteRep = new(VotesRep).Fill(cType, cHash)
		v.v[cHash] = voteRep
	}
	voteRep.Add(c)
	fmt.Println("  >>> VOTE REPRESENTATION:", cHash)
	fmt.Println(voteRep.String())

	return nil
}
