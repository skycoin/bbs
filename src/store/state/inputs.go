package state

import (
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/skycoin/src/cipher"
	"strings"
	"github.com/skycoin/bbs/src/misc/keys"
)

var (
	ErrInvalidSeed = boo.New(boo.InvalidInput, "invalid seed provided")
	ErrInvalidName = boo.New(boo.InvalidInput, "invalid name provided")
	ErrInvalidDesc = boo.New(boo.InvalidInput, "invalid description provided")
)

// NewUserInput represents input required when creating a new user.
type NewUserInput struct {
	Alias    string `json:"alias"`
	Seed     string `json:"seed"`
	Password string `json:"password"`
}

// LoginInput represents input required when logging in.
type LoginInput struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

// SubscriptionInput represents a subscription input.
type SubscriptionInput struct {
	PubKey string `json:"public_key"`
	pubKey cipher.PubKey
}

func (in *SubscriptionInput) process() (e error) {
	in.pubKey, e = keys.GetPubKey(in.PubKey)
	return
}

// NewMasterInput represents input required to create master board.
type NewMasterInput struct {
	Seed                string `json:"seed"`
	Name                string `json:"name"`
	Desc                string `json:"description"`
	SubmissionAddresses string `json:"submission_addresses"` // Separated with commas.
	Connections         string `json:"connections"`          // Separated with commas.

	submissionAddresses []string
	connections         []string
	pubKey              cipher.PubKey
	secKey              cipher.SecKey
}

func (in *NewMasterInput) process() error {
	switch {
	case len(in.Seed) < 1:
		return ErrInvalidSeed
	case len(in.Name) < 1:
		return ErrInvalidName
	case len(in.Desc) < 1:
		return ErrInvalidDesc
	}
	in.submissionAddresses =
		strings.Split(in.SubmissionAddresses, ",")
	for i := range in.submissionAddresses {
		in.submissionAddresses[i] =
			strings.TrimSpace(in.submissionAddresses[i])
	}
	in.connections =
		strings.Split(in.Connections, ",")
	for i := range in.connections {
		in.connections[i] =
			strings.TrimSpace(in.connections[i])
	}
	in.pubKey, in.secKey =
		cipher.GenerateDeterministicKeyPair([]byte(in.Seed))
	return nil
}
