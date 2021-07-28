package main

import "zinx/znet"

func main() {
	// 创建服务器句柄
	s := znet.NewServer()

	// 链接创建和销毁的HOOK钩子函数

	// 注册一些路由业务

	// 启动服务
	s.Server()
}
