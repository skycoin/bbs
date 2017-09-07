package r0

type ContentType string

const (
	ThreadType    = "thread"
	ImagePostType = "image_post"
	TextPostType  = "text_post"
)

type ContentData struct {
	Type         ContentType       `json:"content_type"`
	Name         string            `json:"heading"`
	Body         string            `json:"body"`
	SubAddresses []string          `json:"submission_addresses,omitempty"`
	Image        *ContentImageData `json:"image,omitempty"`
}

type ContentImageData struct {
	URL      string `json:"url"`
	ThumbURL string `json:"thumbnail_url,omitempty"`
	Size     int    `json:"size,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}
