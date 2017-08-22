package tag

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"reflect"
)

const (
	verifyKey = "verify"
	verifyIgn = "-"
	verifySig = "sig"
	verifyPK  = "upk"
	verifyTS  = "time"
)

var (
	ErrInvalidSignatureField = errors.New("signature field has invalid type")
	ErrInvalidPublicKeyField = errors.New("public key field has invalid type")
	ErrInput                 = errors.New("invalid input")
)

// Sign signs an object.
func Sign(obj interface{}, pk cipher.PubKey, sk cipher.SecKey) {
	rVal, rTyp := getReflectPair(obj)
	sig := cipher.Sig{}
	sigInt := -1

	// Look through fields of struct.
	// Apply action based on tag.
	for i := 0; i < rTyp.NumField(); i++ {
		tagVal, has := rTyp.Field(i).Tag.Lookup(verifyKey)
		if !has {
			continue
		}
		field := rVal.Field(i)
		switch tagVal {
		case verifyIgn:
			clearField(field)
		case verifySig:
			if field.Type() != reflect.TypeOf(cipher.Sig{}) {
				panic(ErrInvalidSignatureField)
			}
			clearField(field)
			sigInt = i
		case verifyPK:
			if field.Type() != reflect.TypeOf(cipher.PubKey{}) {
				panic(ErrInvalidPublicKeyField)
			}
			field.Set(reflect.ValueOf(pk))
		case verifyTS:
			// TODO
		default:
			panic(errors.New("invalid tag"))
		}
	}

	// Obtain signature.
	sig = cipher.SignHash(cipher.SumSHA256(encoder.Serialize(obj)), sk)
	if sigInt >= 0 {
		rVal.Field(sigInt).Set(reflect.ValueOf(sig))
	}
}

// Verify checks the signature of object.
func Verify(obj interface{}, pks ...cipher.PubKey) error {
	sig := cipher.Sig{}
	pk := cipher.PubKey{}
	switch len(pks) {
	case 1:
		pk = pks[0]
		fallthrough
	case 0:
		rVal, rTyp := getReflectPair(obj)

		for i := 0; i < rTyp.NumField(); i++ {
			tagVal, has := rTyp.Field(i).Tag.Lookup(verifyKey)
			if !has {
				continue
			}
			field := rVal.Field(i)
			switch tagVal {
			case verifyIgn:
				clearField(field)
			case verifySig:
				if field.Type() != reflect.TypeOf(cipher.Sig{}) {
					return ErrInvalidSignatureField
				}
				sig = field.Interface().(cipher.Sig)
				clearField(field)
			case verifyPK:
				if field.Type() != reflect.TypeOf(cipher.PubKey{}) {
					return ErrInvalidPublicKeyField
				}
				pk = field.Interface().(cipher.PubKey)
			case verifyTS:
				// TODO
			}
		}
		// Verify signature.
		return cipher.VerifySignature(pk, sig,
			cipher.SumSHA256(encoder.Serialize(obj)))

	default:
		return ErrInput
	}
}
