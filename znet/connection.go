package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"net"
)

type Connection struct {
	Conn         *net.TCPConn
	ConnID       uint32
	HandleAPI    ziface.HandleFunc
	IsClosed     bool
	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, handleAPI ziface.HandleFunc) ziface.IConnection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		HandleAPI:    handleAPI,
		IsClosed:     false,
		ExitBuffChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running")
	defer c.Stop()
	defer fmt.Println(c.RemoteAddr().String(), " conn reader is exit.")
	for {
		buf := make([]byte, 512)
		n, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Conn Read err ", err)
			c.ExitBuffChan <- true
			continue
		}
		if err = c.HandleAPI(c.Conn, buf, n); err != nil {
			fmt.Println("HandleAPI err ", err)
			c.ExitBuffChan <- true
			return
		}
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