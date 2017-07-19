package hide

import (
	"encoding/json"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// Create data.
	type Post struct {
		Name string `json:"name"`
		Desc string `json:"description"`
	}
	key := "testKey123"
	post := Post{"Test", "It's testing."}

	{
		// Prepare data.
		data, e := json.Marshal(post)
		if e != nil {
			t.Error(e)
		}

		// Encrypt.
		hiddenData, e := Encrypt([]byte(key), data)
		if e != nil {
			t.Error(e)
		}
		t.Log("ENCRYPTED LENGTH:", len(hiddenData))

		// Decrypt.
		rawData, e := Decrypt([]byte(key), hiddenData)
		if e != nil {
			t.Error(e)
		}

		// Obtain data.
		post2 := Post{}
		if e := json.Unmarshal(rawData, &post2); e != nil {
			t.Error(e)
		}
		t.Log("OBTAINED:", post2)
	}

	{
		// Prepare data.
		data := encoder.Serialize(post)

		// Encrypt.
		hiddenData, e := Encrypt([]byte(key), data)
		if e != nil {
			t.Error(e)
		}
		t.Log("ENCRYPTED LENGTH:", len(hiddenData))

		// Decrypt.
		rawData, e := Decrypt([]byte(key), hiddenData)
		if e != nil {
			t.Error(e)
		}

		// Obtain data.
		post2 := Post{}
		if e := encoder.DeserializeRaw(rawData, &post2); e != nil {
			t.Error(e)
		}
		t.Log("OBTAINED:", post2)
	}
}
