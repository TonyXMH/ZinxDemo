package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/utils"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"net"
	"time"
)

//服务实现模块

type Server struct {
	Name        string //服务器名称
	IPVersion   string
	IP          string
	Port        int
	msgHandler  ziface.IMsgHandler
	connMgr     ziface.IConnManager
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("Server Obj %+v", *s)
	fmt.Printf("GlobalObject %+v", *utils.GlobalObject)
	fmt.Printf("[START] Server Listenner at IP:%s,Port:%d is starting\n", s.IP, s.Port)
	go func() {
		s.msgHandler.StartWorkerPool()
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen ", s.IPVersion, " err ", err)
			return
		}
		cid := uint32(0)
		fmt.Println("start Zinx server ", s.Name, "succ, now listenning...")
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			if s.connMgr.Len() > utils.GlobalObject.MaxConn {
				fmt.Println("Too many connection err ")
				conn.Close()
				continue
			}
			fmt.Println("conn Remote Addr ", conn.RemoteAddr().String())
			dealConn := NewConnection(conn, cid, s.msgHandler, s)
			cid++
			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP]Zinx Server,name ", s.Name)
	s.connMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}

}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.connMgr
}

func (s*Server)SetOnConnStart(hookFunc func(conn ziface.IConnection))  {
	s.OnConnStart = hookFunc
}

func (s*Server)SetOnConnStop(hookFunc func(conn ziface.IConnection))  {
	s.OnConnStop = hookFunc
}

func (s*Server)CallOnConnStart(conn ziface.IConnection)  {
	if s.OnConnStart !=nil{
		fmt.Println("CallOnConnStart")
		s.OnConnStart(conn)
	}

}

func (s*Server)CallOnConnStop(conn ziface.IConnection)  {
	if s.OnConnStop !=nil{
		fmt.Println("CallOnConnStop")
		s.OnConnStop(conn)
	}
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
		connMgr:    NewConnManager(),
	}
}
