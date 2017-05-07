package typ

import "time"

// Board represents a board stored in cxo.
type Board struct {
	Name    string `json:"name"`
	Desc    string `json:"description"`
	URL     string `json:"url"`
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