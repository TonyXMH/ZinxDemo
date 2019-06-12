package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"github.com/TonyXMH/ZinxDemo/znet"
	"io"
	"net"
	"testing"
	"time"
)

//test 默认终止时间10m
func ClientTest(msgID uint32) {
	fmt.Println("Client test Starting...")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777") //回环地址是127.0.0.1
	if err != nil {
		fmt.Println("Dial Server err ", err)
		return
	}

	for i:=0;i<5;i++{
		dp := znet.NewDataPack()

		sendData, err := dp.Pack(znet.NewMsgPacket(msgID, []byte("v0.5 Client Send Message")))
		if err != nil {
			fmt.Println("dp.Pack err ", err)
			return
		}
		if _, err := conn.Write(sendData); err != nil {
			fmt.Println("Conn.Write err ", err)
			return
		}

		dataHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, dataHead); err != nil {
			fmt.Println("io.ReadFull err", err)
			return
		}

		msg, err := dp.Unpack(dataHead)
		if err != nil {
			fmt.Println("dp.Unpack err", err)
			return
		}
		var recvData []byte
		if msg.GetDataLen() > 0 {
			recvData = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn, recvData); err != nil {
				fmt.Println("io.ReadFull err ", err)
				return
			}
			msg.SetData(recvData)
		}

		fmt.Println("Received Data ", string(msg.GetData()))
		time.Sleep(time.Second)
	}
	time.Sleep(3*time.Second)
	conn.Close()
}

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Handle")
	fmt.Printf("MsgID %d,Msg Data%s\n", req.GetMsgID(), string(req.GetData()))
	if err := req.GetConnection().SendMsg(0, []byte("ping ping ping ...")); err != nil {
		fmt.Println(err)
	}

}

type HelloRouter struct {
	znet.BaseRouter
}

func (h *HelloRouter) Handle(req ziface.IRequest) {
	fmt.Println("HelloRouter Call")
	fmt.Println("MsgID ", req.GetMsgID(), " Msg Data ", string(req.GetData()))
	if err := req.GetConnection().SendMsg(1, []byte("hello hello")); err != nil {
		fmt.Println(err)
	}
}
func DoConnectionBegin(conn ziface.IConnection)  {
	fmt.Println("Connection Begin")
	if err:=conn.SendMsg(1,[]byte("Connection Begin"));err!=nil{
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection)  {
	fmt.Println("Connection Lost")
}

func TestServerV5(t *testing.T) {
	server := znet.NewServer("Zinx V0.5")
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)
	go ClientTest(0) //谁go出去都会有影响可以尝试一下将Serve go出去看看
	go ClientTest(1)
	//go ClientTest(0)
	//go ClientTest(1)
	server.Serve()
}
