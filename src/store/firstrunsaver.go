package store

import (
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/skycoin/src/util/file"
	"os"
	"path/filepath"
)

const (
	FirstRunSaverFileName   = "bbs_first_run.json"
	SkycoinCommunityAddress = "34.204.161.180:8210"
	SkycoinCommunityPubKey  = "03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742"
)

// FirstRunFile represents the layout of the first run configuration file.
type FirstRunFile struct {
	FirstRun bool `json:"first_run"`
}

// FirstRunSaver manages first run actions.
type FirstRunSaver struct {
	config *args.Config
	bs     *BoardSaver
	data   *FirstRunFile
}

func NewFirstRunSaver(config *args.Config, boardSaver *BoardSaver) (*FirstRunSaver, error) {
	frs := FirstRunSaver{
		config: config,
		bs:     boardSaver,
	}
	frs.load()
	if frs.data.FirstRun {
		// Subscribe to community board as default on first run.
		bpk, e := misc.GetPubKey(SkycoinCommunityPubKey)
		if e != nil {
			frs.save(false)
			return nil, e
		}
		frs.bs.Add(SkycoinCommunityAddress, bpk)
	}
	frs.save(true)
	return &frs, nil
}

func (s *FirstRunSaver) absConfigDir() string {
	return filepath.Join(s.config.ConfigDir(), FirstRunSaverFileName)
}

func (s *FirstRunSaver) load() {
	s.data = &FirstRunFile{FirstRun: true}
	if s.config.SaveConfig() {
		if e := file.LoadJSON(s.absConfigDir(), s.data); e != nil {
			s.data.FirstRun = true
		}
	} else {
		s.data.FirstRun = false
	}
}

func (s *FirstRunSaver) save(done bool) {
	// Don't save if specified.
	if !s.config.SaveConfig() {
		return
	}
	// If done, save first run as false, else true.
	s.data.FirstRun = !done
	file.SaveJSON(s.absConfigDir(), *s.data, os.FileMode(0700))
}
