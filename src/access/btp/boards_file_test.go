package btp

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"testing"
)

func generateBoards(bFile *BoardsFile, n int) ([]cipher.PubKey, error) {
	pkList := make([]cipher.PubKey, n)
	for i := 0; i < n; i++ {
		pk, _ := cipher.GenerateKeyPair()
		pkList[i] = pk
		if e := bFile.Add(pk, "127.0.0.1:1234"); e != nil {
			return nil, e
		}
	}
	return pkList, nil
}

func TestBoardsFile_Add(t *testing.T) {
	bFile := NewBoardsFile()
	if _, e := generateBoards(bFile, 10); e != nil {
		t.Error(e)
	}
	data, _ := json.MarshalIndent(*bFile, "", "    ")
	fmt.Println(string(data))
}

func TestBoardsFile_Remove(t *testing.T) {
	bFile := NewBoardsFile()
	pks, e := generateBoards(bFile, 10)
	if e != nil {
		t.Error(e)
	}
	for i := 0; i < 10; i++ {
		bFile.Remove(pks[i])
		data, _ := json.Marshal(bFile)
		fmt.Println(string(data))
		if len(bFile.Boards) != (9 - i) {
			t.Error("failed to remove", pks[i].Hex())
		}
	}
}
