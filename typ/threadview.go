package typ

import (
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
)

// ThreadView represents a thread when displayed to user via GUI.
type ThreadView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Created      int64  `json:"created"`
	LastModified int64  `json:"last_modified"`
	Version      uint64 `json:"version"`

	PostCount uint64  `json:"post_count"`
	Posts     []*Post `json:"posts,omitempty"`
}

// NewThreadView obtains a ThreadView from Thread and cxo client.
func NewThreadView(bpk cipher.PubKey, t *Thread, client *node.Client, showPosts bool) (*ThreadView, error) {
	// Extract data from Thread.
	tv := ThreadView{
		ID: t.ID.Hex(),
		Name: t.Name,
		Description: t.Description,
		Created: t.Created,
		LastModified: t.LastModified,
		Version: t.Version,
	}
	// Extract data from cxo.
	tps, tpsv, e := ObtainLatestThreadPosts(bpk, t.ID, client)
	if e != nil {
		return nil, e
	}
	tv.PostCount = tps.Count
	if showPosts == true {
		posts, e := ObtainPostsFromPostsValue(tpsv)
		if e != nil {
			return nil, e
		}
		tv.Posts = posts
	}
	return &tv, nil
}