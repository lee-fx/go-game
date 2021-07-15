package ziface

// TODO 定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Server()
	// 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(uint32, IRouter)
	GetConnMgr() IConnManager

	// hook函数
	SetOnConnStart(func(conn IConn))
	SetOnConnStop(func(conn IConn))
	CallOnConnStart(conn IConn)
	CallOnConnStop(conn IConn)
}

