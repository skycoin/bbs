package keys

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"
)

const (
	// SeedBitSize represents bit size to use for seed.
	SeedBitSize = 128
	// SeedFailedMsg is the message shown on seed generation failure.
	SeedFailedMsg = "failed to generate seed."
)

// GenerateSeed generates a seed.
func GenerateSeed() (string, error) {
	entropy, e := bip39.NewEntropy(SeedBitSize)
	if e != nil {
		return "", boo.WrapType(e, boo.Internal, SeedFailedMsg)
	}
	mnemonic, e := bip39.NewMnemonic(entropy)
	if e != nil {
		return "", boo.WrapType(e, boo.Internal, SeedFailedMsg)
	}
	return mnemonic, nil
}
