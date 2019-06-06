package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)
//test 默认终止时间10m
func ClientTest()  {
	fmt.Println("Client test Starting...")
	time.Sleep(10*time.Second)
	conn,err:=net.Dial("tcp","127.0.0.1:7777")//回环地址是127.0.0.1
	if err!=nil{
		fmt.Println("Dial Server err ",err)
		return
	}

	for  {

		_,err:=conn.Write([]byte("hello zinx"))
		if err!=nil{
			fmt.Println("conn Write err ",err)
			return
		}
		buf:=make([]byte,512)
		n,err:=conn.Read(buf)
		if err!=nil{
			fmt.Println("conn Read err ",err)
			return
		}
		fmt.Printf("server call back %s, cnt%d\n",buf,n)
		time.Sleep(time.Second)
	}
}

func TestServer(t *testing.T)  {
	server:=NewServer("Zinx V0.1")
	go ClientTest()
	server.Serve()
}