package typ

import (
	"errors"
	"github.com/skycoin/cxo/skyobject"
)

// ThreadPage represents a ThreadPage as stored in cxo.
type ThreadPage struct {
	Thread skyobject.Reference  `skyobject:"schema=Thread"`
	Posts  skyobject.References `skyobject:"schema=Post"`
}

func NewThreadPage(tRef skyobject.Reference) *ThreadPage {
	return &ThreadPage{
		Thread: tRef,
	}
}

func (tp *ThreadPage) AddPost(pRef skyobject.Reference) error {
	for _, pr := range tp.Posts {
		if pr == pRef {
			return errors.New("post already exists")
		}
	}
	tp.Posts = append(tp.Posts, pRef)
	return nil
}
