package typ

// Board represents a board stored in cxo.
type Board struct {
	Name    string `json:"name"`
	Desc    string `json:"description"`
	URL     string `json:"url"`
	Created int64  `json:"created"`
}
