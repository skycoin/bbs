package typ

import (
	"errors"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"time"
)

// ThreadPosts references the posts of a specified thread.
type ThreadPosts struct {
	ThreadID     cipher.PubKey
	Posts        skyobject.References `skyobject:"schema=Post"`
	Count        uint64
	LastModified int64
	Version      uint64
}

// NewThreadPosts creates a new ThreadPosts with specified Thread ID.
func NewThreadPosts(tid cipher.PubKey) *ThreadPosts {
	now := time.Now().UnixNano()
	return &ThreadPosts{
		ThreadID:     tid,
		Count:        0,
		LastModified: now,
		Version:      0,
	}
}

// ObtainLatestThreadPosts obtains the latest ThreadPosts of given board public key and thread id.
func ObtainLatestThreadPosts(bpk cipher.PubKey, tid cipher.PubKey, client *node.Client) (*ThreadPosts, *skyobject.Value, error) {
	var tps ThreadPosts
	var val *skyobject.Value
	e := client.Execute(func(ct *node.Container) error {
		// Get values from root.
		values, e := ct.Root(bpk).Values()
		if e != nil {
			return e
		}
		// Loop through values, and if type is ThreadPosts, compare if thread id is tid.
		for _, v := range values {
			if v.Schema().Name() != "ThreadPosts" {
				continue
			}
			// Temporary ThreadPosts.
			temp := ThreadPosts{}
			if e := encoder.DeserializeRaw(v.Data(), &temp); e != nil {
				return e
			}
			// Check thread id.
			if temp.ThreadID != tid {
				continue
			}
			// Compare.
			if temp.Version >= tps.Version {
				tps = temp
				val = v
			}
		}
		return nil
	})
	if e != nil {
		return nil, nil, e
	}
	if len(tps.ThreadID) == 0 {
		return nil, nil, errors.New("thread does not exist")
	}
	return &tps, val, e
}

// ObtainPostsFromPostsValue obtains list of posts from value.
// TODO: Optimise.
func ObtainPostsFromPostsValue(tpsv *skyobject.Value) ([]*Post, error) {
	psv, e := tpsv.FieldByName("Posts")
	if e != nil {
		return nil, e
	}
	// Get number of posts.
	l, e := psv.Len()
	if e != nil {
		return nil, e
	}
	// Loop through extracting posts.
	posts := []*Post{}
	for i := 0; i < l; i++ {
		pv, e := psv.Index(i)
		if e != nil {
			return nil, e
		}
		v, e := pv.Dereference()
		if e != nil {
			return nil, e
		}
		if v.Schema().Name() != "Post" {
			return nil, errors.New("value is not post")
		}
		post := Post{}
		if e := encoder.DeserializeRaw(v.Data(), &post); e != nil {
			return nil, e
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

// Iterate increases the version of ThreadPosts.
func (tp *ThreadPosts) Iterate() {
	tp.Version += 1
	tp.LastModified = time.Now().UnixNano()
}

// AddPost adds a post, and hence, iterates.
func (tp *ThreadPosts) AddPost(bpk cipher.PubKey, client *node.Client, post *Post) error {
	// Save thread to cxo and get reference.
	pRef := skyobject.Reference{}
	client.Execute(func(ct *node.Container) error {
		pRef = ct.Save(post)
		return nil
	})
	// If reference is in ThreadPosts, just return as no modification needed.
	for _, r := range tp.Posts {
		if r == pRef {
			return errors.New("post already exists")
		}
	}
	// Add ref to ThreadPosts and iterate.
	tp.Posts = append(tp.Posts, pRef)
	tp.Count = uint64(len(tp.Posts))
	tp.Iterate()
	return nil
}
