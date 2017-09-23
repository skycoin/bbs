package object

import "sync"

type Tag struct {
	Mode string // [+1, 0, -1]
	Text string // Optional
}

type FollowPage struct {
	sync.Mutex
	UserPubKey string         `json:"user_public_key"` // User's public key.
	Yes        map[string]Tag `json:"yes"`             //
	No         map[string]Tag `json:"no"`              //
}
