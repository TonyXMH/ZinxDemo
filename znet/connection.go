package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"github.com/pkg/errors"
	"io"
	"net"
)

type Connection struct {
	Conn         *net.TCPConn
	ConnID       uint32
	MsgHandler   ziface.IMsgHandler
	IsClosed     bool
	ExitBuffChan chan bool

}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		MsgHandler:   msgHandler,
		IsClosed:     false,
		ExitBuffChan: make(chan bool, 1),

	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running")
	defer c.Stop()
	defer fmt.Println(c.RemoteAddr().String(), " conn reader is exit.")
	for {

		dp := NewDataPack()
		dataHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), dataHead); err != nil {
			fmt.Println("io.ReadFull err", err)
			c.ExitBuffChan <- true
			continue
		}
		msg, err := dp.Unpack(dataHead)
		if err != nil {
			fmt.Println("dp.Unpack err", err)
			c.ExitBuffChan <- true
			continue
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("io.ReadFull err", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		request := &Request{
			conn: c,
			msg:  msg,
		}
		go c.MsgHandler.DoMsgHandler(request)
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Conn Close err ", err)
		return
	}
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("Connection is already closed")
	}
	msg := &Message{
		DataLen: uint32(len(data)),
		ID:      msgID,
		Data:    data,
	}
	dp := NewDataPack()
	sendData, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("dp.Pack err ", err, " msgID", msg.ID)
		return errors.New("Pack msg err")
	}
	if _, err := c.Conn.Write(sendData); err != nil {
		c.ExitBuffChan <- true
		return err
	}
	return nil
}
