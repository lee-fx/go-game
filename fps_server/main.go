package main

import (
	"fmt"
	"zinx/fps_server/apis"
	"zinx/fps_server/core"
	"zinx/ziface"
	"zinx/znet"
)

func OnConnAdd(conn ziface.IConn) {
	// 创建玩家
	player := core.NewPlayer(conn)

	// 给可以端发送MsgID为1的消息
	player.SyncPid()

	// 给客户端发送MsgID:200的消息
	player.BroadCastStartPosition()

	// 将玩家加入世界管理模块
	core.WorldMgrObj.AddPlayer(player)

	// 将该链接绑定一个Pid 玩家的id属性
	conn.SetProperty("pid", player.Pid)

	// 同步周边玩家 告知他们当前玩家已经上线了
	player.SyncSurrounding()

	fmt.Println("==========> Player Pid = ", player.Pid, " is arrived <==============")
}

// 玩家下线业务
func OnConnLost(conn ziface.IConn) {
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println("OnConnLost Get Pid is error: ", err)
	}
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.OffLine()
	fmt.Println("=====> Player pid: ", pid, "OffLine <=====")
}

func main() {
	// 创建服务器句柄
	s := znet.NewServer()

	// 链接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnAdd)

	// 玩家断开链接之前触发的hook函数
	s.SetOnConnStop(OnConnLost)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})
	// 启动服务
	s.Server()
}
