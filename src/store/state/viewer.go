package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/typ"
	"github.com/skycoin/bbs/src/misc/typ/paginatedtypes"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

var ErrViewerNotInitialized = boo.New(boo.NotFound, "viewer is not initialized")

/*
	<<< INDEXER >>>
*/

type Indexer struct {
	Board         string
	Threads       typ.Paginated
	PostsOfThread map[string]typ.Paginated // key (hash of thread or post), value (list of posts)
	Users         typ.Paginated
}

func NewIndexer() *Indexer {
	return &Indexer{
		Threads:       paginatedtypes.NewSimple(),
		PostsOfThread: make(map[string]typ.Paginated),
		Users:         paginatedtypes.NewMapped(),
	}
}

func (i *Indexer) EnsureUsersOfUserVoteBody(body *object.Body) {
	i.Users.Append(body.Creator)
	i.Users.Append(body.OfUser)
}

/*
	<<< CONTAINER >>>
*/

type Container struct {
	content  map[string]*object.ContentRep
	votes    map[string]*VotesRep
	profiles map[string]*Profile
}

func NewContainer() *Container {
	return &Container{
		content:  make(map[string]*object.ContentRep),
		votes:    make(map[string]*VotesRep),
		profiles: make(map[string]*Profile),
	}
}

func (c *Container) GetProfile(upk string) *Profile {
	if profile, ok := c.profiles[upk]; ok {
		return profile
	} else {
		profile = NewProfile()
		c.profiles[upk] = profile
		return profile
	}
}

/*
	<<< VIEWER >>>
*/

type Viewer struct {
	mux sync.Mutex
	pk  cipher.PubKey
	i   *Indexer
	c   *Container
}

func NewViewer(pack *skyobject.Pack) (*Viewer, error) {
	v := &Viewer{
		pk: pack.Root().Pub,
		i:  NewIndexer(),
		c:  NewContainer(),
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
			if e := v.processVote(c); e != nil {
				return e
			}
			return nil
		})
	})
	if e != nil {
		return nil, e
	}

	return v, nil
}

func (v *Viewer) Update(pack *skyobject.Pack, headers *Headers) error {
	if v == nil {
		return ErrViewerNotInitialized
	}
	defer v.lock()()

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
	v.setBoard(board)

	for _, content := range headers.GetChanges().New {
		switch content.GetBody().Type {
		case object.V5ThreadType:
			if _, e := v.addThread(content); e != nil {
				return e
			}
		case object.V5PostType:
			tHash, _ := content.GetBody().GetOfThread()
			if e := v.addPost(tHash, content); e != nil {
				return e
			}
		case object.V5ThreadVoteType, object.V5PostVoteType:
			v.processVote(content)
		}
	}

	return nil
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
	v.i.PostsOfThread[tHash.Hex()] = paginatedtypes.NewMapped()
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
	if posts, ok := v.i.PostsOfThread[tHash.Hex()]; !ok {
		return boo.Newf(boo.Internal, "thread of hash %s not found", tHash.Hex())
	} else {
		posts.Append(pHash)
		v.c.content[pHash] = pc.ToRep()
	}

	if ofPost, _ := body.GetOfPost(); ofPost != (cipher.SHA256{}) {
		pList, ok := v.i.PostsOfThread[ofPost.Hex()]
		if !ok {
			pList = paginatedtypes.NewMapped()
			v.i.PostsOfThread[ofPost.Hex()] = pList
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

	case object.V5UserVoteType:
		return v.processUserVote(c)

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

func (v *Viewer) processUserVote(c *object.Content) error {
	var (
		b              = c.GetBody()
		creatorProfile = v.c.GetProfile(b.Creator)
		ofUserProfile  = v.c.GetProfile(b.OfUser)
	)

	creatorProfile.ClearVotesFor(b.OfUser)
	ofUserProfile.ClearVotesBy(b.Creator)

	switch b.Value {
	case +1:
		if b.HasTag(object.TrustTag) {
			v.i.EnsureUsersOfUserVoteBody(b)
			creatorProfile.Trusted[b.OfUser] = struct{}{}
			ofUserProfile.TrustedBy[b.Creator] = struct{}{}
		}
	case -1:
		if b.HasTag(object.SpamTag) {
			v.i.EnsureUsersOfUserVoteBody(b)
			creatorProfile.MarkedAsSpam[b.OfUser] = struct{}{}
			ofUserProfile.MarkedAsSpamBy[b.Creator] = struct{}{}
		}
		if b.HasTag(object.BlockTag) {
			v.i.EnsureUsersOfUserVoteBody(b)
			creatorProfile.Blocked[b.OfUser] = struct{}{}
			ofUserProfile.BlockedBy[b.Creator] = struct{}{}
		}
	case 0:
		v.i.EnsureUsersOfUserVoteBody(b)
	}
	return nil
}

/*
	<<< GET >>>
*/

func (v *Viewer) GetBoard() (*object.ContentRep, error) {
	if v == nil {
		return nil, ErrViewerNotInitialized
	}
	defer v.lock()()
	return v.c.content[v.i.Board], nil
}

type BoardPageIn struct {
	Perspective    string
	PaginatedInput typ.PaginatedInput
}

type BoardPageOut struct {
	Board       *object.ContentRep   `json:"board"`
	//ThreadsMeta *typ.PaginatedOutput `json:"threads_meta"`
	Threads     []*object.ContentRep `json:"threads"`
}

func (v *Viewer) GetBoardPage(in *BoardPageIn) (*BoardPageOut, error) {
	if v == nil {
		return nil, ErrViewerNotInitialized
	}
	defer v.lock()()

	tHashes, e := v.i.Threads.Get(&in.PaginatedInput)
	if e != nil {
		return nil, e
	}

	out := new(BoardPageOut)
	out.Board = v.c.content[v.i.Board]
	//out.ThreadsMeta = tHashes
	out.Threads = make([]*object.ContentRep, len(tHashes.Data))
	for i, tHash := range tHashes.Data {
		out.Threads[i] = v.c.content[tHash]
		if votes, ok := v.c.votes[tHash]; ok {
			out.Threads[i].Votes = votes.View(in.Perspective)
		}
	}
	return out, nil
}

type ThreadPageIn struct {
	Perspective    string
	ThreadHash     string
	PaginatedInput typ.PaginatedInput
}

type ThreadPageOut struct {
	Board     *object.ContentRep   `json:"board"`
	Thread    *object.ContentRep   `json:"thread"`
	//PostsMeta *typ.PaginatedOutput `json:"posts_meta"`
	Posts     []*object.ContentRep `json:"posts"`
}

func (v *Viewer) GetThreadPage(in *ThreadPageIn) (*ThreadPageOut, error) {
	if v == nil {
		return nil, ErrViewerNotInitialized
	}
	defer v.lock()()
	out := new(ThreadPageOut)
	out.Board = v.c.content[v.i.Board]
	out.Thread = v.c.content[in.ThreadHash]

	if out.Thread == nil {
		return nil, boo.Newf(boo.NotFound, "thread of hash '%s' is not found in board '%s'",
			in.ThreadHash, v.pk.Hex())
	}
	if votes, ok := v.c.votes[in.ThreadHash]; ok {
		out.Thread.Votes = votes.View(in.Perspective)
	}

	pHashes, e := v.i.PostsOfThread[in.ThreadHash].Get(&in.PaginatedInput)
	if e != nil {
		return nil, e
	}
	out.Posts = make([]*object.ContentRep, len(pHashes.Data))
	for i, pHash := range pHashes.Data {
		out.Posts[i] = v.c.content[pHash]
		if votes, ok := v.c.votes[pHash]; ok {
			out.Posts[i].Votes = votes.View(in.Perspective)
		}
	}

	return out, nil
}

type ContentVotesIn struct {
	Perspective string
	ContentHash string
}

type ContentVotesOut struct {
	Votes *VoteRepView `json:"votes"`
}

func (v *Viewer) GetVotes(in *ContentVotesIn) (*ContentVotesOut, error) {
	if v == nil {
		return nil, ErrViewerNotInitialized
	}
	defer v.lock()()
	out := new(ContentVotesOut)
	if votes, ok := v.c.votes[in.ContentHash]; ok {
		out.Votes = votes.View(in.Perspective)
		return out, nil
	}
	if _, ok := v.c.content[in.ContentHash]; ok {
		out.Votes = &VoteRepView{
			Ref: in.ContentHash,
		}
		return out, nil
	}
	return nil, boo.Newf(boo.NotFound, "content of hash '%s' is not found",
		in.ContentHash)
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
