package types

import "testing"

func TestNewUserConfigFromData(t *testing.T) {
	jsonData := `
	{
		"public_key": "03517b80b2889e4de80aae0fa2a4b2a408490f3178857df5b756e690b4524e1e61",
		"secret_key": "3cd98cc9385225f9af47e5ff0dfc073253aa410076cf5f426c19460a1d0de976"
	}`
	config, e := NewUserConfigFromData([]byte(jsonData))
	if e != nil {
		t.Error(e)
	}
	t.Logf("Public Key: %s", config.PublicKey.Hex())
	t.Logf("Secret Key: %s", config.SecretKey.Hex())
}
