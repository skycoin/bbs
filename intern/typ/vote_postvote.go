package typ

import (
	"github.com/pkg/errors"
	"github.com/skycoin/cxo/skyobject"
)

// PostVotes is an element of PostVotesContainer.
type PostVotes struct {
	Post  skyobject.Reference  `skyobject:"schema=Post"`
	Votes skyobject.References `skyobject:"schema=Vote"`
}

// PostVotesContainer contains the votes of posts.
type PostVotesContainer struct {
	Posts []PostVotes
}

// GetPostVotes obtains the thread vote references for specified thread.
func (pvc *PostVotesContainer) GetPostVotes(pRef skyobject.Reference) (*PostVotes, error) {
	for i := range pvc.Posts {
		if pvc.Posts[i].Post == pRef {
			return &pvc.Posts[i], nil
		}
	}
	return nil, errors.New("post votes not found")
}

// AddPost adds a post to PostVotesContainer.
func (pvc *PostVotesContainer) AddPost(pRef skyobject.Reference) {
	for _, p := range pvc.Posts {
		if p.Post == pRef {
			return
		}
	}
	pvc.Posts = append(pvc.Posts, PostVotes{Post: pRef})
}

// RemovePost removes a post from PostVotesContainer.
func (pvc *PostVotesContainer) RemovePost(pRef skyobject.Reference) {
	for i, p := range pvc.Posts {
		if p.Post != pRef {
			continue
		}
		// Swap i'th and last element.
		pvc.Posts[i], pvc.Posts[len(pvc.Posts)-1] =
			pvc.Posts[len(pvc.Posts)-1], pvc.Posts[i]
		// Remove last element.
		pvc.Posts = pvc.Posts[:len(pvc.Posts)-1]
		return
	}
}
