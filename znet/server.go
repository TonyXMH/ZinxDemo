package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"net"
	"time"
)

//服务实现模块

type Server struct {
	Name      string//服务器名称
	IPVersion string
	IP        string
	Port      int
}

func (s*Server)Start()  {
	fmt.Printf("[START] Server Listenner at IP:%s,Port:%d is starting\n",s.IP,s.Port)
	go func() {
		addr,err:=net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err!=nil{
			fmt.Println("resolve tcp addr err:",err)
			return
		}
		listener,err:=net.ListenTCP(s.IPVersion,addr)
		if err!=nil{
			fmt.Println("Listen ", s.IPVersion," err ",err)
			return
		}
		fmt.Println("start Zinx server ",s.Name,"succ, now listenning...")
		for{
			conn,err:=listener.AcceptTCP()
			if err!=nil{
				fmt.Println("Accept err",err)
				continue
			}

			go func() {
				for{
					buf:=make([]byte,512)
					cnt,err:=conn.Read(buf)
					if err!=nil{
						fmt.Println("recv buf err ",err)
						continue
					}
					if _,err:=conn.Write(buf[:cnt]);err!=nil{
						fmt.Println("write back buf err",err)
						continue
					}
				}
			}()
		}
	}()
}

func (s*Server)Stop()  {
	fmt.Println("[STOP]Zinx Server,name ",s.Name)
}

func (s*Server)Serve()  {
	s.Start()
	for   {
		time.Sleep(10*time.Second)
	}
	
}

func NewServer(name string)ziface.IServer  {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
}