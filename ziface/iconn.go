package ziface

import "net"

// 定义链接模块的抽象层
type IConn interface {
	// 启动链接
	Start()
	// 停止链接
	Stop()
	// 获取当前链接绑定的socket conn
	GetTCPConn() *net.TCPConn
	// 获取当前链接模块的链接ID
	GetConnID() uint32
	// 获取远程客户端的 TCP状态 IP port
	GetRemoteAddr() net.Addr
	// 发送数据给远端客户端
	SendMsg(uint32, []byte) error

	// 设置链接属性
	SetProperty(key string, value interface{})
	// 获取链接属性
	GetProperty(key string) (interface{}, error)
	// 移除链接属性
	RemoveProperty(key string)
}

// 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
