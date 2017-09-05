package tag

import (
	"github.com/skycoin/skycoin/src/cipher"
	"testing"
)

func TestTransfer(t *testing.T) {
	t.Run("case_one", func(t *testing.T) {
		pk, sk := cipher.GenerateKeyPair()
		src := struct {
			XXX string
			XX  int
			x   bool
			A   string        `transfer:"a"`
			B   int           `transfer:"b"`
			C   cipher.PubKey `transfer:"c"`
			D   cipher.SecKey `transfer:"d"`
		}{
			XXX: "what",
			XX:  -1,
			x:   true,
			A:   "test test",
			B:   100,
			C:   pk,
			D:   sk,
		}
		dst := struct {
			Number int           `transfer:"b"`
			PK     cipher.PubKey `transfer:"c"`
			SK     cipher.SecKey `transfer:"d"`
			What   string        `transfer:"a"`
		}{}
		Transfer(&src, &dst)
		switch {
		case src.A != dst.What:
			t.Error("strings not equal", src.A, dst.What)
			fallthrough
		case src.B != dst.Number:
			t.Error("ints not equal", src.B, dst.Number)
			fallthrough
		case src.C != dst.PK:
			t.Error("public keys not equal", src.C, dst.PK)
			fallthrough
		case src.D != dst.SK:
			t.Error("secret keys not equal", src.D, dst.SK)
			fallthrough
		default:
			t.Log("OKAY!")
		}
	})
}
