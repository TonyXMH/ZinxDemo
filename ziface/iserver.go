package ziface

//服务抽象模块层

type IServer interface {
	//启动服务
	Start()
	//停止服务
	Stop()
	//开启业务服务
	Serve()

	AddRouter(msgID uint32, router IRouter)
}
