package factory

import (
	"sync"

	"net"

	"errors"

	"io"

	log "github.com/sirupsen/logrus"
	cn "github.com/skycoin/net/conn"
	"github.com/skycoin/skycoin/src/cipher"
)

type transport struct {
	creator *MessengerFactory
	// manager
	managerConn *Connection

	// node
	net *MessengerFactory
	// node B 2 A
	conn *Connection

	// app
	appNet  net.Listener
	appConn net.Conn

	fromNode, toNode cipher.PubKey
	fromApp, toApp   cipher.PubKey

	fieldsMutex sync.RWMutex
}

func NewTransport(creator *MessengerFactory, fromNode, toNode, fromApp, toApp cipher.PubKey) *transport {
	t := &transport{
		creator:  creator,
		fromNode: fromNode,
		toNode:   toNode,
		fromApp:  fromApp,
		toApp:    toApp,
		net:      NewMessengerFactory(),
	}
	return t
}

func (t *transport) ListenAndConnect(address string) (conn *Connection, err error) {
	conn, err = t.net.connectUDPWithConfig(address, &ConnConfig{
		OnConnected: func(connection *Connection) {
			connection.Reg()
		},
		Creator: t.creator,
	})
	t.managerConn = conn
	return
}

func (t *transport) Connect(address, appAddress string) (err error) {
	conn, err := t.net.connectUDPWithConfig(address, &ConnConfig{
		OnConnected: func(connection *Connection) {
			connection.writeOP(OP_BUILD_APP_CONN_OK,
				&buildConnResp{
					FromNode: t.fromNode,
					Node:     t.toNode,
					FromApp:  t.fromApp,
					App:      t.toApp,
				})
		},
		Creator: t.creator,
	})
	if err != nil {
		return
	}
	t.fieldsMutex.Lock()
	t.conn = conn
	t.fieldsMutex.Unlock()

	appConn, err := net.Dial("tcp", appAddress)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case m, ok := <-conn.GetChanIn():
				if !ok {
					log.Debugf("node conn read err %v", err)
					return
				}
				log.Debugf("Connect from app udp %x", m)
				err = writeAll(appConn, m)
				log.Debugf("Connect to server tcp %x", m)
				if err != nil {
					log.Debugf("app conn write err %v", err)
					return
				}
			}
		}
	}()

	go func() {
		buf := make([]byte, cn.MAX_UDP_PACKAGE_SIZE-100)
		for {
			n, err := appConn.Read(buf)
			if err != nil {
				log.Debugf("app conn read err %v, %d", err, n)
				return
			}
			log.Debugf("Connect from server tcp %x", buf[:n])
			err = conn.Write(buf[:n])
			log.Debugf("Connect to app udp %x", buf[:n])
			if err != nil {
				log.Debugf("node conn write err %v", err)
				return
			}
		}
	}()
	return
}

func (t *transport) setUDPConn(conn *Connection) {
	t.fieldsMutex.Lock()
	t.conn = conn
	t.fieldsMutex.Unlock()
}

func (t *transport) ListenForApp(address string, fn func()) (err error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	t.fieldsMutex.Lock()
	t.appNet = ln
	t.fieldsMutex.Unlock()

	fn()

	conn, err := ln.Accept()
	if err != nil {
		return
	}
	log.Debug("ListenForApp accepted")

	t.fieldsMutex.Lock()
	t.appConn = conn
	tConn := t.conn
	t.fieldsMutex.Unlock()

	go func() {
		buf := make([]byte, cn.MAX_UDP_PACKAGE_SIZE-100)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Debugf("app conn read err %v, %d", err, n)
				return
			}
			log.Debugf("ListenForApp from app tcp %x", buf[:n])
			err = tConn.Write(buf[:n])
			log.Debugf("ListenForApp write to node b %x", buf[:n])
			if err != nil {
				log.Debugf("node conn write err %v", err)
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case m, ok := <-tConn.GetChanIn():
				if !ok {
					err = errors.New("transport closed")
					return
				}
				log.Debugf("ListenForApp from node B udp %x", m)
				err = writeAll(conn, m)
				log.Debugf("ListenForApp write to app %x", m)
				if err != nil {
					return
				}
			}
		}
	}()
	return
}

func (t *transport) OK() {

}

func writeAll(conn io.Writer, m []byte) error {
	for i := 0; i < len(m); {
		n, err := conn.Write(m[i:])
		if err != nil {
			return err
		}
		i += n
	}
	return nil
}
