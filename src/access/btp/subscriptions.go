package btp

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc"
)

// GetSubscriptionsView obtains the BoardFile as json.
func (a *BoardAccessor) GetSubscriptionsView() ([]byte, error) {
	defer a.lock()()
	return json.Marshal(*a.bFile)
}

// NewSubscriptionInput provides input for NewSubscription.
type NewSubscriptionInput struct {
	Address string `json:"address"`
	PubKey  string `json:"public_key"`
}

// NewSubscription creates a new subscription.
func (a *BoardAccessor) NewSubscription(in *NewSubscriptionInput) error {
	defer a.lock()()
	pk, e := misc.GetPubKey(in.PubKey)
	if e != nil {
		return e
	}
	if e := a.cxo.Subscribe(in.Address, pk); e != nil {
		return e
	}
	if e := a.bFile.Add(pk, in.Address); e != nil {
		return e
	}
	return nil
}

// RemoveSubscriptionInput provides input to remove subscription.
type RemoveSubscriptionInput struct {
	PubKey string `json:"public_key"`
}

// RemoveSubscription removes a subscription.
func (a *BoardAccessor) RemoveSubscription(in *RemoveSubscriptionInput) error {
	defer a.lock()()
	pk, e := misc.GetPubKey(in.PubKey)
	if e != nil {
		return e
	}
	a.cxo.Unsubscribe("", pk)
	a.bFile.Remove(pk)
	a.stateSaver.Remove(pk)
	return nil
}
