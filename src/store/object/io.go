package object

import (
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"strings"
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

// LoginIO represents input required when logging io.
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

// BoardIO represents a subscription input.
type BoardIO struct {
	PubKey string `json:"public_key"`
	SecKey cipher.SecKey
	pubKey cipher.PubKey
}

// Process checks and processes input for BoardIO.
func (io *BoardIO) Process() (e error) {
	io.pubKey, e = keys.GetPubKey(io.PubKey)
	return
}

func (io *BoardIO) GetPK() cipher.PubKey {
	return io.pubKey
}

func (io *BoardIO) GetSK() cipher.SecKey {
	return io.SecKey
}

// AddressIO represents input/output regarding addresses.
type AddressIO struct {
	PubKey  string `json:"public_key"`
	SecKey  cipher.SecKey
	Address string `json:"string"`
	pubKey  cipher.PubKey
}

func (io *AddressIO) Process() error {
	var e error
	io.pubKey, e = keys.GetPubKey(io.PubKey)
	if e != nil {
		return e
	}
	if e = CheckAddress(io.Address); e != nil {
		return e
	}
	return nil
}

func (io *AddressIO) GetPK() cipher.PubKey {
	return io.pubKey
}

// NewBoardIO represents input required to create master board.
type NewBoardIO struct {
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

// Process checks and processes input for NewBoardIO.
func (io *NewBoardIO) Process() error {
	if e := CheckSeed(io.Seed); e != nil {
		return e
	}
	if e := CheckName(io.Name); e != nil {
		return e
	}
	if e := CheckDesc(io.Desc); e != nil {
		return e
	}
	// TODO: Fix empty addresses.
	io.submissionAddresses =
		strings.Split(io.SubmissionAddresses, ",")
	for i := range io.submissionAddresses {
		io.submissionAddresses[i] =
			strings.TrimSpace(io.submissionAddresses[i])
	}
	io.connections =
		strings.Split(io.Connections, ",")
	for i := range io.connections {
		io.connections[i] =
			strings.TrimSpace(io.connections[i])
	}
	io.pubKey, io.secKey =
		cipher.GenerateDeterministicKeyPair([]byte(io.Seed))
	return nil
}

func (io *NewBoardIO) GetSubmissionAddresses() []string {
	return io.submissionAddresses
}

func (io *NewBoardIO) GetConnections() []string {
	return io.connections
}

func (io *NewBoardIO) GetPK() cipher.PubKey {
	return io.pubKey
}

func (io *NewBoardIO) GetSK() cipher.SecKey {
	return io.secKey
}

type ThreadIO struct {
	BoardPubKey string `json:"board_public_key"`
	BoardSecKey cipher.SecKey
	ThreadRef   string `json:"thread_reference"`

	boardPubKey cipher.PubKey
	threadRef   skyobject.Reference
}

func (io *ThreadIO) Process() (e error) {
	io.boardPubKey, e = keys.GetPubKey(io.BoardPubKey)
	if e != nil {
		return
	}
	io.threadRef, e = keys.GetReference(io.ThreadRef)
	if e != nil {
		return
	}
	return
}

func (io *ThreadIO) GetBoardPK() cipher.PubKey {
	return io.boardPubKey
}

func (io *ThreadIO) GetThreadRef() skyobject.Reference {
	return io.threadRef
}

type NewThreadIO struct {
	BoardPubKey string `json:"board_public_key"`
	BoardSecKey cipher.SecKey
	UserPubKey  cipher.PubKey
	UserSecKey  cipher.SecKey
	Title       string `json:"title"`
	Body        string `json:"body"`

	boardPubKey cipher.PubKey
}

func (io *NewThreadIO) Process() (e error) {
	io.boardPubKey, e = keys.GetPubKey(io.BoardPubKey)
	if e != nil {
		return e
	}
	if e := CheckName(io.Title); e != nil {
		return e
	}
	if e := CheckDesc(io.Body); e != nil {
		return e
	}
	return nil
}

func (io *NewThreadIO) GetBoardPK() cipher.PubKey {
	return io.boardPubKey
}

type PostIO struct {
	ThreadIO
	PostRef string `json:"post_reference"`
	postRef skyobject.Reference
}

func (io *PostIO) Process() (e error) {
	if e = io.ThreadIO.Process(); e != nil {
		return
	}
	if io.postRef, e = keys.GetReference(io.PostRef); e != nil {
		return
	}
	return
}

func (io *PostIO) GetPostRef() skyobject.Reference {
	return io.postRef
}

type NewPostIO struct {
	NewThreadIO
	ThreadRef string `json:"thread_reference"`
	threadRef skyobject.Reference
}

func (io *NewPostIO) Process() (e error) {
	if e = io.NewThreadIO.Process(); e != nil {
		return
	}
	if io.threadRef, e = keys.GetReference(io.ThreadRef); e != nil {
		return
	}
	return
}

func (io *NewPostIO) GetThreadRef() skyobject.Reference {
	return io.threadRef
}
