package typ

import (
	"errors"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"time"
)

// BoardThreads references the threads of the board.
type BoardThreads struct {
	Threads      skyobject.References `skyobject:"schema=Thread"`
	Count        uint64
	LastModified int64
	Version      uint64
}

// NewBoardThreads creates a new board thread with no threads.
func NewBoardThreads() *BoardThreads {
	now := time.Now().UnixNano()
	return &BoardThreads{
		Count:        0,
		LastModified: now,
		Version:      0,
	}
}

// ObtainLatestBoardThreads obtains the latest BoardThreads of given public key from cxo.
func ObtainLatestBoardThreads(bpk cipher.PubKey, client *node.Client) (*BoardThreads, *skyobject.Value, error) {
	var bts BoardThreads
	var val *skyobject.Value
	e := client.Execute(func(ct *node.Container) error {
		// Get values from root.
		values, e := ct.Root(bpk).Values()
		if e != nil {
			return e
		}
		// Loop through values, and if type is BoardThreads, compare.
		for _, v := range values {
			if v.Schema().Name() != "BoardThreads" {
				continue
			}
			// Temporary BoardThreads.
			temp := BoardThreads{}
			if e := encoder.DeserializeRaw(v.Data(), &temp); e != nil {
				return e
			}
			if temp.Version >= bts.Version {
				bts = temp
				val = v
			}
		}
		return nil
	})
	return &bts, val, e
}

// Iterate increases the version of BoardThreads.
func (bt *BoardThreads) Iterate() {
	bt.Version += 1
	bt.LastModified = time.Now().UnixNano()
}

// AddThread adds a thread, and hence, iterates.
func (bt *BoardThreads) AddThread(bpk cipher.PubKey, client *node.Client, thread *Thread) error {
	// Save thread to cxo and get reference.
	tRef := skyobject.Reference{}
	client.Execute(func(ct *node.Container) error {
		tRef = ct.Save(thread)
		return nil
	})
	// If reference is in BoardThreads, just return as no modification needed.
	for _, r := range bt.Threads {
		if r == tRef {
			return errors.New("thread already exists")
		}
	}
	// Add ref to BoardThreads and iterate.
	bt.Threads = append(bt.Threads, tRef)
	bt.Count = uint64(len(bt.Threads))
	bt.Iterate()
	return nil
}

// ObtainThreadFromBoardThreadsValue obtains thread from value of specified id.
func ObtainThreadFromBoardThreadsValue(btsv *skyobject.Value, tid cipher.PubKey) (*Thread, error) {
	tsv, e := btsv.FieldByName("Threads")
	if e != nil {
		return nil, e
	}
	// Get number of threads.
	l, e := tsv.Len()
	if e != nil {
		return nil, e
	}
	// Loop through and find thread.
	for i := 0; i < l; i++ {
		tv, e := tsv.Index(i)
		if e != nil {
			return nil, e
		}
		v, e := tv.Dereference()
		if e != nil {
			return nil, e
		}
		if v.Schema().Name() != "Thread" {
			return nil, errors.New("value is not thread")
		}
		thread := Thread{}
		if e := encoder.DeserializeRaw(v.Data(), &thread); e != nil {
			return nil, e
		}
		if thread.ID == tid {
			return &thread, nil
		}
	}
	return nil, errors.New("thread not found")
}

// ObtainThreadsFromBoardThreadsValue obtains list of threads from value.
// TODO: Optimise.
func ObtainThreadsFromBoardThreadsValue(btsv *skyobject.Value) ([]*Thread, error) {
	tsv, e := btsv.FieldByName("Threads")
	if e != nil {
		return nil, e
	}
	// Get number of threads.
	l, e := tsv.Len()
	if e != nil {
		return nil, e
	}
	// Loop through extracting threads.
	threads := []*Thread{}
	for i := 0; i < l; i++ {
		tv, e := tsv.Index(i)
		if e != nil {
			return nil, e
		}
		v, e := tv.Dereference()
		if e != nil {
			return nil, e
		}
		if v.Schema().Name() != "Thread" {
			return nil, errors.New("value is not thread")
		}
		thread := Thread{}
		if e := encoder.DeserializeRaw(v.Data(), &thread); e != nil {
			return nil, e
		}
		threads = append(threads, &thread)
	}
	return threads, nil
}