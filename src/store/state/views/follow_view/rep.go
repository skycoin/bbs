package follow_view

type UserRep struct {
	UserPubKey string `json:"user_public_key"`
	Tag        string `json:"tag"`
}

type FollowRep struct {
	UserPubKey string              `json:"user_public_key"`
	Following  map[string]*UserRep `json:"following"` // key(user public key)
	Avoiding   map[string]*UserRep `json:"avoiding"`  // key(user public key)
}

func NewFollowRep(upk string) *FollowRep {
	return &FollowRep{
		UserPubKey: upk,
		Following:  make(map[string]*UserRep),
		Avoiding:   make(map[string]*UserRep),
	}
}

func (r *FollowRep) Set(upk string, mode int, tag string) {
	// Remove existing.
	delete(r.Following, upk)
	delete(r.Avoiding, upk)

	// Add in.
	switch mode {
	case +1:
		r.Following[upk] = &UserRep{
			UserPubKey: upk,
			Tag:        tag,
		}
	case -1:
		r.Avoiding[upk] = &UserRep{
			UserPubKey: upk,
			Tag:        tag,
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
