package content

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

type Result struct {
	e               error
	root            *node.Root
	BoardPage       *object.BoardPage
	Board           *object.Board
	ThreadPage      *object.ThreadPage
	ThreadPageIndex int
	ThreadPages     []*object.ThreadPage
	Thread          *object.Thread
	Threads         []*object.Thread
	Posts           []*object.Post
}

func NewResult(cxo *state.CXO, pk cipher.PubKey, sk ...cipher.SecKey) *Result {
	root, e := cxo.GetRoot(pk)
	if e != nil {
		return &Result{e: boo.WrapType(e, boo.Internal,
			"failed to obtain board root")}
	}
	if root == nil {
		return &Result{e: boo.New(boo.NotFound,
			"this board is not yet downloaded or does not exist")}
	}
	if len(root.Refs()) != 3 {
		return &Result{e: boo.New(boo.InvalidRead,
			"corrupt board - reference count is not 3")}
	}
	if len(sk) == 1 {
		root.Edit(sk[0])
	}
	return &Result{root: root, ThreadPageIndex: -1}
}

func (r *Result) Error() error {
	return r.e
}

func (r *Result) GetPK() cipher.PubKey {
	return r.root.Pub()
}

func (r *Result) getBoardPage() *Result {
	if r.e != nil {
		return r
	}
	r.BoardPage = &object.BoardPage{
		R: toSHA256(r.root.Refs()[0].Object),
	}
	if e := r.deserialize(toRef(r.BoardPage.R), r.BoardPage); e != nil {
		r.e = boo.Wrap(e, "invalid board page")
		return r
	}
	return r
}

func (r *Result) getBoard() *Result {
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

func (r *Result) getThreadPage(tRef skyobject.Reference) *Result {
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
			r.ThreadPageIndex = i
			r.ThreadPage = &tp
			r.ThreadPage.R = toSHA256(tpRef)
			return r
		}
	}
	r.e = boo.Newf(boo.NotFound,
		"thread of reference %s is not found under board %s", tRef.String(), r.Board.R.Hex())
	return r
}

func (r *Result) getThread() *Result {
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

func (r *Result) getThreadPages() *Result {
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

func (r *Result) getThreads() *Result {
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
				"thread %s of board %s is corrupt", tPage.Thread.String(), r.Board.R.Hex())
			return r
		}
	}
	return r
}

func (r *Result) saveThread() *Result {
	if r.e != nil {
		return r
	}
	r.Thread.R = toSHA256(r.root.Save(r.Thread))
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
	r.BoardPage.ThreadPages = append(
		r.BoardPage.ThreadPages, toRef(r.ThreadPage.R))
	return r
}

func (r *Result) saveThreadPages() *Result {
	if r.e != nil {
		return r
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

func (r *Result) saveBoardPage() *Result {
	if r.e != nil {
		return r
	}
	r.BoardPage.R = toSHA256(r.root.Save(r.BoardPage))
	refs := r.root.Refs()
	refs[0].Object = toRef(r.BoardPage.R)
	if _, e := r.root.Replace(refs); e != nil {
		r.e = e
	}
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
