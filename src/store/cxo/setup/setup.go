package setup

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
)

// PrepareRegistry sets up the CXO Registry.
func PrepareRegistry(r *skyobject.Reg) {
	r.Register(
		object.RootPageName,
		object.RootPage{})

	r.Register(
		object.BoardPageName,
		object.BoardPage{})

	r.Register(
		object.ThreadPageName,
		object.ThreadPage{})

	r.Register(
		object.DiffPageName,
		object.DiffPage{})

	r.Register(
		object.UsersPageName,
		object.UsersPage{})

	r.Register(
		object.UserProfileName,
		object.UserProfile{})

	r.Register(
		object.ContentName,
		object.Content{})
}

// NewBoard generates a new board.
func NewBoard(node *node.Node, in *object.NewBoardIO) (*skyobject.Root, error) {
	pack, e := node.Container().NewRoot(
		in.BoardPubKey,
		in.BoardSecKey,
		skyobject.HashTableIndex|skyobject.EntireTree,
		node.Container().CoreRegistry().Types(),
	)
	if e != nil {
		return nil, e
	}

	if e := SetBoard(pack, in); e != nil {
		return nil, e
	}
	node.Publish(pack.Root())
	pack.Close()

	return node.Container().LastRoot(in.BoardPubKey)
}

func SetBoard(pack *skyobject.Pack, in *object.NewBoardIO) error {
	pack.Clear()
	pack.Append(
		&object.RootPage{
			Typ: object.RootTypeBoard,
			Rev: 0,
			Del: false,
			Sum: in.Content.Body,
		},
		&object.BoardPage{
			Board: pack.Ref(in.Content),
		},
		&object.DiffPage{},
		&object.UsersPage{},
	)
	return pack.Save()
}
