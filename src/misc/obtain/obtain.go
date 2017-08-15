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
	content.R = ref.Hash
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

func ThreadPage(ref *skyobject.Ref) (*object.ThreadPage, error) {
	value, e := ref.Value()
	if e != nil {
		return nil, valErr(e, ref.Hash)
	}
	tPage, ok := value.(*object.ThreadPage)
	if !ok {
		return nil, extErr("ThreadPage", ref.Hash)
	}
	return tPage, nil
}

func ThreadPages(ref *skyobject.Ref) (*object.ThreadPages, error) {
	value, e := ref.Value()
	if e != nil {
		return nil, valErr(e, ref.Hash)
	}
	tPages, ok := value.(*object.ThreadPages)
	if !ok {
		return nil, extErr("ThreadPages", ref.Hash)
	}
	return tPages, nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

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
