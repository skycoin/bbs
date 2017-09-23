package r0

import (
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type BoardData struct {
	Name    string   `json:"name"`            // Name of board.
	Body    string   `json:"body"`            // Description of board.
	Created int64    `json:"created"`         // Time of creation (unix time in ns).
	SubKeys []string `json:"submission_keys"` // Submission public keys.
	Tags    []string `json:"tags"`            // Tags used for searching.
}

func (bd *BoardData) GetSubKeys() []cipher.PubKey {
	out := make([]cipher.PubKey, len(bd.SubKeys))
	for i, subPK := range bd.SubKeys {
		var e error
		if out[i], e = keys.GetPubKey(subPK); e != nil {
			log.Printf("error obtaining 'submission_keys'[%d]", i)
		}
	}
	return out
}

func (bd *BoardData) SetSubKeys(sPKs []cipher.PubKey) {
	bd.SubKeys = make([]string, len(sPKs))
	for i, sPK := range sPKs {
		bd.SubKeys[i] = sPK.Hex()
	}
}

type ThreadData struct {
	OfBoard string `json:"of_board"` // Public key of board of which thread belongs.
	Name    string `json:"name"`     // Name of thread.
	Body    string `json:"body"`     // Body of thread.
	Created int64  `json:"created"`  // Time of creation (unix time in ns).
	Creator string `json:"creator"`  // Public key of creator (in hex).
}

func (td *ThreadData) GetOfBoard() cipher.PubKey {
	pk, e := keys.GetPubKey(td.OfBoard)
	if e != nil {
		log.Println("failed to get 'of_board' from thread:", e)
	}
	return pk
}

func (td *ThreadData) GetCreator() cipher.PubKey {
	pk, e := keys.GetPubKey(td.Creator)
	if e != nil {
		log.Println("failed to get 'creator' from thread:", e)
	}
	return pk
}

type PostData struct {
	OfBoard  string       `json:"of_board"`          // Public key of board in which post belongs.
	OfThread string       `json:"of_thread"`         // SHA256 of thread in which post belongs.
	OfPost   string       `json:"of_post,omitempty"` // SHA256 of post this post is replying to (optional).
	Name     string       `json:"name"`              // Name of post.
	Body     string       `json:"body"`              // Body of post.
	Images   []*ImageData `json:"images,omitempty"`  // Images of post (optional).
	Created  int64        `json:"created"`           // Time of creation (unix time in ns).
	Creator  string       `json:"creator"`           // Public key of creator (in hex).
}

func (pd *PostData) GetOfBoard() cipher.PubKey {
	pk, e := keys.GetPubKey(pd.OfBoard)
	if e != nil {
		log.Println("failed to get 'of_board' from post:", e)
	}
	return pk
}

func (pd *PostData) GetOfThread() cipher.SHA256 {
	tr, e := keys.GetHash(pd.OfThread)
	if e != nil {
		log.Println("failed to get 'of_thread' from post:", e)
	}
	return tr
}

func (pd *PostData) GetOfPost() (cipher.SHA256, bool) {
	empty := cipher.SHA256{}
	if pd.OfPost == empty.Hex() || pd.OfPost == "" {
		return empty, false
	}
	pRef, e := keys.GetHash(pd.OfPost)
	if e != nil {
		log.Println("failed to get 'of_post' from post:", e)
	}
	return pRef, true
}

func (pd *PostData) GetCreator() cipher.PubKey {
	pk, e := keys.GetPubKey(pd.Creator)
	if e != nil {
		log.Println("failed to get 'creator' from post:", pk)
	}
	return pk
}

type ImageData struct {
	Name   string       `json:"name"`
	Hash   string       `json:"hash"`
	URL    string       `json:"url,omitempty"`
	Size   int          `json:"size,omitempty"`
	Height int          `json:"height,omitempty"`
	Width  int          `json:"width,omitempty"`
	Thumbs []*ImageData `json:"thumbs,omitempty"`
}
