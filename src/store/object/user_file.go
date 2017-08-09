package object

type UserFile struct {
	User User
	Seed string
}

func (f *UserFile) View() UserFileView {
	return UserFileView{
		User: UserView{
			User:   f.User,
			PubKey: f.User.PubKey.Hex(),
			SecKey: f.User.SecKey.Hex(),
		},
		Seed: f.Seed,
	}
}

type UserFileView struct {
	User UserView `json:"user"`
	Seed string   `json:"seed"`
}
