package object

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type SubscriptionView struct {
	PK string `json:"public_key"`
	SK string `json:"secret_key,omitempty"`
}

type Subscription struct {
	PK cipher.PubKey
	SK cipher.SecKey
}

func (s *Subscription) View() SubscriptionView {
	if s.SK == (cipher.SecKey{}) {
		return SubscriptionView{
			PK: s.PK.Hex(),
		}
	} else {
		return SubscriptionView{
			PK: s.PK.Hex(),
			SK: s.SK.Hex(),
		}
	}
}

type CXOFile struct {
	MasterSubs []SubscriptionView `json:"master_subscriptions"`
	RemoteSubs []SubscriptionView `json:"remote_subscriptions"`
}
