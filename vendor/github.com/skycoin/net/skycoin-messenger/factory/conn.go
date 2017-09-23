package factory

import (
	"sync"

	"encoding/json"
	"errors"

	"github.com/skycoin/net/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

type Connection struct {
	*factory.Connection
	factory     *MessengerFactory
	key         cipher.PubKey
	keySetCond  *sync.Cond
	keySet      bool
	services    *NodeServices
	servicesMap map[cipher.PubKey]*Service
	fieldsMutex sync.RWMutex

	in           chan []byte
	disconnected chan struct{}

	proxyConnections map[uint32]*Connection

	appTransports      map[cipher.PubKey]*transport
	appTransportsMutex sync.RWMutex

	// callbacks

	// call after received response for FindServiceNodesByKeys
	findServiceNodesByKeysCallback func(resp *QueryResp)

	// call after received response for FindServiceNodesByAttributes
	findServiceNodesByAttributesCallback func(resp *QueryByAttrsResp)
}

// Used by factory to spawn connections for server side
func newConnection(c *factory.Connection, factory *MessengerFactory) *Connection {
	connection := &Connection{
		Connection:    c,
		factory:       factory,
		appTransports: make(map[cipher.PubKey]*transport),
	}
	c.RealObject = connection
	connection.keySetCond = sync.NewCond(connection.fieldsMutex.RLocker())
	return connection
}

// Used by factory to spawn connections for client side
func newClientConnection(c *factory.Connection, factory *MessengerFactory) *Connection {
	connection := newUDPClientConnection(c, factory)
	go func() {
		connection.preprocessor()
		close(connection.disconnected)
	}()
	return connection
}

// Used by factory to spawn connections for udp client side
func newUDPClientConnection(c *factory.Connection, factory *MessengerFactory) *Connection {
	connection := &Connection{
		Connection:       c,
		factory:          factory,
		in:               make(chan []byte),
		disconnected:     make(chan struct{}),
		proxyConnections: make(map[uint32]*Connection),
		appTransports:    make(map[cipher.PubKey]*transport),
	}
	c.RealObject = connection
	connection.keySetCond = sync.NewCond(connection.fieldsMutex.RLocker())
	return connection
}

func (c *Connection) setProxyConnection(seq uint32, conn *Connection) {
	c.fieldsMutex.Lock()
	c.proxyConnections[seq] = conn
	c.fieldsMutex.Unlock()
}

func (c *Connection) removeProxyConnection(seq uint32) (conn *Connection, ok bool) {
	c.fieldsMutex.Lock()
	conn, ok = c.proxyConnections[seq]
	if ok {
		delete(c.proxyConnections, seq)
	}
	c.fieldsMutex.Unlock()
	return
}

func (c *Connection) WaitForDisconnected() {
	<-c.disconnected
}

func (c *Connection) SetKey(key cipher.PubKey) {
	c.fieldsMutex.Lock()
	c.key = key
	c.keySet = true
	c.fieldsMutex.Unlock()
	c.keySetCond.Broadcast()
}

func (c *Connection) IsKeySet() bool {
	c.fieldsMutex.Lock()
	defer c.fieldsMutex.Unlock()
	return c.keySet
}

func (c *Connection) GetKey() cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	if !c.keySet {
		c.keySetCond.Wait()
	}
	return c.key
}

func (c *Connection) setServices(s *NodeServices) {
	if s == nil {
		c.fieldsMutex.Lock()
		c.services = nil
		c.servicesMap = nil
		c.fieldsMutex.Unlock()
		return
	}
	m := make(map[cipher.PubKey]*Service)
	for _, v := range s.Services {
		m[v.Key] = v
	}
	c.fieldsMutex.Lock()
	c.services = s
	c.servicesMap = m
	c.fieldsMutex.Unlock()
}

func (c *Connection) getService(key cipher.PubKey) (service *Service, ok bool) {
	c.fieldsMutex.Lock()
	defer c.fieldsMutex.Unlock()
	service, ok = c.servicesMap[key]
	return
}

func (c *Connection) GetServices() *NodeServices {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.services
}

func (c *Connection) Reg() error {
	return c.Write(GenRegMsg())
}

func (c *Connection) UpdateServices(ns *NodeServices) error {
	if ns == nil || len(ns.Services) < 1 {
		return errors.New("invalid arguments")
	}
	err := c.writeOP(OP_OFFER_SERVICE, ns)
	if err != nil {
		return err
	}
	c.setServices(ns)
	return nil
}

func (c *Connection) OfferService(attrs ...string) error {
	return c.UpdateServices(&NodeServices{Services: []*Service{{Key: c.GetKey(), Attributes: attrs}}})
}

func (c *Connection) OfferServiceWithAddress(address string, attrs ...string) error {
	return c.UpdateServices(&NodeServices{Services: []*Service{{Key: c.GetKey(), Attributes: attrs, Address: address}}})
}

func (c *Connection) FindServiceNodesByAttributes(attrs ...string) error {
	return c.writeOP(OP_QUERY_BY_ATTRS, newQueryByAttrs(attrs))
}

func (c *Connection) FindServiceNodesByKeys(keys []cipher.PubKey) error {
	return c.writeOP(OP_QUERY_SERVICE_NODES, newQuery(keys))
}

func (c *Connection) BuildConnection(node, app cipher.PubKey) error {
	return c.writeOP(OP_BUILD_APP_CONN, &appConn{Node: node, App: app})
}

func (c *Connection) Send(to cipher.PubKey, msg []byte) error {
	return c.Write(GenSendMsg(c.GetKey(), to, msg))
}

func (c *Connection) SendCustom(msg []byte) error {
	return c.writeOPBytes(OP_CUSTOM, msg)
}

func (c *Connection) preprocessor() (err error) {
	defer func() {
		if e := recover(); e != nil {
			c.GetContextLogger().Debugf("panic in preprocessor %v", e)
		}
		if err != nil {
			c.GetContextLogger().Debugf("preprocessor err %v", err)
		}
		c.Close()
	}()
	for {
		select {
		case m, ok := <-c.Connection.GetChanIn():
			if !ok {
				return
			}
			c.GetContextLogger().Debugf("read %x", m)
			if len(m) < MSG_HEADER_END {
				return
			}
			opn := m[MSG_OP_BEGIN]
			if opn&RESP_PREFIX > 0 {
				i := int(opn &^ RESP_PREFIX)
				r := getResp(i)
				if r != nil {
					body := m[MSG_HEADER_END:]
					if len(body) > 0 {
						err = json.Unmarshal(body, r)
						if err != nil {
							return
						}
					}
					err = r.Run(c)
					if err != nil {
						return
					}
					putResp(i, r)
					continue
				}
			}

			c.in <- m
		}
	}
}

func (c *Connection) GetChanIn() <-chan []byte {
	if c.in == nil {
		return c.Connection.GetChanIn()
	}
	return c.in
}

func (c *Connection) Close() {
	c.keySetCond.Broadcast()
	c.fieldsMutex.Lock()
	defer c.fieldsMutex.Unlock()
	if c.Connection.IsClosed() {
		return
	}
	if c.keySet {
		c.factory.unregister(c.key, c)
	}
	if c.in != nil {
		close(c.in)
		c.in = nil
	}
	c.Connection.Close()
}

func (c *Connection) writeOPBytes(op byte, body []byte) error {
	data := make([]byte, MSG_HEADER_END+len(body))
	data[MSG_OP_BEGIN] = op
	copy(data[MSG_HEADER_END:], body)
	return c.Write(data)
}

func (c *Connection) writeOP(op byte, object interface{}) error {
	js, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return c.writeOPBytes(op, js)
}

func (f *Connection) setTransport(to cipher.PubKey, tr *transport) {
	f.appTransportsMutex.Lock()
	defer f.appTransportsMutex.Unlock()

	f.appTransports[to] = tr
}

func (f *Connection) getTransport(to cipher.PubKey) (tr *transport, ok bool) {
	f.appTransportsMutex.RLock()
	defer f.appTransportsMutex.RUnlock()

	tr, ok = f.appTransports[to]
	return
}
