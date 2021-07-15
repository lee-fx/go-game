package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
	存储一切全局参数，供其他模块使用，数据来源可配置
*/
type GlobalObj struct {
	// 全局server对象
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version string
	// 当前服务器主机最大链接数
	MaxConn int
	// 数据包最大值
	MaxPackageSize uint32
	// worker池的数量（goroutine）
	WorkerPoolSize uint32
	// 框架允许开辟的最大worker
	MaxWorkerTaskLen uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 初始化当前GlobalObject对象
func init() {
	// 配置默认值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServer",
		Version:        "v0.10",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        365,
		MaxPackageSize: 4096,
		WorkerPoolSize: 8,      // worker工作池worker数量
		MaxWorkerTaskLen: 1024, // 每个worker对应queue的任务数量最大值
	}

	// 加载自定义参数
	GlobalObject.Reload()
}
