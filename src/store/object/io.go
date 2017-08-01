package object

import (
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// RetryIO represents io required on connection/subscription retries.
type RetryIO struct {
	PublicKeys []cipher.PubKey
	Addresses  []string
}

func (io *RetryIO) Process() error {
	return nil
}

func (io *RetryIO) IsEmpty() bool {
	return len(io.PublicKeys) == 0 && len(io.Addresses) == 0
}

func (io *RetryIO) HasPK(pk cipher.PubKey) int {
	out := -1
	for i, gotPK := range io.PublicKeys {
		if gotPK == pk {
			return i
		}
	}
	return out
}

func (io *RetryIO) HasAddress(address string) int {
	out := -1
	for i, gotAddress := range io.Addresses {
		if gotAddress == address {
			return i
		}
	}
	return out
}

func (io *RetryIO) Add(in *RetryIO) {
	if in == nil {
		return
	}
	// Add unique public keys.
	pksToAdd := []cipher.PubKey{}
	for _, newPK := range in.PublicKeys {
		got := false
		for _, oldPK := range io.PublicKeys {
			if newPK == oldPK {
				got = true
				break
			}
		}
		if !got {
			pksToAdd = append(pksToAdd, newPK)
		}
	}
	io.PublicKeys = append(io.PublicKeys, pksToAdd...)

	// Add unique addresses.
	addressesToAdd := []string{}
	for _, newAddress := range in.Addresses {
		got := false
		for _, oldAddress := range io.Addresses {
			if newAddress == oldAddress {
				got = true
				break
			}
		}
		if !got {
			addressesToAdd = append(addressesToAdd, newAddress)
		}
	}
	io.Addresses = append(io.Addresses, addressesToAdd...)
}

func (io *RetryIO) Fill(subscriptions []Subscription, connections []string) *RetryIO {
	for _, sub := range subscriptions {
		if io.HasPK(sub.PubKey) == -1 {
			io.PublicKeys = append(io.PublicKeys, sub.PubKey)
		}
	}
	for _, address := range connections {
		if io.HasAddress(address) == -1 {
			io.Addresses = append(io.Addresses, address)
		}
	}
	return io
}

// NewUserIO represents io required when creating a new user.
type NewUserIO struct {
	Alias      string        `bbs:"alias"`
	Seed       string        `bbs:"uSeed"`
	Password   string        `bbs:"password"`
	UserPubKey cipher.PubKey `bbs:"upk"`
	UserSecKey cipher.SecKey `bbs:"usk"`
}

// LoginIO represents input required when logging io.
type LoginIO struct {
	Alias    string `bbs:"alias"`
	Password string `bbs:"password"`
}

// ConnectionIO represents input/output required when connection/disconnecting from address.
type ConnectionIO struct {
	Address string `bbs:"address"`
}

// BoardIO represents a subscription input.
type BoardIO struct {
	PubKeyStr string        `bbs:"bpkStr"`
	PubKey    cipher.PubKey `bbs:"bpk"`
	SecKey    cipher.SecKey `bbs:"bsk"`
}

// AddressIO represents input/output regarding addresses.
type AddressIO struct {
	PubKeyStr string        `bbs:"bpkStr"`
	PubKey    cipher.PubKey `bbs:"bpk"`
	SecKey    cipher.SecKey `bbs:"bsk"`
	Address   string        `bbs:"address"`
}

// NewBoardIO represents input required to create master board.
type NewBoardIO struct {
	Seed                   string        `bbs:"bSeed"`
	Name                   string        `bbs:"name"`
	Desc                   string        `bbs:"description"`
	SubmissionAddressesStr string        `bbs:"subAddrsStr"` // Separated with commas.
	SubmissionAddresses    []string      `bbs:"subAddrs"`
	ConnectionsStr         string        `bbs:"consStr"` // Separated with commas.
	Connections            []string      `bbs:"cons"`
	BoardPubKey            cipher.PubKey `bbs:"bpk"`
	BoardSecKey            cipher.SecKey `bbs:"bsk"`
}

type ThreadIO struct {
	BoardPubKeyStr string              `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey       `bbs:"bpk"`
	BoardSecKey    cipher.SecKey       `bbs:"bsk"`
	ThreadRefStr   string              `bbs:"tRefStr"`
	ThreadRef      skyobject.Reference `bbs:"tRef"`
}

type NewThreadIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	BoardSecKey    cipher.SecKey `bbs:"bsk"`
	UserPubKey     cipher.PubKey `bbs:"upk"`
	UserSecKey     cipher.SecKey `bbs:"usk"`
	Title          string        `bbs:"heading"`
	Body           string        `bbs:"body"`
}

type VoteThreadIO struct {
	BoardPubKeyStr string              `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey       `bbs:"bpk"`
	BoardSecKey    cipher.SecKey       `bbs:"bsk"`
	ThreadRefStr   string              `bbs:"tRefStr"`
	ThreadRef      skyobject.Reference `bbs:"tRef"`
	UserPubKey     cipher.PubKey       `bbs:"upk"`
	UserSecKey     cipher.SecKey       `bbs:"usk"`
	ModeStr        string              `bbs:"modeStr"`
	Mode           int8                `bbs:"mode"`
	TagStr         string              `bbs:"tagStr"`
	Tag            []byte              `bbs:"tag"`
}

type PostIO struct {
	ThreadIO
	PostRefStr string              `bbs:"pRefStr"`
	PostRef    skyobject.Reference `bbs:"pRef"`
}

type NewPostIO struct {
	NewThreadIO
	ThreadRefStr string              `bbs:"tRefStr"`
	ThreadRef    skyobject.Reference `bbs:"tRef"`
}

type VotePostIO struct {
	BoardPubKeyStr string              `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey       `bbs:"bpk"`
	BoardSecKey    cipher.SecKey       `bbs:"bsk"`
	PostRefStr     string              `bbs:"pRefStr"`
	PostRef        skyobject.Reference `bbs:"pRef"`
	UserPubKey     cipher.PubKey       `bbs:"upk"`
	UserSecKey     cipher.SecKey       `bbs:"usk"`
	ModeStr        string              `bbs:"modeStr"`
	Mode           int8                `bbs:"mode"`
	TagStr         string              `bbs:"tagStr"`
	Tag            []byte              `bbs:"tag"`
}