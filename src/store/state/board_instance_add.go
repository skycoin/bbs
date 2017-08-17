package state

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

func (bi *BoardInstance) NewContent(thread *object.Content) error {
	if e := thread.Verify(); e != nil {
		return e
	}
	// TODO: Check user permissions.

	bi.packMux.Lock()
	defer bi.packMux.Unlock()
	return nil
}

func (bi *BoardInstance) NewVote(post *object.Content) error {
	return nil
}

func (bi *BoardInstance) DeleteContent(cHash cipher.SHA256) error {
	return nil
}

