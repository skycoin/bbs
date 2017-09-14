package r0

import "github.com/skycoin/skycoin/src/cipher"

type ContentData struct {
	Name    string              `json:"name,omitempty"`
	Body    string              `json:"body,omitempty"`
	SubKeys []cipher.PubKey     `json:"submission_keys,omitempty"`
	Images  []*ContentImageData `json:"images,omitempty"`
}

type ContentImageData struct {
	Name   string              `json:"name"`
	Hash   string              `json:"hash"`
	URL    string              `json:"url,omitempty"`
	Size   int                 `json:"size,omitempty"`
	Height int                 `json:"height,omitempty"`
	Width  int                 `json:"width,omitempty"`
	Thumbs []*ContentImageData `json:"thumbs,omitempty"`
}