package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/skycoin/src/cipher"
	"strings"
)

var (
	ErrInvalidSeed     = boo.New(boo.InvalidInput, "invalid seed provided")
	ErrInvalidName     = boo.New(boo.InvalidInput, "invalid name provided")
	ErrInvalidDesc     = boo.New(boo.InvalidInput, "invalid description provided")
	ErrInvalidAlias    = boo.New(boo.InvalidInput, "invalid alias provided")
	ErrInvalidPassword = boo.New(boo.InvalidInput, "invalid password provided")
)

// CheckSeed ensures validity of seed. TODO
func CheckSeed(seed string) error {
	return nil
}

// CheckName ensures validity of board/thread/post name. TODO
func CheckName(name string) error {
	return nil
}

// CheckDesc ensures validity of board/thread/post description. TODO
func CheckDesc(desc string) error {
	return nil
}

// CheckAlias ensures validity of user alias. TODO
func CheckAlias(alias string) error {
	return nil
}

// CheckPassword ensures validity of password. TODO
func CheckPassword(password string) error {
	return nil
}

// CheckAddress ensures validity of address. TODO
func CheckAddress(address string) error {
	return nil
}

// RetryIO represents io required on connection/subscription retries.
type RetryIO struct {
	pks       []cipher.PubKey
	addresses []string
}

func (io *RetryIO) Process() error {
	return nil
}

func (io *RetryIO) IsEmpty() bool {
	return len(io.pks) == 0 && len(io.addresses) == 0
}

func (io *RetryIO) HasPK(pk cipher.PubKey) int {
	out := -1
	for i, gotPK := range io.pks {
		if gotPK == pk {
			return i
		}
	}
	return out
}

func (io *RetryIO) HasAddress(address string) int {
	out := -1
	for i, gotAddress := range io.addresses {
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
	for _, newPK := range in.pks {
		got := false
		for _, oldPK := range io.pks {
			if newPK == oldPK {
				got = true
				break
			}
		}
		if !got {
			pksToAdd = append(pksToAdd, newPK)
		}
	}
	io.pks = append(io.pks, pksToAdd...)

	// Add unique addresses.
	addressesToAdd := []string{}
	for _, newAddress := range in.addresses {
		got := false
		for _, oldAddress := range io.addresses {
			if newAddress == oldAddress {
				got = true
				break
			}
		}
		if !got {
			addressesToAdd = append(addressesToAdd, newAddress)
		}
	}
	io.addresses = append(io.addresses, addressesToAdd...)
}

func (io *RetryIO) Fill(file *UserFile) *RetryIO {
	if file == nil {
		return io
	}
	for _, sub := range append(file.Subscriptions, file.Masters...) {
		if io.HasPK(sub.PubKey) == -1 {
			io.pks = append(io.pks, sub.PubKey)
		}
	}
	for _, address := range file.Connections {
		if io.HasAddress(address) == -1 {
			io.addresses = append(io.addresses, address)
		}
	}
	return io
}

// NewUserIO represents io required when creating a new user.
type NewUserIO struct {
	Alias    string `json:"alias"`
	Seed     string `json:"seed"`
	Password string `json:"password"`
}

func (io *NewUserIO) Process() error {
	if e := CheckAlias(io.Alias); e != nil {
		return e
	}
	if e := CheckSeed(io.Seed); e != nil {
		return e
	}
	if e := CheckPassword(io.Password); e != nil {
		return e
	}
	return nil
}

// LoginIO represents input required when logging in.
type LoginIO struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

func (io *LoginIO) Process() error {
	if e := CheckAlias(io.Alias); e != nil {
		return e
	}
	if e := CheckPassword(io.Password); e != nil {
		return e
	}
	return nil
}

// ConnectionIO represents input/output required when connection/disconnecting from address.
type ConnectionIO struct {
	Address string `json:"address"`
}

func (io *ConnectionIO) Process() error {
	if e := CheckAddress(io.Address); e != nil {
		return e
	}
	return nil
}

// SubscriptionIO represents a subscription input.
type SubscriptionIO struct {
	PubKey string `json:"public_key"`
	pubKey cipher.PubKey
}

// Process checks and processes input for SubscriptionIO.
func (in *SubscriptionIO) Process() (e error) {
	in.pubKey, e = keys.GetPubKey(in.PubKey)
	return
}

// NewMasterIO represents input required to create master board.
type NewMasterIO struct {
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

// Process checks and processes input for NewMasterIO.
func (in *NewMasterIO) Process() error {
	if e := CheckSeed(in.Seed); e != nil {
		return e
	}
	if e := CheckName(in.Name); e != nil {
		return e
	}
	if e := CheckDesc(in.Desc); e != nil {
		return e
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
