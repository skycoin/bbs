package state

import "github.com/skycoin/bbs/src/store/obj"

// UserFile represents a file of user configuration.
type UserFile struct {
	User          obj.User           `json:"user"`
	Subscriptions []obj.Subscription `json:"subscriptions"`
	Masters       []obj.Subscription `json:"masters"`
}

// GenerateView generates something readable for front end.
func (f *UserFile) GenerateView() *UserFileView {
	view := &UserFileView{
		User: obj.UserView{
			Alias:     f.User.Alias,
			PublicKey: f.User.PublicKey.Hex(),
			SecretKey: f.User.SecretKey.Hex(),
		},
	}

	subscriptions := make([]obj.SubscriptionView, len(f.Subscriptions))
	for i, s := range f.Subscriptions {
		subscriptions[i] = obj.SubscriptionView{
			PubKey:      s.PubKey.Hex(),
			SecKey:      s.SecKey.Hex(),
			Connections: s.Connections,
		}
	}
	view.Subscriptions = subscriptions

	masters := make([]obj.SubscriptionView, len(f.Masters))
	for i, m := range f.Masters {
		masters[i] = obj.SubscriptionView{
			PubKey:      m.PubKey.Hex(),
			SecKey:      m.SecKey.Hex(),
			Connections: m.Connections,
		}
	}
	view.Masters = masters

	return view
}

type UserFileView struct {
	User          obj.UserView           `json:"user"`
	Subscriptions []obj.SubscriptionView `json:"subscriptions"`
	Masters       []obj.SubscriptionView `json:"subscriptions"`
}
