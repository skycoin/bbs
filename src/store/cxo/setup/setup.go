package setup

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/node/log"
)

func NewCXOConfig() node.Config {
	var sc node.Config
	sc.Config = gnet.NewConfig()
	sc.Log = log.NewConfig()
	sc.Skyobject = skyobject.NewConfig()
	sc.EnableRPC = node.EnableRPC
	sc.RPCAddress = node.RPCAddress
	sc.Listen = node.Listen
	sc.EnableListener = node.EnableListener
	sc.RemoteClose = node.RemoteClose
	sc.PingInterval = node.PingInterval
	sc.InMemoryDB = node.InMemoryDB
	//sc.DataDir = node.DataDir()
	sc.DBPath = ""
	sc.ResponseTimeout = node.ResponseTimeout
	sc.PublicServer = node.PublicServer
	sc.Config.OnDial = node.OnDialFilter
	return sc
}

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
func NewBoard(node *node.Node, content *object.Content, pk cipher.PubKey, sk cipher.SecKey) (*skyobject.Root, error) {
	pack, e := node.Container().NewRoot(
		pk,
		sk,
		skyobject.HashTableIndex|skyobject.EntireTree,
		node.Container().CoreRegistry().Types(),
	)
	if e != nil {
		return nil, e
	}

	if e := SetBoard(pack, content); e != nil {
		return nil, e
	}
	node.Publish(pack.Root())
	pack.Close()

	return node.Container().LastRoot(pk)
}

func SetBoard(pack *skyobject.Pack, content *object.Content) error {
	pack.Clear()
	pack.Append(
		&object.RootPage{
			Typ: object.RootTypeBoard,
			Rev: 0,
			Del: false,
			Sum: content.Body,
		},
		&object.BoardPage{
			Board: pack.Ref(content),
		},
		&object.DiffPage{},
		&object.UsersPage{},
	)
	return pack.Save()
}
