package state

import (
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
)

type InitBoardInstance func(ct *skyobject.Container, root *skyobject.Root) (
	*BoardInstance, error)

type BoardInstanceConfig struct {
	Master bool
	PK     cipher.PubKey
	SK     cipher.SecKey
}

type BoardInstance struct {
	c    *BoardInstanceConfig
	l    *log.Logger
	flag skyobject.Flag // Used for compiling pack.

	seqMux  sync.RWMutex
	seqChan chan struct{} // Triggered when sequence increments.
	seq     uint64        // Last completed sequence.

	packMux sync.Mutex
	pack    *skyobject.Pack


}

func NewBoardInstance(config *BoardInstanceConfig) InitBoardInstance {
	return func(ct *skyobject.Container, root *skyobject.Root) (*BoardInstance, error) {

		// Prepare output.
		bi := &BoardInstance{
			c:       config,
			l:       inform.NewLogger(true, os.Stdout, "INSTANCE:"+config.PK.Hex()),
			seqChan: make(chan struct{}),
		}

		// Prepare flags.
		bi.flag = skyobject.HashTableIndex | skyobject.EntireTree
		if !bi.c.Master {
			bi.flag |= skyobject.ViewOnly
		}

		// Prepare pack.
		var e error
		bi.pack, e = ct.Unpack(root, bi.flag, ct.CoreRegistry().Types(), config.SK)
		if e != nil {
			return nil, e
		}

		// Output.
		return bi, nil
	}
}

func (bi *BoardInstance) compile(newPack *skyobject.Pack) error {
	bi.packMux.Lock()
	defer bi.packMux.Unlock()

	// TODO: Implement.

	// Success.
	bi.SetSeq(bi.pack.Root().Seq)
	return nil
}

func (bi *BoardInstance) GetSeq() uint64 {
	bi.seqMux.RLock()
	defer bi.seqMux.RUnlock()
	return bi.seq
}

func (bi *BoardInstance) SetSeq(seq uint64) {
	bi.seqMux.Lock()
	defer bi.seqMux.Unlock()
	if seq > bi.seq {
		bi.seq = seq
		for {
			select {
			case bi.seqChan<- struct{}{}:
			default:
				return
			}
		}
	}
}
