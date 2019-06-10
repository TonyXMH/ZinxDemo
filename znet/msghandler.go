package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandler) DoMsgHandler(req ziface.IRequest) {
	handler,ok:=m.Apis[req.GetMsgID()]
	if !ok{
		fmt.Printf("MsgID %d has not bind Router\n",req.GetMsgID())
		return
	}
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgID]; ok {
		panic("repeated api msgID" + strconv.Itoa(int(msgID)))
	}
	m.Apis[msgID] = router
	fmt.Println("Add api msgID = ",msgID)
}
