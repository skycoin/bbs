package btp

import (
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
	"os"
)

// BoardInfo represents a board configuration.
type BoardInfo struct {
	Connection string `json:"connection"`
}

// MasterBoardInfo represents a master board configuration.
type MasterBoardInfo struct {
	SecretKey string              `json:"secret_key"`
	DependsOn map[string][]string `json:"depends_on"`
}

// BoardsFile saves configurations of boards.
type BoardsFile struct {
	unsaved      bool                        `json:"-"`
	Boards       map[string]*BoardInfo       `json:"boards,omitempty"`
	MasterBoards map[string]*MasterBoardInfo `json:"master_boards,omitempty"`
}

// NewBoardsFile creates a new BoardsFile.
func NewBoardsFile() *BoardsFile {
	return &BoardsFile{
		Boards:       make(map[string]*BoardInfo),
		MasterBoards: make(map[string]*MasterBoardInfo),
	}
}

// Loads loads BoardsFile from specified location.
func (f *BoardsFile) Load(loc string) error {
	return file.LoadJSON(loc, f)
}

// Save saves BoardsFile to specified location.
func (f *BoardsFile) Save(loc string) {
	file.SaveJSON(loc, *f, os.FileMode(0700))
	f.unsaved = false
}

// Add adds a board configuration to BoardsFile.
func (f *BoardsFile) Add(pk cipher.PubKey, connection string) error {
	if _, has := f.Boards[pk.Hex()]; has {
		return boo.Newf(boo.ObjectAlreadyExists,
			"configuration of board of public key '%s' already exists", pk.Hex())
	}
	f.Boards[pk.Hex()] = &BoardInfo{
		Connection: connection,
	}
	f.unsaved = true
	return nil
}

// Remove removes a board configuration from BoardsFile.
func (f *BoardsFile) Remove(pk cipher.PubKey) {
	delete(f.Boards, pk.Hex())
	f.unsaved = true
}

// AddMaster adds a master board configuration to BoardsFile.
func (f *BoardsFile) AddMaster(pk cipher.PubKey, sk cipher.SecKey) {
	f.MasterBoards[pk.Hex()] = &MasterBoardInfo{
		SecretKey: sk.Hex(),
		DependsOn: make(map[string][]string),
	}
	f.unsaved = true
}

// RemoveMaster removes a board configuration from BoardsFile.
func (f *BoardsFile) RemoveMaster(pk cipher.PubKey) {
	delete(f.MasterBoards, pk.Hex())
	f.unsaved = true
}

// AddMasterDep adds a dependency to a master board.
func (f *BoardsFile) AddMasterDep(masterPK, boardPK cipher.PubKey, threadRef skyobject.Reference) error {
	mbInfo, has := f.MasterBoards[masterPK.Hex()]
	if !has {
		return boo.Newf(boo.ObjectNotFound,
			"configuration for master board '%s' not found", masterPK.Hex())
	}
	bpkStr := boardPK.Hex()
	tDeps, _ := mbInfo.DependsOn[bpkStr]

	tRefStr := threadRef.String()
	for _, tDep := range tDeps {
		if tDep == tRefStr {
			return boo.Newf(boo.ObjectAlreadyExists,
				"dependency of b'%s':t'%s' already exists for board '%s'",
				boardPK.Hex(), threadRef.String(), masterPK.Hex())
		}
	}
	mbInfo.DependsOn[bpkStr] = append(tDeps, tRefStr)
	f.unsaved = true
	return nil
}

// RemoveMasterDep removes a dependency from a master board.
func (f *BoardsFile) RemoveMasterDep(masterPK, boardPK cipher.PubKey, threadRef skyobject.Reference) error {
	mbInfo, has := f.MasterBoards[masterPK.Hex()]
	if !has {
		return boo.Newf(boo.ObjectNotFound,
			"configuration for master board '%s' not found", masterPK.Hex())
	}
	pkStr := boardPK.Hex()
	tDeps, _ := mbInfo.DependsOn[pkStr]

	tRefStr := threadRef.String()
	for i, tDep := range tDeps {
		if tDep == tRefStr {
			tDeps = append(tDeps[:i], tDeps[i+1:]...)
			break
		}
	}
	if len(tDeps) == 0 {
		delete(mbInfo.DependsOn, boardPK.Hex())
	} else {
		mbInfo.DependsOn[pkStr] = tDeps
	}
	f.unsaved = true
	return nil
}
