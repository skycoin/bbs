package types

import (
	"fmt"
	"github.com/skycoin/cxo/skyobject"
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
