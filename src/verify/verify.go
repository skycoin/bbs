package verify

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"reflect"
)

const (
	tagKey         = "verify"
	ignoreValue    = "-"
	signatureValue = "sig"
	publicKeyValue = "pk"
)

var (
	ErrInterfaceNotPointer       = errors.New("interface is not pointer")
	ErrInterfaceNotStructPointer = errors.New("interface is not struct pointer")
	ErrInvalidSignatureField     = errors.New("signature field has invalid type")
	ErrInvalidPublicKeyField     = errors.New("public key field has invalid type")
	ErrInput                     = errors.New("invalid input")
)

// Sign signs an object.
func Sign(obj interface{}, pk cipher.PubKey, sk cipher.SecKey) (cipher.Sig, error) {
	sig := cipher.Sig{}
	rVal := reflect.ValueOf(obj)

	// Check if interface is pointer.
	if rVal.Kind() != reflect.Ptr || rVal.IsNil() {
		return sig, ErrInterfaceNotPointer
	}

	rVal = rVal.Elem()

	// Check if pointer points to struct.
	if rVal.Kind() != reflect.Struct {
		return sig, ErrInterfaceNotStructPointer
	}

	rTyp := rVal.Type()
	sigInt := -1

	// Look through fields of struct.
	// Apply action based on tag.
	for i := 0; i < rTyp.NumField(); i++ {
		tagVal, has := rTyp.Field(i).Tag.Lookup(tagKey)
		if !has {
			continue
		}
		field := rVal.Field(i)
		switch tagVal {
		case ignoreValue:
			clearField(field)
		case signatureValue:
			if field.Type() != reflect.TypeOf(cipher.Sig{}) {
				return sig, ErrInvalidSignatureField
			}
			clearField(field)
			sigInt = i
		case publicKeyValue:
			if field.Type() != reflect.TypeOf(cipher.PubKey{}) {
				return sig, ErrInvalidPublicKeyField
			}
			field.Set(reflect.ValueOf(pk))
		default:
			return sig, errors.New("invalid tag")
		}
	}

	// Obtain signature.
	sig = cipher.SignHash(cipher.SumSHA256(encoder.Serialize(obj)), sk)
	if sigInt >= 0 {
		rVal.Field(sigInt).Set(reflect.ValueOf(sig))
	}
	return sig, nil
}

// Check checks the signature of object.
func Check(obj interface{}, pks ...cipher.PubKey) error {
	sig := cipher.Sig{}
	pk := cipher.PubKey{}
	switch len(pks) {
	case 1:
		pk = pks[0]
		fallthrough
	case 0:
		rVal := reflect.ValueOf(obj)

		// Check if interface is pointer.
		if rVal.Kind() != reflect.Ptr || rVal.IsNil() {
			return ErrInterfaceNotPointer
		}

		rVal = rVal.Elem()

		// Check if pointer points to struct.
		if rVal.Kind() != reflect.Struct {
			return ErrInterfaceNotStructPointer
		}

		rTyp := rVal.Type()

		for i := 0; i < rTyp.NumField(); i++ {
			tagVal, has := rTyp.Field(i).Tag.Lookup(tagKey)
			if !has {
				continue
			}
			field := rVal.Field(i)
			switch tagVal {
			case ignoreValue:
				clearField(field)
			case signatureValue:
				if field.Type() != reflect.TypeOf(cipher.Sig{}) {
					return ErrInvalidSignatureField
				}
				sig = field.Interface().(cipher.Sig)
				clearField(field)
			case publicKeyValue:
				if field.Type() != reflect.TypeOf(cipher.PubKey{}) {
					return ErrInvalidPublicKeyField
				}
				pk = field.Interface().(cipher.PubKey)
			}
		}
		// Check signature.
		return cipher.VerifySignature(pk, sig,
			cipher.SumSHA256(encoder.Serialize(obj)))

	default:
		return ErrInput
	}
}

func clearField(fv reflect.Value) {
	fv.Set(reflect.Zero(fv.Type()))
}
