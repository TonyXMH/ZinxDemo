package utils

import (
	"encoding/json"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"io/ioutil"
)

type GlobalObj struct {
	TcpServer     ziface.IServer
	Host          string
	TcpPort       int
	Name          string
	Version       string
	MaxPacketSize uint32
	MaxConn       int
}

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		TcpServer:     nil,
		Host:          "0.0.0.0",
		TcpPort:       7777,
		Name:          "ZinxServer",
		Version:       "v0.4",
		MaxPacketSize: 12000,
		MaxConn:       4096,
	}
}
