package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
	// 消息队列
	TaskQueue []chan ziface.IRequest
	// worker池数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度/执行 对应的router
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " not found need register")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("repeat api exist, msgID = ", msgID)
		return
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, "succ!")
}

// 启动一个worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 启动一个worker
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 阻塞当前worker，等待消息从channel中传进来
		go mh.startWorker(i, mh.TaskQueue[i])
	}
}
// 启动一个worker工作流程
func (mh *MsgHandle) startWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID = ", workerID, "is started...")
	// 不断阻塞等待消息队列
	for{
		select {
		case req := <- taskQueue:
			// 执行业务
			mh.DoMsgHandler(req)
		}
	}
}

// 将请求交给TaskQueue 由worker来处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 平均分配给worker
	workerID := request.GetConn().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add connID = ", request.GetConn().GetConnID(), "request msgID = ", request.GetMsgID(), " to workerID = ", workerID)

	// 将消息发送给对应的taskQueue
	mh.TaskQueue[workerID] <- request
}

