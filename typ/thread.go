package typ

import (
	"github.com/skycoin/skycoin/src/cipher"
	"time"
	"github.com/skycoin/cxo/node"
)

// Thread represents a thread.
type Thread struct {
	ID           cipher.PubKey `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Created      int64         `json:"created"`
	LastModified int64         `json:"last_modified"`
	Version      uint64        `json:"version"`
}

// NewThread creates a new thread.
func NewThread(name, desc string) *Thread {
	now := time.Now().UnixNano()
	thread := Thread{
		ID:           MakeTimeStampedRandomID(128),
		Name:         name,
		Description:  desc,
		Created:      now,
		LastModified: now,
		Version:      0,
	}
	return &thread
}

// ObtainLatestThread obtains the latest thread of given board and thread id from cxo.
func ObtainLatestThread(pbk, tid cipher.PubKey, client *node.Client) (*Thread, error) {
	var thread Thread
	e := client.Execute(func(ct *node.Container) error {

		return nil
	})
	return &thread, e
}