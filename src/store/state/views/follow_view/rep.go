package follow_view

import "github.com/skycoin/skycoin/src/cipher"

type UserRep struct {
	UserPubKey string `json:"user_public_key"`
	Tag        string `json:"tag"`
}

type FollowRep struct {
	UserPubKey string                     `json:"user_public_key"`
	Following  map[cipher.PubKey]*UserRep `json:"following"`
	Avoiding   map[cipher.PubKey]*UserRep `json:"avoiding"`
}

func NewFollowRep(upk cipher.PubKey) *FollowRep {
	return &FollowRep{
		UserPubKey: upk.Hex(),
		Following:  make(map[cipher.PubKey]*UserRep),
		Avoiding:   make(map[cipher.PubKey]*UserRep),
	}
}

func (r *FollowRep) Set(upk cipher.PubKey, mode int8, tag []byte) {
	// Remove existing.
	delete(r.Following, upk)
	delete(r.Avoiding, upk)

	// Add in.
	switch mode {
	case +1:
		r.Following[upk] = &UserRep{
			UserPubKey: upk.Hex(),
			Tag:        string(tag),
		}
	case -1:
		r.Avoiding[upk] = &UserRep{
			UserPubKey: upk.Hex(),
			Tag:        string(tag),
		}
	}
}

type FollowRepView struct {
	UserPubKey string     `json:"user_public_key"`
	Following  []*UserRep `json:"following"`
	Avoiding   []*UserRep `json:"avoiding"`
}

func (r *FollowRep) View() *FollowRepView {
	view := &FollowRepView{
		UserPubKey: r.UserPubKey,
		Following:  make([]*UserRep, len(r.Following)),
		Avoiding:   make([]*UserRep, len(r.Avoiding)),
	}
	i := 0
	for _, u := range r.Following {
		view.Following[i] = u
		i++
	}
	i = 0
	for _, u := range r.Avoiding {
		view.Avoiding[i] = u
		i++
	}
	return view
}
