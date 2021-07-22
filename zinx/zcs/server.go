package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call ping Router Handel...")
	// 先读取客户端的数据 再诙谐ping ping ping
	fmt.Println("recv from client msgID: ", request.GetMsgID(), ", data: ", string(request.GetData()))

	err := request.GetConn().SendMsg(200, []byte("ping..."))
	if err != nil {
		fmt.Println("err: ", err)
	}
}

// pihng test 自定义路由
type HelloZinx struct {
	znet.BaseRouter
}

func (p *HelloZinx) Handle(request ziface.IRequest) {
	fmt.Println("Call Hello Router Handel...")
	// 先读取客户端的数据 再诙谐ping ping ping
	fmt.Println("recv from client msgID: ", request.GetMsgID(), ", data: ", string(request.GetData()))

	err := request.GetConn().SendMsg(201, []byte("hello zinx..."))
	if err != nil {
		fmt.Println("err: ", err)
	}
}

// 创建链接之后hook
func DoConnBegin(conn ziface.IConn) {
	fmt.Println("=====> DoConnBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConn Begin")); err != nil {
		fmt.Println(err)
	}

	// 给当前链接设置属性
	conn.SetProperty("name", "leefx")
	conn.SetProperty("github", "https://github.com/lee-fx/go-im")

}

// 链接断开之前
func DoConnLost(conn ziface.IConn) {

	fmt.Println("=====> DoConnLost is Called...")

	// 获取链接属性
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("---------------- > name < : ", name)
	}

	if github, err := conn.GetProperty("github"); err == nil {
		fmt.Println("---------------- > github < : ", github)
	}

	if err := conn.SendMsg(203, []byte("DoConn Lost")); err != nil {
		fmt.Println(err)
	}

}

func main() {
	s := znet.NewServer()

	// 注册conn钩子
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinx{})

	s.Server()
}
