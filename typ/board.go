package typ

import (
	"errors"
	"strings"
	"time"
)

// Board represents a board stored in cxo.
type Board struct {
	Name    string `json:"name"`
	Desc    string `json:"description"`
	PubKey  string `json:"public_key" enc:"-"`
	URL     string `json:"-"`
	Created int64  `json:"created"`
}

// NewBoard creates a new Board.
func NewBoard(name, desc, url string) *Board {
	now := time.Now().UnixNano()
	return &Board{
		Name:    name,
		Desc:    desc,
		URL:     url,
		Created: now,
	}
}

// NewBoardFromConfig creates a new Board from BoardConfig.
func NewBoardFromConfig(bc *BoardConfig) *Board {
	return &Board{
		Name:    bc.Name,
		Desc:    bc.Desc,
		PubKey:  bc.PubKeyStr,
		URL:     bc.URL,
		Created: time.Now().UnixNano(),
	}
}

func (b *Board) CheckAndPrep() error {
	if b == nil {
		return errors.New("nil board")
	}
	b.Name = strings.TrimSpace(b.Name)
	b.Desc = strings.TrimSpace(b.Desc)
	b.Created = time.Now().UnixNano()
	if len(b.Name) == 0 {
		return errors.New("invalid board name")
	}
	return nil
}
