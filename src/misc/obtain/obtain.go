package obtain

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

func Content(ref *skyobject.Ref) (*object.Content, error) {
	value, e := ref.Value()
	if e != nil {
		return nil, valErr(e, ref.Hash)
	}
	content, ok := value.(*object.Content)
	if !ok {
		return nil, extErr("Content", ref.Hash)
	}
	return content, nil
}

func Vote(ref *skyobject.Ref) (*object.Vote, error) {
	value, e := ref.Value()
	if e != nil {
		return nil, valErr(e, ref.Hash)
	}
	vote, ok := value.(*object.Vote)
	if !ok {
		return nil, extErr("Vote", ref.Hash)
	}
	return vote, nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func savePackErr(e error, what string) error {
	return boo.WrapTypef(e, boo.NotAllowed,
		"failed to save '%s'", what)
}

func valIndexErr(e error, i int) error {
	return boo.WrapType(e, boo.InvalidRead,
		"failed to obtain root child value of index %d", i)
}

func extIndexErr(name string, i int) error {
	return boo.Newf(boo.InvalidRead,
		"failed to extract '%s' from root child of index %d",
		name, i)
}

func valErr(e error, ref cipher.SHA256) error {
	return boo.WrapTypef(e, boo.InvalidRead,
		"failed to obtain value from object of ref '%s'",
		ref.Hex())
}

func extErr(name string, ref cipher.SHA256) error {
	return boo.Newf(boo.InvalidRead,
		"failed to extract '%s' from value of ref '%s'",
		name, ref.Hex())
}
