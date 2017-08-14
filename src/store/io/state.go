package io

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type State struct {
	Node *node.Node
	Conn *gnet.Conn
	Root *skyobject.Root
}

type Changes struct {
	sync.Mutex
	BoardPubKey cipher.PubKey
	Threads     map[cipher.SHA256]*ThreadChanges
}

func NewChanges(bpk cipher.PubKey, enable bool) *Changes {
	if enable {
		return &Changes{
			BoardPubKey: bpk,
			Threads:     make(map[cipher.SHA256]*ThreadChanges),
		}
	} else {
		return nil
	}
}

func (c *Changes) getThreadChange(tRef cipher.SHA256) *ThreadChanges {
	c.Lock()
	defer c.Unlock()

	tChanges, has := c.Threads[tRef]
	if !has {
		tChanges = &ThreadChanges{
			ThreadRef: tRef,
		}
		c.Threads[tRef] = tChanges
	}
	return tChanges
}

func (c *Changes) RecordThreadVoteChanges(tRef cipher.SHA256, vs *object.VotesSummary) {
	if c == nil {
		return
	}
	c.getThreadChange(tRef).Do(func(tc *ThreadChanges) {
		tc.ThreadVotesChange = vs
	})
}

func (c *Changes) RecordPostVoteChanges(tRef cipher.SHA256, vs *object.VotesSummary) {
	if c == nil {
		return
	}
	c.getThreadChange(tRef).Do(func(tc *ThreadChanges) {
		tc.PostVotesChanges = append(tc.PostVotesChanges, vs)
	})
}

func (c *Changes) RecordNewPost(tRef cipher.SHA256, post *object.Content) {
	if c == nil {
		return
	}
	c.getThreadChange(tRef).Do(func(tc *ThreadChanges) {
		tc.NewPosts = append(tc.NewPosts, post)
	})
}

func (c *Changes) RecordDeleteThread(tRef cipher.SHA256) {
	if c == nil {
		return
	}
	c.getThreadChange(tRef).Do(func(tc *ThreadChanges) {
		tc.ThreadDeleted = true
	})
}

func (c *Changes) RecordDeletePost(tRef, pRef cipher.SHA256) {
	if c == nil {
		return
	}
	c.getThreadChange(tRef).Do(func(tc *ThreadChanges) {
		tc.DeletedPosts = append(tc.DeletedPosts, pRef)
	})
}

type ThreadChanges struct {
	sync.Mutex

	ThreadRef         cipher.SHA256
	ThreadDeleted     bool
	ThreadVotesChange *object.VotesSummary

	PostVotesChanges []*object.VotesSummary
	NewPosts         []*object.Content
	DeletedPosts     []cipher.SHA256
}

func (c *ThreadChanges) Do(action func(tc *ThreadChanges)) {
	c.Lock()
	defer c.Unlock()
	action(c)
}
