package utils

import (
	"encoding/json"
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"io/ioutil"
	"os"
)

type GlobalObj struct {
	TcpServer        ziface.IServer
	Host             string
	TcpPort          int
	Name             string
	Version          string
	MaxPacketSize    uint32
	MaxConn          int
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	ConfFilePath     string
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (g *GlobalObj) Reload() {

	if ok, err := PathExists(g.ConfFilePath); !ok {
		fmt.Println(g.ConfFilePath, " Path is not exists err ", err)
		return
	}
	data, err := ioutil.ReadFile(g.ConfFilePath)
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
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          7777,
		Name:             "ZinxServer",
		Version:          "v0.4",
		MaxPacketSize:    12000,
		MaxConn:          4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "conf/zinx.json",
	}
	GlobalObject.Reload()
}
