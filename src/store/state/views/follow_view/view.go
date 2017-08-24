package follow_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type FollowView struct {
	uMap map[cipher.PubKey]*FollowRep
}

func (v *FollowView) Init(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error {

	// Init map.
	v.uMap = make(map[cipher.PubKey]*FollowRep)

	// Get pages.
	pages, e := object.GetPages(pack, mux, false, false, true)
	if e != nil {
		return e
	}

	return pages.UsersPage.RangeUserActivityPages(func(i int, uap *object.UserActivityPage) error {
		return uap.RangeVoteActions(func(j int, vote *object.Vote) error {

			// Only parse if vote is of user.
			if vote.GetType() == object.UserVote {

				// Ensure creator's follow page exists.
				followRef, has := v.uMap[vote.Creator]
				if !has {
					followRef = NewFollowRep(vote.Creator)
					v.uMap[vote.Creator] = followRef
				}

				// Add stuff.
				followRef.Set(vote.OfUser, vote.Mode, vote.Tag)
			}

			return nil
		}, nil)
	}, mux)
}

func (v *FollowView) Update(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error {

	for _, vote := range headers.GetChanges().NewVotes {

		// Only parse if vote is of user.
		if vote.GetType() == object.UserVote {

			// Ensure creator's follow page exists.
			followRef, has := v.uMap[vote.Creator]
			if !has {
				followRef = NewFollowRep(vote.Creator)
				v.uMap[vote.Creator] = followRef
			}

			// Add stuff.
			followRef.Set(vote.OfUser, vote.Mode, vote.Tag)

		}
	}

	return nil
}

const (
	FollowPage = "FollowPage"
)

func (v *FollowView) Get(id string, a ...interface{}) (interface{}, error) {
	upk := a[0].(cipher.PubKey)
	switch {
	case id == FollowPage && len(a) == 1:
		fr, has := v.uMap[upk]
		if !has {
			return &FollowRepView{UserPubKey: upk.Hex()}, nil
		}
		return fr.View(), nil

	default:
		return nil, boo.Newf(boo.NotAllowed,
			"invalid get request 's' (%v)", id, a)
	}
}
