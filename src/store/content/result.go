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
	e           error
	root        *node.Root
	BoardPage   *object.BoardPage
	Board       *object.Board
	ThreadPages []*object.ThreadPage
	Posts       []*object.Post
}

func NewResult(cxo *state.CXO, pk cipher.PubKey, sk ...cipher.SecKey) *Result {
	root, e := cxo.GetRoot(pk)
	if e != nil {
		return &Result{e: boo.WrapType(e, boo.NotFound,
			"this board is not yet downloaded or does not exist")}
	}
	if len(root.Refs()) != 3 {
		return &Result{e: boo.New(boo.InvalidRead,
			"corrupt board - reference count is not 3")}
	}
	if len(sk) == 1 {
		root.Edit(sk[0])
	}
	return &Result{root: root}
}

func (r *Result) Error() error {
	return r.e
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
