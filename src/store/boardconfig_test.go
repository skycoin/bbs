package store

import (
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"strconv"
	"testing"
	//"log"
	"fmt"
)

var (
	nPK = 5
	pks []cipher.PubKey

	nRef = 30
	refs []skyobject.Reference
)

func init() {
	// Public keys for boards.
	pks = make([]cipher.PubKey, nPK)
	for i := 0; i < nPK; i++ {
		pks[i], _ = cipher.GenerateDeterministicKeyPair([]byte(strconv.Itoa(i)))
	}
	// References for threads.
	refs = make([]skyobject.Reference, nRef)
	for i := 0; i < nRef; i++ {
		refs[i] = misc.MakeDeterministicRef(byte(i))
	}
}

func TestBoardConfig_AddDep(t *testing.T) {
	t.Run("single_threads", func(t *testing.T) {
		bc := &BoardConfig{}
		for i := 0; i < 5; i++ {
			bc.AddDep(pks[i], refs[i])
			if len(bc.Deps[i].Threads) != 1 {
				t.Errorf(`len(bc.Deps[%d].Threads) != 1`, i)
			}
		}
		fmt.Println(bc.String(true))
		if len(bc.Deps) != 5 {
			t.Error(`len(bc.Deps) != 5`)
		}
	})
	t.Run("multi_threads", func(t *testing.T) {
		bc := &BoardConfig{}
		for i := 1; i < 5; i++ {
			for j := 1; j < 4; j++ {
				bc.AddDep(pks[i], refs[i*j])
			}
			if len(bc.Deps[i-1].Threads) != 3 {
				t.Errorf(`len(bc.Deps[%d].Threads) != 3`, i)
			}
		}
		fmt.Println(bc.String(true))
		if len(bc.Deps) != 4 {
			t.Error(`len(bc.Deps) != 5`)
		}
	})
	t.Run("duplicate_threads", func(t *testing.T) {
		bc := &BoardConfig{}
		for i := 0; i < 5; i++ {
			bc.AddDep(pks[0], refs[0])
		}
		fmt.Println(bc.String(true))
		if len(bc.Deps[0].Threads) != 1 {
			t.Error(`len(bc.Deps[0].Threads) != 1`)
		}
	})
}

func TestBoardConfig_RemoveDep(t *testing.T) {
	t.Run("remove_many_threads_from_single_board", func(t *testing.T) {
		bc := &BoardConfig{}
		for i := 0; i < 10; i++ {
			bc.AddDep(pks[0], refs[i])
		}
		for i := 0; i < 5; i++ {
			bc.RemoveDep(pks[0], refs[int(float32(i)+0.5)*2])
		}
		fmt.Println(bc.String(true))
		if len(bc.Deps[0].Threads) != 5 {
			t.Error(`len(bc.Deps[0].Threads) != 5`)
		}
	})
	t.Run("empty_board", func(t *testing.T) {
		bc := &BoardConfig{}
		for i := 0; i < 10; i++ {
			bc.AddDep(pks[0], refs[i])
		}
		for i := 9; i >= 0; i-- {
			bc.RemoveDep(pks[0], refs[i])
		}
		fmt.Println(bc.String(true))
		if len(bc.Deps) != 0 {
			t.Error(`len(bc.Deps) != 0`)
		}
	})
}
