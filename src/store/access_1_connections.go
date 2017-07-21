package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/bbs/src/store/state"
)

type ConnectionsOutput struct {
	Connections []obj.ConnectionView `json:"connections"`
}

func generateConnectionsOutput(cxo *state.CXO, file *state.UserFile) (*ConnectionsOutput, error) {
	actives, e := cxo.GetConnections()
	if e != nil {
		return nil, e
	}
	activeMap := make(map[string]bool)
	for _, address := range actives {
		activeMap[address] = true
	}

	out := new(ConnectionsOutput)
	for _, address := range file.Connections {
		out.Connections = append(out.Connections, obj.ConnectionView{
			Address: address,
			Active:  activeMap[address],
		})
	}

	return out, nil
}

func (a *Access) GetConnections(ctx context.Context) (*ConnectionsOutput, error) {
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	return generateConnectionsOutput(a.Session.GetCXO(), file)
}

func (a *Access) NewConnection(ctx context.Context, in *state.ConnectionIO) (*ConnectionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	file, e := a.Session.NewConnection(ctx, in)
	if e != nil {
		return nil, e
	}
	return generateConnectionsOutput(a.Session.GetCXO(), file)
}

func (a *Access) DeleteConnection(ctx context.Context, in *state.ConnectionIO) (*ConnectionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	file, e := a.Session.DeleteConnection(ctx, in)
	if e != nil {
		return nil, e
	}
	return generateConnectionsOutput(a.Session.GetCXO(), file)
}
