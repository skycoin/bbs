package datastore

import (
	"github.com/skycoin/cxo/skyobject"
	"time"
)

// Thread represents a thread.
type Thread struct {
	ID           []byte `json:"id"`
	Name        string `json:"name"`
	Description  string `json:"description"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"last_modified"`
	Version      uint64 `json:"version"`
}

// NewThread creates a new thread.
func NewThread(name, desc string) *Thread {
	now := time.Now().UnixNano()
	thread := Thread{
		ID:          MakeTimeStampedRandomID(128),
		Name:       name,
		Description: desc,
		Created:     now,
		LastModified: now,
		Version: 0,
	}
	return &thread
}

// GetThreadFromSkyValue obtains a Thread from *skyobject.Value.
func GetThreadFromSkyValue(v *skyobject.Value) (*Thread, error) {
	thread := Thread{}
	e := v.RangeFields(func(fn string, mv *skyobject.Value) error {
		var e error
		switch fn {
		case "ID":
			thread.ID, e = mv.Bytes()
		case "Name":
			thread.Name, e = mv.String()
		case "Description":
			thread.Description, e = mv.String()
		case "Created":
			thread.Created, e = mv.Int()
		case "LastModified":
			thread.LastModified, e = mv.Int()
		case "Version":
			thread.Version, e = mv.Uint()
		}
		return e
	})
	return &thread, e
}

// ThreadPage represents a page of posts in a thread.
type ThreadPage struct {
	Thread *Thread `json:"thread"`
	Posts  []*Post `json:"posts"`
}

// NewThreadPageFromThread creates a ThreadPage from Thread.
func NewThreadPageFromThread(t *Thread) *ThreadPage {
	return &ThreadPage{
		Thread: t,
	}
}
