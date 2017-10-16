package keys

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"
)

const (
	// SeedBitSize represents bit size to use for seed.
	SeedBitSize = 128
	// SeedFailedMsg is the message shown on seed generation failure.
	SeedFailedMsg = "failed to generate seed."
)

type GenerateSeedOut struct {
	Seed string `json:"seed"`
}

func GenerateSeed() (*GenerateSeedOut, error) {
	entropy, e := bip39.NewEntropy(SeedBitSize)
	if e != nil {
		return nil, boo.WrapType(e, boo.Internal, SeedFailedMsg)
	}
	mnemonic, e := bip39.NewMnemonic(entropy)
	if e != nil {
		return nil, boo.WrapType(e, boo.Internal, SeedFailedMsg)
	}
	return &GenerateSeedOut{
		Seed: mnemonic,
	}, nil
}

type GenerateKeyPairIn struct {
	Seed string `json:"seed"`
}

type GenerateKeyPairOut struct {
	Input  *GenerateKeyPairIn `json:"input"`
	PubKey string             `json:"public_key"`
	SecKey string             `json:"secret_key"`
}

func GenerateKeyPair(in *GenerateKeyPairIn) (*GenerateKeyPairOut, error) {
	var (
		pk cipher.PubKey
		sk cipher.SecKey
	)
	switch in.Seed {
	case "":
		pk, sk = cipher.GenerateKeyPair()
	default:
		pk, sk = cipher.GenerateDeterministicKeyPair([]byte(in.Seed))
	}
	return &GenerateKeyPairOut{
		Input:  in,
		PubKey: pk.Hex(),
		SecKey: sk.Hex(),
	}, nil
}

type SumSHA256In struct {
	Data string `json:"data"`
}

type SumSHA256Out struct {
	Input *SumSHA256In `json:"input"`
	Hash  string       `json:"hash"`
}

func SumSHA256(in *SumSHA256In) (*SumSHA256Out, error) {
	return &SumSHA256Out{
		Input: in,
		Hash:  cipher.SumSHA256([]byte(in.Data)).Hex(),
	}, nil
}

type SignHashIn struct {
	Hash   string `json:"hash"`
	SecKey string `json:"secret_key"`
}

type SignHashOut struct {
	Input *SignHashIn `json:"input"`
	Sig   string      `json:"sig"`
}

func SignHash(in *SignHashIn) (*SignHashOut, error) {
	hash, e := GetHash(in.Hash)
	if e != nil {
		return nil, e
	}
	sk, e := GetSecKey(in.SecKey)
	if e != nil {
		return nil, e
	}
	return &SignHashOut{
		Input: in,
		Sig:   cipher.SignHash(hash, sk).Hex(),
	}, nil
}
