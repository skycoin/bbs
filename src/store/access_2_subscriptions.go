package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/bbs/src/store/state"
)

type SubsOutput struct {
	Subscriptions       []obj.SubscriptionView `json:"subscriptions"`
	MasterSubscriptions []obj.SubscriptionView `json:"master_subscriptions"`
}

func (a *Access) GetSubs(ctx context.Context) (*SubsOutput, error) {
	view, e := a.GetSession(ctx)
	if e != nil {
		return nil, e
	}
	out := &SubsOutput{
		Subscriptions:       view.Subscriptions,
		MasterSubscriptions: view.Masters,
	}
	return out, nil
}

func (a *Access) NewSub(ctx context.Context, in *state.SubscriptionIO) (*SubsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	file, e := a.Session.NewSubscription(ctx, in)
	if e != nil {
		return nil, e
	}
	view := file.GenerateView(a.Session.GetCXO())
	out := &SubsOutput{
		Subscriptions:       view.Subscriptions,
		MasterSubscriptions: view.Masters,
	}
	return out, nil
}

func (a *Access) DeleteSub(ctx context.Context, in *state.SubscriptionIO) (*SubsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	file, e := a.Session.DeleteSubscription(ctx, in)
	if e != nil {
		return nil, e
	}
	view := file.GenerateView(a.Session.GetCXO())
	out := &SubsOutput{
		Subscriptions:       view.Subscriptions,
		MasterSubscriptions: view.Masters,
	}
	return out, nil
}
