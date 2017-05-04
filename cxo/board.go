package cxo

import (
	"fmt"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

// Board represents a board as stored in cxo.
type Board struct {
	Name         string
	Threads      skyobject.References `skyobject:"schema=Thread"`
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

// GetBoardFromSkyValue obtains a Board from skyobject.Value.
func GetBoardFromSkyValue(v *skyobject.Value) (*Board, error) {
	board := Board{}

	e := v.RangeFields(func(fn string, mv *skyobject.Value) error {
		var e error
		switch fn {
		case "Name":
			board.Name, e = mv.String()
		case "Threads":
			fmt.Println("Unprocessed field: Threads")
		case "URL":
			board.URL, e = mv.String()
		case "Created":
			board.Created, e = mv.Int()
		case "LastModified":
			board.LastModified, e = mv.Int()
		case "Version":
			board.Version, e = mv.Uint()
		}
		return e
	})
	return &board, e
}

// BoardConfig represents a board's configuration as stored on a local file.
type BoardConfig struct {
	Master    bool   `json:"master"`
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key,omitempty"` // Empty if Master = false.
	URL       string `json:"url"`
}

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
		PublicKey: bc.PublicKey,
		SecretKey: bc.SecretKey,
		URL:       bc.URL,
	}
	e := client.Execute(func(ct *node.Container) error {
		root := ct.Root(pk)

		// Find latest version of board.
		boardValue, e := FindLatestBoardValueFromRoot(root)
		if e != nil {
			return e
		}
		board, e := GetBoardFromSkyValue(boardValue)
		if e != nil {
			return e
		}
		bp.Name = board.Name
		bp.Created = board.Created
		bp.LastModified = board.LastModified
		bp.Version = board.Version

		// Find thread values from latest board.
		threadValues, e := GetThreadValuesFromBoardValue(boardValue)
		if e != nil {
			return e
		}
		// Get threads from threadValues, and append to BoardPage.
		for _, v := range threadValues {
			thread, e := GetThreadFromSkyValue(v)
			if e != nil {
				return e
			}
			bp.Threads = append(bp.Threads, thread)
		}
		return nil
	})
	return bp, e
}

// FindLatestBoardValueFromRoot finds the latest board from root.
// Note that the latest board is represented as *skyobject.Value.
func FindLatestBoardValueFromRoot(r *node.Root) (*skyobject.Value, error) {
	var (
		boardValue   *skyobject.Value
		boardLastMod int64
	)
	values, e := r.Values()
	if e != nil {
		return nil, e
	}
	for _, v := range values {
		if v.Schema().Name() != "Board" {
			continue
		}
		vLastMod, e := v.FieldByName("LastModified")
		if e != nil {
			return nil, e
		}
		lastMod, _ := vLastMod.Int()
		if lastMod > boardLastMod {
			boardValue = v
			boardLastMod = lastMod
		}
	}
	return boardValue, nil
}

// GetThreadValuesFromBoardValue gets thread values from a board value.
func GetThreadValuesFromBoardValue(vBoard *skyobject.Value) ([]*skyobject.Value, error) {
	vThreads, e := vBoard.FieldByName("Threads")
	if e != nil {
		return nil, e
	}
	// Get number of threads.
	l, e := vThreads.Len()
	if e != nil {
		return nil, e
	}
	var vThreadList = make([]*skyobject.Value, l)
	for i := 0; i < l; i++ {
		vThreadList[i], e = vThreads.Index(i)
		if e != nil {
			return nil, e
		}
	}
	return vThreadList, nil
}
