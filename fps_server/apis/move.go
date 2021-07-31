package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinx/fps_server/core"
	"zinx/fps_server/pb"
	"zinx/ziface"
	"zinx/znet"
)

/*
	玩家移动
 */

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(req ziface.IRequest) {
	// 解析客户端传输过来的数据
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(req.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move Position Unmarshal error: ", err)
		return
	}

	// 获取发送数据的pid
	pid, err := req.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println("Move Get Pid error: ", err)
		return
	}
	fmt.Printf("Player pid: %d, move(%f, %f, %f, %f) \n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	// 给其他玩家同步移动(广播)
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
