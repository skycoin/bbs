package typ

import (
	"github.com/pkg/errors"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Vote represents a post by a user.
type Vote struct {
	User cipher.PubKey // User who voted.
	Mode int8          // +1 is up, -1 is down.
	Tag  []byte        // What's this?
	Sig  cipher.Sig    // Signature.
}

func (v *Vote) checkContent() error {
	switch v.Mode {
	case 0, +1, -1:
		return nil
	default:
		return errors.New("invalid vote mode")
	}
}

func (v *Vote) checkAuthor() error {
	if v.User == (cipher.PubKey{}) {
		return errors.New("empty user")
	}
	return nil
}

// Up sets vote to up-vote.
func (v *Vote) Up() { v.Mode = 1 }

// Down sets vote to down-vote.
func (v *Vote) Down() { v.Mode = -1 }

// Sign checks and signs the post.
func (v *Vote) Sign(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := v.checkContent(); e != nil {
		return e
	}
	v.User = upk
	v.Sig = cipher.Sig{}
	v.Sig = cipher.SignHash(cipher.SumSHA256(encoder.Serialize(*v)), usk)
	return nil
}

// Verify checks the legitimacy of the post.
func (v Vote) Verify() error {
	// Check contents.
	if e := v.checkContent(); e != nil {
		return e
	}
	// Check user.
	if e := v.checkAuthor(); e != nil {
		return e
	}
	// Check signature.
	sig := v.Sig
	v.Sig = cipher.Sig{}

	return cipher.VerifySignature(
		v.User, sig,
		cipher.SumSHA256(encoder.Serialize(v)),
	)
}
