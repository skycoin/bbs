package users

import "github.com/skycoin/bbs/src/store/object"

type FileView struct {
	User object.UserView `json:"user"`
	Seed string          `json:"seed"`
}

type File struct {
	User object.User
	Seed string
}

func (f *File) GenerateView() *FileView {
	return &FileView{
		User: object.UserView{
			User:      f.User,
			PublicKey: f.User.PublicKey.Hex(),
			SecretKey: f.User.SecretKey.Hex(),
		},
		Seed: f.Seed,
	}
}
