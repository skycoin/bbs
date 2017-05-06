package typ

// Thread represents a thread stored in cxo.
type Thread struct {
	Name    string `json:"name"`
	Desc    string `json:"description"`
	Created int64  `json:"created"`
}
