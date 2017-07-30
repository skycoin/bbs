package content

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"sync"
)

type Result struct {
	e                    error
	root                 *node.Root
	BoardPage            *object.BoardPage
	Board                *object.Board
	ThreadPage           *object.ThreadPage
	ThreadPages          []*object.ThreadPage
	Thread               *object.Thread
	Threads              []*object.Thread
	ThreadVotesPage      *object.ThreadVotesPage
	ThreadVotesPageIndex int
	ThreadVote           *object.Vote
	ThreadIndex          int
	ThreadRefMap         map[cipher.SHA256]int
	Post                 *object.Post
	Posts                []*object.Post
	PostVotesPage        *object.PostVotesPage
	PostVotesPageIndex   int
	PostVote             *object.Vote
	PostIndex            int
	PostRefMap           map[cipher.SHA256]int
}

func NewResult(root *node.Root) *Result {
	if len(root.Refs()) != 3 {
		return &Result{e: boo.New(boo.InvalidRead,
			"corrupt board - reference count is not 3")}
	}
	return &Result{
		root:                 root,
		ThreadIndex:          -1,
		PostIndex:            -1,
		ThreadVotesPageIndex: -1,
		PostVotesPageIndex:   -1,
	}
}

func (r *Result) Error() error {
	return r.e
}

func (r *Result) GetPK() cipher.PubKey {
	return r.root.Pub()
}

func (r *Result) GetSeq() uint64 {
	return r.root.Seq()
}

func (r *Result) GetPages(b, t, p bool) *Result {
	if r.e != nil {
		return r
	}
	var wg sync.WaitGroup
	if b {
		r.BoardPage = &object.BoardPage{
			R: toSHA256(r.root.Refs()[0].Object),
		}
		if e := r.deserialize(toRef(r.BoardPage.R), r.BoardPage); e != nil {
			r.e = boo.WrapType(e, boo.InvalidRead, "invalid board page")
			return r
		}
	}
	if t {
		r.ThreadVotesPage = &object.ThreadVotesPage{
			R: toSHA256(r.root.Refs()[1].Object),
		}
		if e := r.deserialize(toRef(r.ThreadVotesPage.R), r.ThreadVotesPage); e != nil {
			r.e = boo.WrapType(e, boo.InvalidRead, "invalid thread votes page")
			return r
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.ThreadRefMap = make(map[cipher.SHA256]int)
			for i, tvp := range r.ThreadVotesPage.Store {
				r.ThreadRefMap[tvp.Ref] = i
			}
		}()
	}
	if p {
		r.PostVotesPage = &object.PostVotesPage{
			R: toSHA256(r.root.Refs()[2].Object),
		}
		if e := r.deserialize(toRef(r.PostVotesPage.R), r.PostVotesPage); e != nil {
			r.e = boo.WrapType(e, boo.InvalidRead, "invalid post votes page")
			return r
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.PostRefMap = make(map[cipher.SHA256]int)
			for i, pvp := range r.PostVotesPage.Store {
				r.PostRefMap[pvp.Ref] = i
			}
		}()
	}
	wg.Wait()
	return r
}

func (r *Result) GetBoard() *Result {
	if r.e != nil {
		return r
	}
	r.Board = &object.Board{
		R: toSHA256(r.BoardPage.Board),
	}
	if e := r.deserialize(toRef(r.Board.R), r.Board); e != nil {
		r.e = boo.Wrap(e, "invalid board")
		return r
	}
	return r
}

func (r *Result) GetThreadPage(tRef skyobject.Reference) *Result {
	if r.e != nil {
		return r
	}
	for i, tpRef := range r.BoardPage.ThreadPages {
		var tp object.ThreadPage
		if e := r.deserialize(tpRef, &tp); e != nil {
			r.e = boo.WrapTypef(e, boo.InvalidRead,
				"thread page %s of board %s is corrupt", tpRef.String(), r.Board.R.Hex())
			return r
		}
		if tp.Thread == tRef {
			r.ThreadIndex = i
			r.ThreadPage = &tp
			r.ThreadPage.R = toSHA256(tpRef)
			return r
		}
	}
	r.e = boo.Newf(boo.NotFound,
		"thread of reference %s is not found under board %s", tRef.String(), r.Board.R.Hex())
	return r
}

func (r *Result) GetThread() *Result {
	if r.e != nil {
		return r
	}
	r.Thread = &object.Thread{
		R: toSHA256(r.ThreadPage.Thread),
	}
	if e := r.deserialize(r.ThreadPage.Thread, r.Thread); e != nil {
		r.e = boo.WrapTypef(e, boo.InvalidRead,
			"thread object of reference %s is corrupt", r.Thread.R.Hex())
		return r
	}
	return r
}

func (r *Result) GetThreadPages() *Result {
	if r.e != nil {
		return r
	}
	refs := r.BoardPage.ThreadPages
	r.ThreadPages = make([]*object.ThreadPage, len(refs))
	for i, ref := range refs {
		r.ThreadPages[i] = &object.ThreadPage{
			R: toSHA256(ref),
		}
		if e := r.deserialize(ref, r.ThreadPages[i]); e != nil {
			r.e = boo.WrapTypef(e, boo.InvalidRead,
				"thread page %s of board %s is corrupt", ref.String(), r.Board.R.Hex())
			return r
		}
	}
	return r
}

func (r *Result) GetThreads() *Result {
	if r.e != nil {
		return r
	}
	r.Threads = make([]*object.Thread, len(r.ThreadPages))
	for i, tPage := range r.ThreadPages {
		r.Threads[i] = &object.Thread{
			R: toSHA256(tPage.Thread),
		}
		if e := r.deserialize(tPage.Thread, r.Threads[i]); e != nil {
			r.e = boo.WrapTypef(e, boo.InvalidRead,
				"thread %s of board %s is corrupt",
				tPage.Thread.String(), r.Board.R.Hex())
			return r
		}
	}
	return r
}

func (r *Result) GetPosts() *Result {
	if r.e != nil {
		return r
	}
	r.Posts = make([]*object.Post, len(r.ThreadPage.Posts))
	for i, pRef := range r.ThreadPage.Posts {
		r.Posts[i] = &object.Post{
			R: toSHA256(pRef),
		}
		if e := r.deserialize(pRef, r.Posts[i]); e != nil {
			r.e = boo.WrapTypef(e, boo.InvalidRead,
				"post %s of thread %s of board %s is corrupt",
				pRef.String(), r.Thread.R.Hex(), r.Board.R.Hex())
			return r
		}
	}
	return r
}

func (r *Result) savePost() *Result {
	if r.e != nil {
		return r
	}
	r.Post.R = toSHA256(r.root.Save(r.Post))
	if _, has := r.PostRefMap[r.Post.R]; has {
		r.e = boo.Newf(boo.AlreadyExists,
			"this exact post of reference %s already exists in board %s",
			r.Post.R.Hex(), r.root.Pub().Hex())
		return r
	} else {
		r.PostVotesPage.Store = append(
			r.PostVotesPage.Store,
			object.VotesPage{Ref: r.Post.R},
		)
	}
	if r.PostIndex == -1 {
		r.ThreadPage.Posts = append(
			r.ThreadPage.Posts, toRef(r.Post.R))
	} else {
		r.ThreadPage.Posts[r.PostIndex] =
			toRef(r.Post.R)
	}
	return r
}

func (r *Result) saveThread() *Result {
	if r.e != nil {
		return r
	}
	r.Thread.R = toSHA256(r.root.Save(r.Thread))
	if _, has := r.ThreadRefMap[r.Thread.R]; has {
		r.e = boo.Newf(boo.AlreadyExists,
			"this exact thread of reference %s already exists in board %s",
			r.Thread.R.Hex(), r.root.Pub().Hex())
		return r
	} else {
		r.ThreadVotesPage.Store = append(
			r.ThreadVotesPage.Store,
			object.VotesPage{Ref: r.Thread.R},
		)
	}
	if r.ThreadPage == nil {
		r.ThreadPage = new(object.ThreadPage)
	}
	r.ThreadPage.Thread = toRef(r.Thread.R)
	return r
}

func (r *Result) saveThreadPage() *Result {
	if r.e != nil {
		return r
	}
	r.ThreadPage.R = toSHA256(r.root.Save(r.ThreadPage))
	if r.ThreadIndex == -1 {
		r.BoardPage.ThreadPages = append(
			r.BoardPage.ThreadPages, toRef(r.ThreadPage.R))
	} else {
		r.BoardPage.ThreadPages[r.ThreadIndex] =
			toRef(r.ThreadPage.R)
	}
	return r
}

func (r *Result) saveBoard() *Result {
	if r.e != nil {
		return r
	}
	r.Board.R = toSHA256(r.root.Save(r.Board))
	r.BoardPage.Board = toRef(r.Board.R)
	return r
}

func (r *Result) savePages(b, t, p bool) *Result {
	if r.e != nil {
		return r
	}
	refs := r.root.Refs()
	if b {
		r.BoardPage.R = toSHA256(r.root.Save(r.BoardPage))
		refs[0].Object = toRef(r.BoardPage.R)
	}
	if t {
		r.ThreadVotesPage.R = toSHA256(r.root.Save(r.ThreadVotesPage))
		refs[1].Object = toRef(r.ThreadVotesPage.R)
	}
	if p {
		r.PostVotesPage.R = toSHA256(r.root.Save(r.PostVotesPage))
		refs[2].Object = toRef(r.PostVotesPage.R)
	}
	if _, e := r.root.Replace(refs); e != nil {
		r.e = e
	}
	return r
}

func (r *Result) saveThreadVote(tRef skyobject.Reference) *Result {
	if r.e != nil {
		return r
	}
	r.ThreadVote.R = toSHA256(r.root.Save(r.ThreadVote))

	tvi, has := r.ThreadRefMap[toSHA256(tRef)]
	if !has {
		r.e = boo.Newf(boo.NotFound,
			"thread of reference %s not found in board %s",
			tRef.String(), r.root.Pub().Hex())
		return r
	}

	var temp object.Vote
	for i, vRef := range r.ThreadVotesPage.Store[tvi].Votes {
		if e := r.deserialize(vRef, &temp); e != nil {
			r.e = boo.WrapTypef(e, boo.InvalidRead,
				"vote %d from thread %s of board %s is corrupt",
				i, tRef.String(), r.root.Pub().Hex())
			return r
		}
		if temp.User == r.ThreadVote.User {
			r.ThreadVotesPage.Store[tvi].Votes[i] =
				toRef(r.ThreadVote.R)
			return r
		}
	}
	r.ThreadVotesPage.Store[tvi].Votes = append(
		r.ThreadVotesPage.Store[tvi].Votes,
		toRef(r.ThreadVote.R),
	)
	return r
}

func (r *Result) savePostVote(pRef skyobject.Reference) *Result {
	if r.e != nil {
		return r
	}
	r.PostVote.R = toSHA256(r.root.Save(r.PostVote))

	pvi, has := r.PostRefMap[toSHA256(pRef)]
	if !has {
		r.e = boo.Newf(boo.NotFound,
			"post of reference %s not found in board %s",
			pRef.String(), r.root.Pub().Hex())
		return r
	}

	var temp object.Vote
	for i, vRef := range r.PostVotesPage.Store[pvi].Votes {
		if e := r.deserialize(vRef, &temp); e != nil {
			r.e = boo.WrapTypef(e, boo.InvalidRead,
				"vote %d from thread %s of board %s is corrupt",
				pRef.String(), r.root.Pub().Hex())
			return r
		}
		if temp.User == r.PostVote.User {
			r.PostVotesPage.Store[pvi].Votes[i] =
				toRef(r.PostVote.R)
			return r
		}
	}
	r.PostVotesPage.Store[pvi].Votes = append(
		r.PostVotesPage.Store[pvi].Votes,
		toRef(r.PostVote.R),
	)
	return r
}

func (r *Result) deleteThreadVote(tRef skyobject.Reference) *Result {
	if r.e != nil {
		return r
	}
	di := r.ThreadRefMap[toSHA256(tRef)]
	r.ThreadVotesPage.Store = append(
		r.ThreadVotesPage.Store[:di],
		r.ThreadVotesPage.Store[di+1:]...,
	)
	r.ThreadVotesPage.Deleted = append(
		r.ThreadVotesPage.Deleted,
		toSHA256(tRef),
	)
	return r
}

func (r *Result) deletePostVote(pRef skyobject.Reference) *Result {
	if r.e != nil {
		return r
	}
	di := r.PostRefMap[toSHA256(pRef)]
	r.PostVotesPage.Store = append(
		r.PostVotesPage.Store[:di],
		r.PostVotesPage.Store[di+1:]...,
	)
	r.PostVotesPage.Deleted = append(
		r.PostVotesPage.Deleted,
		toSHA256(pRef),
	)
	return r
}

func (r *Result) deletePostVotes(pRefs skyobject.References) *Result {
	if r.e != nil {
		return r
	}
	for _, pRef := range pRefs {
		di := r.PostRefMap[toSHA256(pRef)]
		r.PostVotesPage.Store = append(
			[]object.VotesPage{{}},
			append(
				r.PostVotesPage.Store[:di],
				r.PostVotesPage.Store[di+1:]...,
			)...,
		)
		r.PostVotesPage.Deleted = append(
			r.PostVotesPage.Deleted,
			toSHA256(pRef),
		)
	}
	r.PostVotesPage.Store =
		r.PostVotesPage.Store[len(pRefs):]
	return r
}

func (r *Result) deleteThread(i int) *Result {
	if r.e != nil {
		return r
	}
	r.BoardPage.ThreadPages = append(
		r.BoardPage.ThreadPages[:i],
		r.BoardPage.ThreadPages[i+1:]...,
	)
	r.ThreadPages = append(
		r.ThreadPages[:i],
		r.ThreadPages[i+1:]...,
	)
	r.Threads = append(
		r.Threads[:i],
		r.Threads[i+1:]...,
	)
	return r
}

func (r *Result) deletePost(i int) *Result {
	if r.e != nil {
		return r
	}
	r.ThreadPage.Posts = append(
		r.ThreadPage.Posts[:i],
		r.ThreadPage.Posts[i+1:]...,
	)
	r.Posts = append(
		r.Posts[:i],
		r.Posts[i+1:]...,
	)
	return r
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func toRef(sha256 cipher.SHA256) skyobject.Reference {
	return skyobject.Reference(sha256)
}

func toSHA256(reference skyobject.Reference) cipher.SHA256 {
	return cipher.SHA256(reference)
}

func (r *Result) deserialize(ref skyobject.Reference, v interface{}) error {
	data, _ := r.root.Get(ref)
	if e := encoder.DeserializeRaw(data, v); e != nil {
		return boo.WrapType(e, boo.InvalidRead, "corrupt board")
	}
	return nil
}
