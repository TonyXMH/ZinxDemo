package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func ServerTest() {
	fmt.Println("ServerTest")
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("net Listen err", err)
		return
	}
	fmt.Println("Listen successful")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err ", err)
		}
		go func(conn net.Conn) {
			dp := NewDataPack()
			for {
				dataHead := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, dataHead)
				if err != nil {
					fmt.Println("ReadFull err", err)
					break
				}
				msgHead, err := dp.Unpack(dataHead)
				if err != nil {
					fmt.Println("dp.Unpack err", err)
					break
				}
				if msgHead.GetDataLen() > 0 {
					msg := msgHead.(*Message)
					fmt.Printf("msgHead addr %p\n", msgHead)
					fmt.Printf("msg addr %p\n", msg)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("ReadFull err", err)
						return
					}
					fmt.Printf("==> Msg:ID=%d Msg:Len:%d,Msg:data:%s\n", msg.ID, msg.DataLen, msg.Data)
				}
			}
		}(conn)
	}

}

func ClientTest() {
	fmt.Println("ClientTest")
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Dial err ", err)
		return
	}
	fmt.Println("Client Dial Successfully!")
	dp := NewDataPack()
	msg1 := &Message{
		ID:      0,
		DataLen: 5,
		Data:    []byte("hello"),
	}
	data1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("dp.Pack err ", err)
		return
	}
	msg2 := &Message{
		ID:      1,
		DataLen: 7,
		Data:    []byte("world!!"),
	}
	data2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("Dial err ", err)
		return
	}
	data1 = append(data1, data2...)
	_, err = conn.Write(data1)
	if err != nil {
		fmt.Println("conn.Write err ", err)
		return
	}
	select {}
}

func TestMessage(t *testing.T) {
	go ServerTest()
	ClientTest()

}
