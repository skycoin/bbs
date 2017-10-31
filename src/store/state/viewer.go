package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/typ"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

/*
	<<< INDEXER >>>
*/

type Indexer struct {
	Board   string
	Threads typ.Paginated
	Posts   map[string]typ.Paginated // key (hash of thread or post), value (list of posts)
}

func NewIndexer(init typ.PaginatedCreator) *Indexer {
	return &Indexer{
		Threads: init(),
		Posts:   make(map[string]typ.Paginated),
	}
}

/*
	<<< CONTAINER >>>
*/

type Container struct {
	content map[string]*object.ContentRep
	votes   map[string]*VotesRep
}

func NewContainer() *Container {
	return &Container{
		content: make(map[string]*object.ContentRep),
		votes:   make(map[string]*VotesRep),
	}
}

/*
	<<< VIEWER >>>
*/

type Viewer struct {
	mux   sync.Mutex
	pk    cipher.PubKey
	i     *Indexer
	c     *Container
	pInit typ.PaginatedCreator
}

func NewViewer(pack *skyobject.Pack, headers *pack.Headers, pInit typ.PaginatedCreator) (*Viewer, error) {
	v := &Viewer{
		pk:    pack.Root().Pub,
		i:     NewIndexer(pInit),
		c:     NewContainer(),
		pInit: pInit,
	}

	pages, e := object.GetPages(pack, &object.GetPagesIn{
		RootPage:  false,
		BoardPage: true,
		DiffPage:  false,
		UsersPage: true,
	})
	if e != nil {
		return nil, e
	}

	// Set board.
	if board, e := pages.BoardPage.GetBoard(); e != nil {
		return nil, e
	} else {
		v.setBoard(board)
	}

	e = pages.BoardPage.RangeThreadPages(func(i int, tp *object.ThreadPage) error {
		thread, e := tp.GetThread()
		if e != nil {
			return e
		}
		tHash, e := v.addThread(thread)
		if e != nil {
			return e
		}
		return tp.RangePosts(func(i int, post *object.Content) error {
			return v.addPost(tHash, post)
		})
	})
	if e != nil {
		return nil, e
	}

	e = pages.UsersPage.RangeUserProfiles(func(i int, uap *object.UserProfile) error {
		return uap.RangeSubmissions(func(i int, c *object.Content) error {
			return v.processVote(c)
		})
	})
	if e != nil {
		return nil, e
	}

	return v, nil
}

func (v *Viewer) lock() func() {
	v.mux.Lock()
	return v.mux.Unlock
}

func (v *Viewer) setBoard(bc *object.Content) {
	delete(v.c.content, v.i.Board)
	v.i.Board = bc.GetHeader().Hash
	rep := bc.ToRep()
	rep.PubKey = v.pk.Hex()
	v.c.content[v.i.Board] = rep
}

func (v *Viewer) addThread(tc *object.Content) (cipher.SHA256, error) {
	body := tc.GetBody()

	// Check board public key.
	if e := checkBoardRef(v.pk, body, "thread"); e != nil {
		return cipher.SHA256{}, e
	}

	tHash := tc.GetHeader().GetHash()
	v.i.Threads.Append(tHash.Hex())
	v.c.content[tHash.Hex()] = tc.ToRep()
	v.i.Posts[tHash.Hex()] = v.pInit()
	return tHash, nil
}

func (v *Viewer) addPost(tHash cipher.SHA256, pc *object.Content) error {
	body := pc.GetBody()

	// Check board public key.
	if e := checkBoardRef(v.pk, body, "post"); e != nil {
		return e
	}

	// Check thread ref.
	if e := checkThreadRef(tHash, body, "post"); e != nil {
		return e
	}

	pHash := pc.GetHeader().Hash
	v.i.Posts[tHash.Hex()].Append(pHash)
	v.c.content[pHash] = pc.ToRep()

	if ofPost, _ := body.GetOfPost(); ofPost != (cipher.SHA256{}) {
		pList, ok := v.i.Posts[ofPost.Hex()]
		if !ok {
			pList = v.pInit()
			v.i.Posts[ofPost.Hex()] = pList
		}
		pList.Append(pHash)
	}

	return nil
}

func (v *Viewer) processVote(c *object.Content) error {
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

		// TODO: User vote.

	default:
		return nil
	}

	if v.c.content[cHash] == nil {
		return nil
	}

	// Add to votes map.
	voteRep, has := v.c.votes[cHash]
	if !has {
		voteRep = new(VotesRep).Fill(cType, cHash)
		v.c.votes[cHash] = voteRep
	}
	voteRep.Add(c)

	return nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func checkBoardRef(expected cipher.PubKey, body *object.Body, what string) error {
	if got, e := body.GetOfBoard(); e != nil {
		return boo.WrapTypef(e, boo.InvalidRead, "corrupt %s", what)
	} else if got != expected {
		return boo.Newf(boo.InvalidRead,
			"misplaced %s, unmatched board public key", what)
	} else {
		return nil
	}
}

func checkThreadRef(expected cipher.SHA256, body *object.Body, what string) error {
	if got, e := body.GetOfThread(); e != nil {
		return boo.WrapTypef(e, boo.InvalidRead, "corrupt %s", what)
	} else if got != expected {
		return boo.Newf(boo.InvalidRead,
			"misplaced %s, unmatched board public key", what)
	} else {
		return nil
	}
}