package state

import (
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/skycoin/src/cipher"
)

var (
	// ErrEmpty occurs when object is nil.
	ErrEmpty = boo.New(boo.ObjectNotFound, "nil error")
)

func corruptWrap(e error) error {
	return boo.WrapType(e, boo.InvalidRead, "corrupt user file")
}

// UserFile represents a user of user configuration.
type UserFile struct {
	User          obj.User           `json:"user"`
	Subscriptions []obj.Subscription `json:"subscriptions"`
	Masters       []obj.Subscription `json:"masters"`
}

// Check ensures the validity of the UserFile.
func (f *UserFile) Check() error {
	if f == nil {
		return ErrEmpty
	}
	if e := f.User.PublicKey.Verify(); e != nil {
		return corruptWrap(e)
	}
	if e := f.User.SecretKey.Verify(); e != nil {
		return corruptWrap(e)
	}
	for i, sub := range f.Subscriptions {
		if e := sub.PubKey.Verify(); e != nil {
			return corruptWrap(e)
		}
		f.Subscriptions[i].SecKey = cipher.SecKey{}
	}
	for _, sub := range f.Masters {
		if e := sub.PubKey.Verify(); e != nil {
			return corruptWrap(e)
		}
		if e := sub.SecKey.Verify(); e != nil {
			return corruptWrap(e)
		}
	}
	return nil
}

// GenerateView generates something readable for front end.
func (f *UserFile) GenerateView() *UserFileView {
	if e := f.Check(); e != nil {
		return &UserFileView{}
	}
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

// UserFileView represents a user user as displayed to end user.
type UserFileView struct {
	User          obj.UserView           `json:"user"`
	Subscriptions []obj.SubscriptionView `json:"subscriptions"`
	Masters       []obj.SubscriptionView `json:"subscriptions"`
}
