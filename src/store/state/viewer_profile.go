package state

type Profile struct {
	Trusted      map[string]struct{}
	MarkedAsSpam map[string]struct{}
	Blocked      map[string]struct{}

	TrustedBy      map[string]struct{}
	MarkedAsSpamBy map[string]struct{}
	BlockedBy      map[string]struct{}
}

func NewProfile() *Profile {
	return &Profile{
		Trusted:        make(map[string]struct{}),
		MarkedAsSpam:   make(map[string]struct{}),
		Blocked:        make(map[string]struct{}),
		TrustedBy:      make(map[string]struct{}),
		MarkedAsSpamBy: make(map[string]struct{}),
		BlockedBy:      make(map[string]struct{}),
	}
}

type ProfileView struct {
	TrustedCount      int      `json:"trusted_count"`
	Trusted           []string `json:"trusted"`
	MarkedAsSpamCount int      `json:"marked_as_spam_count"`
	MarkedAsSpam      []string `json:"marked_as_spam"`
	BlockedCount      int      `json:"blocked_count"`
	Blocked           []string `json:"blocked"`

	TrustedByCount      int      `json:"trusted_by_count"`
	TrustedBy           []string `json:"trusted_by"`
	MarkedAsSpamByCount int      `json:"marked_as_spam_by_count"`
	MarkedAsSpamBy      []string `json:"marked_as_spam_by"`
	BlockedByCount      int      `json:"blocked_by_count"`
	BlockedBy           []string `json:"blocked_by"`
}

func (p *Profile) View() *ProfileView {
	view := &ProfileView{
		TrustedCount:        len(p.Trusted),
		Trusted:             make([]string, len(p.Trusted)),
		MarkedAsSpamCount:   len(p.MarkedAsSpam),
		MarkedAsSpam:        make([]string, len(p.MarkedAsSpam)),
		BlockedCount:        len(p.Blocked),
		Blocked:             make([]string, len(p.Blocked)),
		TrustedByCount:      len(p.TrustedBy),
		TrustedBy:           make([]string, len(p.TrustedBy)),
		MarkedAsSpamByCount: len(p.MarkedAsSpamBy),
		MarkedAsSpamBy:      make([]string, len(p.MarkedAsSpamBy)),
		BlockedByCount:      len(p.BlockedBy),
		BlockedBy:           make([]string, len(p.BlockedBy)),
	}

	i := 0
	for k := range p.Trusted {
		view.Trusted[i] = k
		i++
	}

	i = 0
	for k := range p.MarkedAsSpam {
		view.MarkedAsSpam[i] = k
		i++
	}

	i = 0
	for k := range p.Blocked {
		view.Blocked[i] = k
		i++
	}

	i = 0
	for k := range p.TrustedBy {
		view.TrustedBy[i] = k
		i++
	}

	i = 0
	for k := range p.MarkedAsSpamBy {
		view.MarkedAsSpamBy[i] = k
		i++
	}

	i = 0
	for k := range p.BlockedBy {
		view.BlockedBy[i] = k
		i++
	}

	return view
}

func (p *Profile) ClearVotesFor(user string) {
	delete(p.Trusted, user)
	delete(p.MarkedAsSpam, user)
	delete(p.Blocked, user)
}

func (p *Profile) ClearVotesBy(user string) {
	delete(p.TrustedBy, user)
	delete(p.MarkedAsSpamBy, user)
	delete(p.BlockedBy, user)
}
