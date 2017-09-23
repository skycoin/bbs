package websocket

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/net/skycoin-messenger/rpc"
)

type Client struct {
	rpc.Client

	seq uint32
	PendingMap

	conn *websocket.Conn
}

func (c *Client) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			c.Logger.Errorf("readLoop recovered err %v", err)
		}
		c.SetConnection(nil)
		c.conn.Close()
		close(c.Push)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.Logger.Errorf("error: %v", err)
			}
			c.Logger.Errorf("error: %v", err)
			return
		}
		if len(m) < msg.MSG_HEADER_END {
			return
		}
		c.Logger.Debugf("recv %x", m)
		opn := int(m[msg.MSG_OP_BEGIN])
		if opn == msg.OP_ACK {
			c.DelMsg(binary.BigEndian.Uint32(m[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END]))
			continue
		}
		op := msg.GetOP(opn)
		if op == nil {
			c.Logger.Errorf("op not found, %d", opn)
			return
		}

		c.ack(m[msg.MSG_OP_BEGIN:msg.MSG_SEQ_END])

		err = json.Unmarshal(m[msg.MSG_HEADER_END:], op)
		if err == nil {
			err = op.Execute(c)
			if err != nil {
				c.Logger.Errorf("websocket readLoop execute err: %v", err)
			}
		} else {
			c.Logger.Errorf("websocket readLoop json Unmarshal err: %v", err)
		}
		msg.PutOP(opn, op)
	}
}

func (c *Client) writeLoop() (err error) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if err := recover(); err != nil {
			c.Logger.Errorf("writeLoop recovered err %v", err)
		}
		ticker.Stop()
		c.SetConnection(nil)
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Push:
			c.Logger.Debug("Push", message)
			if !ok {
				c.Logger.Debug("closed c.Push")
				err = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					c.Logger.Error(err)
					return
				}
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				c.Logger.Error(err)
				return err
			}
			switch m := message.(type) {
			case *msg.PushMsg:
				err = c.write(w, msg.PUSH_MSG, m)
			case *msg.Reg:
				err = c.write(w, msg.PUSH_REG, m)
			default:
				c.Logger.Errorf("not implemented msg %v", m)
			}
			if err != nil {
				c.Logger.Error(err)
				return err
			}
			if err = w.Close(); err != nil {
				c.Logger.Error(err)
				return err
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error(err)
				return err
			}
		}
	}
}

func (c *Client) write(w io.WriteCloser, op byte, m interface{}) (err error) {
	_, err = w.Write([]byte{op})
	c.Logger.Debugf("op %d", op)
	if err != nil {
		return
	}
	ss := make([]byte, 4)
	nseq := atomic.AddUint32(&c.seq, 1)
	c.AddMsg(nseq, m)
	binary.BigEndian.PutUint32(ss, nseq)
	_, err = w.Write(ss)
	c.Logger.Debugf("seq %x", ss)
	if err != nil {
		return
	}
	jbs, err := json.Marshal(m)
	if err != nil {
		return
	}
	_, err = w.Write(jbs)
	c.Logger.Debugf("json %x", jbs)
	if err != nil {
		return
	}

	return nil
}

func (c *Client) ack(data []byte) error {
	data[msg.MSG_OP_BEGIN] = msg.PUSH_ACK
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
