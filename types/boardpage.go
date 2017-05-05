package types

import (
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
)

// BoardPage represents a board page as displayed to user.
// Lists all threads and board meta.
type BoardPage struct {
	Name         string    `json:"name"`
	Master       bool      `json:"master"`
	PublicKey    string    `json:"public_key"`
	SecretKey    string    `json:"secret_key,omitempty"`
	URL          string    `json:"url"`
	Created      int64     `json:"created"`
	LastModified int64     `json:"last_modified"`
	Version      uint64    `json:"version"`
	Threads      []*Thread `json:"threads"`
}

// NewBoardPage creates a new BoardPage from cipher.PubKey and *node.Client.
func (bc *BoardConfig) NewBoardPage(pk cipher.PubKey, client *node.Client) (*BoardPage, error) {
	bp := &BoardPage{
		Master:    bc.Master,
		PublicKey: bc.PublicKeyStr,
		SecretKey: bc.SecretKeyStr,
		URL:       bc.URL,
	}
	// Get latest board from cxo and add board data to board page.
	board, e := ObtainLatestBoard(pk, client)
	if e != nil {
		return nil, e
	}
	bp.Name = board.Name
	bp.Created = board.Created
	bp.LastModified = board.LastModified
	bp.Version = board.Version

	// TODO: ObtainLatestBoardThreads. Obtain threads.

	return bp, nil
}
