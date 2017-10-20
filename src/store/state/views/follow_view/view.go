package follow_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
)

type FollowView struct {
	uMap map[string]*FollowRep // key(user public key), value(rep)
}

func (v *FollowView) Init(pack *skyobject.Pack, headers *pack.Headers) error {

	// Init map.
	v.uMap = make(map[string]*FollowRep)

	// Get pages.
	pages, e := object.GetPages(pack, &object.GetPagesIn{
		RootPage:  false,
		BoardPage: false,
		DiffPage:  false,
		UsersPage: true,
	})
	if e != nil {
		return e
	}

	return pages.UsersPage.RangeUserProfiles(func(i int, profile *object.UserProfile) error {
		return profile.RangeSubmissions(func(i int, c *object.Content) error {
			cBody := c.GetBody()

			// Only parse if content is user vote.
			if cBody.Type == object.V5UserVoteType {

				// Ensure creator's profile exists.
				followRep, ok := v.uMap[cBody.Creator]
				if !ok {
					followRep = NewFollowRep(cBody.Creator)
					v.uMap[cBody.Creator] = followRep
				}

				// Add stuff.
				followRep.Set(cBody.Creator, cBody.Value, cBody.Tag)
			}

			return nil
		})
	})
}

func (v *FollowView) Update(pack *skyobject.Pack, headers *pack.Headers) error {

	for _, c := range headers.GetChanges().New {
		cBody := c.GetBody()

		// Only parse if content is user type.
		if cBody.Type == object.V5UserVoteType {
			// Ensure creator's profile exists.
			followRep, ok := v.uMap[cBody.Creator]
			if !ok {
				followRep = NewFollowRep(cBody.Creator)
				v.uMap[cBody.Creator] = followRep
			}

			// Add stuff.
			followRep.Set(cBody.Creator, cBody.Value, cBody.Tag)
		}
	}

	return nil
}

const (
	FollowPage = "FollowPage"
)

func (v *FollowView) Get(id string, a ...interface{}) (interface{}, error) {
	upk := a[0].(string)
	switch {
	case id == FollowPage && len(a) == 1:
		fr, has := v.uMap[upk]
		if !has {
			return NewFollowRep(upk), nil
		}
		return fr.View(), nil

	default:
		return nil, boo.Newf(boo.NotAllowed,
			"invalid get request 's' (%v)", id, a)
	}
}
