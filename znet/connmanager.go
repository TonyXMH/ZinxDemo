package znet

import (
	"errors"
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func (c*ConnManager)Add(conn ziface.IConnection)  {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[conn.GetConnID()]=conn
	fmt.Println("Connection connID",conn.GetConnID()," is added to ConnManager, num is",c.Len())
}

func (c*ConnManager)Remove(conn ziface.IConnection)  {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections,conn.GetConnID())
	fmt.Println("Connection connID",conn.GetConnID()," is removed from ConnManager, num is",c.Len())
}

func (c*ConnManager)GetConnection(connID uint32)(ziface.IConnection,error)  {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if con,ok:=c.connections[connID];ok{
		return con,nil
	}
	return nil,errors.New("Connection not found")
}

func (c*ConnManager)Len()int  {
	return len(c.connections)
}

func (c*ConnManager)ClearConn()  {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	for connID,conn:=range c.connections{
		conn.Stop()
		delete(c.connections,connID)
	}
	fmt.Println("Remove all connection num is",c.Len())
}


func NewConnManager()ziface.IConnManager  {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}