package ziface

type IMsgHandle interface {
	// 调度/执行 对应的router
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
	// 启动worker池
	StartWorkerPool()
	// 将请求交给TaskQueue 由worker来处理
	SendMsgToTaskQueue(IRequest)
}
