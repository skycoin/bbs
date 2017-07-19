package hide

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
)

var (
	// ErrKeyTooLong occurs when key phrase is too long.
	ErrKeyTooLong = errors.New("key is too long, maximum length allowed is 32 bytes")
)

// Encrypt encrypts 'data' with key.
func Encrypt(key, data []byte) ([]byte, error) {
	gcm, e := obtainGCM(key)
	if e != nil {
		return nil, e
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, e := io.ReadFull(rand.Reader, nonce); e != nil {
		return nil, e
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts 'data' with key.
func Decrypt(key, data []byte) ([]byte, error) {
	gcm, e := obtainGCM(key)
	if e != nil {
		return nil, e
	}
	nonceSize := gcm.NonceSize()
	nonce, data := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, data, nil)
}

func obtainGCM(key []byte) (cipher.AEAD, error) {
	key, e := ensureKey32(key)
	if e != nil {
		return nil, e
	}
	block, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	return cipher.NewGCM(block)
}

func ensureKey32(key []byte) ([]byte, error) {
	keyLen := len(key)
	switch {
	case keyLen == 32:
		return key, nil
	case keyLen > 32:
		return nil, ErrKeyTooLong
	default:
		return append(key, make([]byte, 32-keyLen)...), nil
	}
}
