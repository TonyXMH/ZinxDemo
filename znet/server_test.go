package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"net"
	"testing"
	"time"
)

//test 默认终止时间10m
func ClientTest() {
	fmt.Println("Client test Starting...")
	time.Sleep(10 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777") //回环地址是127.0.0.1
	if err != nil {
		fmt.Println("Dial Server err ", err)
		return
	}

	for {

		_, err := conn.Write([]byte("hello zinx"))
		if err != nil {
			fmt.Println("conn Write err ", err)
			return
		}
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn Read err ", err)
			return
		}
		fmt.Printf("server call back %s, cnt%d\n", buf, n)
		fmt.Println("conn Remote Addr ", conn.RemoteAddr().String())
		time.Sleep(time.Second)
	}
}

//func TestServer(t *testing.T) {
//	server := NewServer("Zinx V0.1")
//	go ClientTest()
//	server.Serve()
//}

type PingRouter struct {
	BaseRouter
}

func (p*PingRouter)PreHandle(req ziface.IRequest)  {
	fmt.Println("Call PreHandle")
	_,err:=req.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	if err!=nil{
		fmt.Println("Conn Write Err",err)
	}
}

func (p*PingRouter)Handle(req ziface.IRequest)  {
	fmt.Println("Call Handle")
	_,err:=req.GetConnection().GetTCPConnection().Write([]byte("ping ping ...\n"))
	if err!=nil{
		fmt.Println("Conn Write Err",err)
	}
}

func (p*PingRouter)PostHandle(req ziface.IRequest)  {
	fmt.Println("Call PostHandle")
	_,err:=req.GetConnection().GetTCPConnection().Write([]byte("after ping ...\n"))
	if err!=nil{
		fmt.Println("Conn Write Err",err)
	}
}

func TestServerV3(t*testing.T)  {
	server:=NewServer("Zinx V0.3")
	server.AddRouter(&PingRouter{})
	go ClientTest()//谁go出去都会有影响可以尝试一下将Serve go出去看看
	server.Serve()
}

