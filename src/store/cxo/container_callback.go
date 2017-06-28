package cxo

import (
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type MsgMode int

const (
	RootFilled MsgMode = iota
	SubAccepted
	SubRejected

	ConnCreated
	ConnClosed
)

type Msg struct {
	pk cipher.PubKey
	c  *gnet.Conn
	m  MsgMode
}

func (m *Msg) PubKey() cipher.PubKey { return m.pk }
func (m *Msg) Conn() *gnet.Conn      { return m.c }
func (m *Msg) Mode() MsgMode         { return m.m }

func (c *Container) GetUpdatesChan() chan *Msg {
	c.Lock(c.GetUpdatesChan)
	defer c.Unlock()
	return c.msgs
}

func (c *Container) rootFilledInternalCB(root *node.Root) {
	log.Printf("[CONTAINER] Recieved filled board '%s'", root.Pub().Hex())
	go c.sendRootMsg(root.Pub(), RootFilled)
}

func (c *Container) subAcceptedInternalCB(conn *gnet.Conn, feed cipher.PubKey) {
	log.Printf("[CONTAINER] SUBSCRIPTION ACCEPTED: '%s (%s)'.",
		feed.Hex(), conn.Address())
	go c.sendRootMsg(feed, SubAccepted)
}

func (c *Container) subRejectedInternalCB(conn *gnet.Conn, feed cipher.PubKey) {
	log.Printf("[CONTAINER] SUBSCRIPTION REJECTED: '%s (%s)'.",
		feed.Hex(), conn.Address())
	go c.sendRootMsg(feed, SubRejected)
}

func (c *Container) connCreatedInternalCB(conn *gnet.Conn) {
	log.Printf("[CONTAINER] CONNECTION CREATED: '%s'.", conn.Address())
	go c.sendConnMsg(conn, ConnCreated)
}

func (c *Container) connClosedInternalCB(conn *gnet.Conn) {
	log.Printf("[CONTAINER] CONNECTION CLOSED: '%s'.", conn.Address())
	go c.sendConnMsg(conn, ConnClosed)
}

func (c *Container) sendRootMsg(pk cipher.PubKey, m MsgMode) {
	c.msgs <- &Msg{pk: pk, m: m}
}

func (c *Container) sendConnMsg(conn *gnet.Conn, m MsgMode) {
	c.msgs <- &Msg{c: conn, m: m}
}
