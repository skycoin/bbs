package object

import "github.com/skycoin/bbs/src/store/object/revisions/r0"

type UserFile struct {
	User r0.User
	Seed string
}

func (f *UserFile) View() *UserFileView {
	return &UserFileView{
		User: r0.UserView{
			User:   f.User,
			PubKey: f.User.PubKey.Hex(),
			SecKey: f.User.SecKey.Hex(),
		},
		Seed: f.Seed,
	}
}

type UserFileView struct {
	User r0.UserView `json:"user"`
	Seed string      `json:"seed"`
}
