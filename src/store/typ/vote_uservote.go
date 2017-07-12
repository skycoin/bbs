package typ

import (
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// UserVotes represents votes given to a user.
type UserVotes struct {
	User  cipher.PubKey
	Votes skyobject.References `skyobject:"schema=Vote"`
}

// UserVotesContainer contains the votes given to users.
type UserVotesContainer struct {
	Users []UserVotes
}

func (c *UserVotesContainer) GetUser(upk cipher.PubKey) *UserVotes {
	for i := range c.Users {
		if c.Users[i].User == upk {
			return &c.Users[i]
		}
	}
	c.Users = append(c.Users, UserVotes{User: upk})
	return &c.Users[len(c.Users)-1]
}

func (c *UserVotesContainer) RemoveUser(upk cipher.PubKey) {
	for i, u := range c.Users {
		if u.User == upk {
			c.Users[i], c.Users[0] = c.Users[0], c.Users[i]
			c.Users = c.Users[1:]
			return
		}
	}
}
