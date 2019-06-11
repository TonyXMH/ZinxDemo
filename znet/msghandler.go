package znet

import (
	"fmt"
	"github.com/TonyXMH/ZinxDemo/utils"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest //chan的切片  每个chan都是缓冲chan
}

func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandler) DoMsgHandler(req ziface.IRequest) {
	handler, ok := m.Apis[req.GetMsgID()]
	if !ok {
		fmt.Printf("MsgID %d has not bind Router\n", req.GetMsgID())
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
	fmt.Println("Add api msgID = ", msgID)
}

func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("workerID ", workerID, " is running")
	for {
		select {
		case req := <-taskQueue:
			fmt.Println("taskQueue is coming req")
			m.DoMsgHandler(req)
		}
	}
}

func (m *MsgHandler) StartWorkerPool() {
	//for i, taskQueue := range m.TaskQueue {
	//	taskQueue = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
	//	go m.StartOneWorker(i, taskQueue)
	//}//for range的坑留作警示
	for i := 0; i < int(utils.GlobalObject.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandler) SendMsgToTaskQueue(req ziface.IRequest) {
	workerID := req.GetConnection().GetConnID() % utils.GlobalObject.WorkerPoolSize
	fmt.Println("ConnID ", req.GetConnection().GetConnID(), "Add MsgID", req.GetMsgID(), "To WorkerID ", workerID)
	m.TaskQueue[workerID] <- req
}
