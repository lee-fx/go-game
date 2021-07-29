package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinx/fps_server/core"
	"zinx/fps_server/pb"
	"zinx/ziface"
	"zinx/znet"
)

/*
	世界聊天的路由业务
 */

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(req ziface.IRequest) {
	// 解析客户端传输的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(req.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error: ", err)
		return
	}
	// 获取哪个用户发送的
	pid, err := req.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println("Get Player ID is ERR: ", err)
		return
	}

	// 得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 将消息广播
	player.Talk(proto_msg.Content)

}
