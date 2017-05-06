package typ

import (
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"time"
)

// Board represents a board as stored in cxo.
type Board struct {
	Name         string
	URL          string
	Created      int64
	LastModified int64
	Version      uint64
}

// NewBoard creates a new board with given name and url.
func NewBoard(name, url string) *Board {
	now := time.Now().UnixNano()
	return &Board{
		Name:         name,
		URL:          url,
		Created:      now,
		LastModified: now,
		Version:      0,
	}
}

// ObtainLatestBoard obtains the latest board of given public key from cxo.
func ObtainLatestBoard(bpk cipher.PubKey, client *node.Client) (*Board, error) {
	var board Board
	e := client.Execute(func(ct *node.Container) error {
		// Get values from root.
		values, e := ct.Root(bpk).Values()
		if e != nil {
			return e
		}
		// Loop through values, and if type is Board, compare.
		for _, v := range values {
			if v.Schema().Name() != "Board" {
				continue
			}
			// Temporary board.
			temp := Board{}
			if e := encoder.DeserializeRaw(v.Data(), &temp); e != nil {
				return e
			}
			if temp.Version >= board.Version {
				board = temp
			}
		}
		return nil
	})
	return &board, e
}
