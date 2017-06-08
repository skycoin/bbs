package main

import (
	"github.com/skycoin/skycoin/src/cipher"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) != 2 {
		panic("invalid arguments")
	}
	pk, _ := cipher.GenerateDeterministicKeyPair([]byte(os.Args[1]))
	fmt.Print(pk.Hex())
}
