package rpc

import (
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/bbs/src/store/object"
)

var (
	ErrEmptyInput = boo.New(boo.InvalidInput, "empty input error")
)

type Gateway struct {
	CXO      *cxo.Manager
	QuitChan chan int
}

func (g *Gateway) Quit(code int, ok *bool) error {
	if ok == nil {
		return ErrEmptyInput
	}
	g.QuitChan <- code
	*ok = true
	return nil
}

type ConnectionsOutput struct {
	Connections []r0.Connection `json:"connections"`
}

func (g *Gateway) Connections(_ struct{}, out *ConnectionsOutput) error {
	out.Connections = g.CXO.GetConnections()
	return nil
}

func (g *Gateway) AddConnection(address string, _ *struct{}) error {
	return g.CXO.Connect(address)
}

func (g *Gateway) RemoveConnection(address string, _ *struct{}) error {
	return g.CXO.Disconnect(address)
}

type FeedsOutput struct {
	Feeds []string `json:"feeds"`
}

func (g *Gateway) Subscriptions(_ struct{}, out *FeedsOutput) error {
	pks := g.CXO.GetSubscriptions()
	out.Feeds = make([]string, len(pks))
	for i, pk := range pks {
		out.Feeds[i] = pk.Hex()
	}
	return nil
}

func (g *Gateway) SubscribeRemote(pk cipher.PubKey, _ *struct{}) error {
	return g.CXO.SubscribeRemote(pk)
}

type SubscribeMasterInput struct {
	PubKey cipher.PubKey
	SecKey cipher.SecKey
}

func (g *Gateway) SubscribeMaster(in SubscribeMasterInput, _ *struct{}) error {
	return g.CXO.SubscribeMaster(in.PubKey, in.SecKey)
}

func (g *Gateway) NewBoard(in *object.NewBoardIO, _ *struct{}) error {
	if e := in.Process(g.CXO.Relay().GetKeys()); e != nil {
		return e
	}
	return g.CXO.NewBoard(in)
}

func (g *Gateway) DeleteBoard(in *object.BoardIO, _ *struct{}) error {
	if e := in.Process(); e != nil {
		return e
	}
	return g.CXO.UnsubscribeMaster(in.PubKey)
}

func (g *Gateway) ExportBoard(in *object.ExportBoardIO, _ *struct{}) error {
	if e := in.Process(); e != nil {
		return e
	}
	_, _, e := g.CXO.ExportBoard(in.PubKey, in.Name)
	return e
}