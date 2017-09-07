package r0

type ContentData struct {
	Name         string              `json:"heading"`
	Body         string              `json:"body"`
	SubAddresses []string            `json:"submission_addresses,omitempty"`
	Images       []*ContentImageData `json:"images,omitempty"`
}

type ContentImageData struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	ThumbURL string `json:"thumbnail_url,omitempty"`
	Size     int    `json:"size,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}
