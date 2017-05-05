package types

import "github.com/skycoin/cxo/node"

// BoardView represents a board when displayed to user via GUI.
type BoardView struct {
	Name      string `json:"name"`
	Master    bool   `json:"master"`
	PublicKey string `json:"public_key"`

	Created      int64  `json:"created"`
	LastModified int64  `json:"last_modified"`
	Version      uint64 `json:"version"`

	ThreadCount uint64    `json:"thread_count"`
	Threads     []*Thread `json:"threads,omitempty"`
}

// NewBoardView obtains a BoardView from BoardConfig and cxo client.
// parameter "withThreads" determines whether or not to list threads in BoardView.
func NewBoardView(bc *BoardConfig, client *node.Client) (*BoardView, error) {
	// Extract data from BoardConfig.
	bv := BoardView{
		Name:      bc.Name,
		Master:    bc.Master,
		PublicKey: bc.PublicKeyStr,
	}
	// Extract data from cxo.
	board, e := ObtainLatestBoard(bc.PublicKey, client)
	if e != nil {
		return nil, e
	}
	bv.Name = board.Name
	bv.Created = board.Created
	bv.LastModified = board.LastModified
	bv.Version = board.Version

	// TODO: ObtainLatestBoardThreads. Obtain threads.
	bv.ThreadCount = 0

	return &bv, nil
}
