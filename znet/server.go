package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"github.com/pkg/errors"
	"net"
	"time"
)

//服务实现模块

type Server struct {
	Name      string //服务器名称
	IPVersion string
	IP        string
	Port      int
	Router    ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[START] Server Listenner at IP:%s,Port:%d is starting\n", s.IP, s.Port)
	go func() {
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
			fmt.Println("conn Remote Addr ", conn.RemoteAddr().String())
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP]Zinx Server,name ", s.Name)
}

func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}

}

func (s*Server)AddRouter(router ziface.IRouter)  {
	s.Router = router
	fmt.Println("AddRouter Successful")
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
		Router:nil,
	}
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("Conn Write err ", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}
