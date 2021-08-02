package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"zinx/ziface"
	"zinx/zlog"
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

	/*
		config file path
	*/
	ConfFilePath string

	/*
		logger
	*/
	LogDir        string //日志所在文件夹 默认"./log"
	LogFile       string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {

	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		//fmt.Println("Config File ", g.ConfFilePath , " is not exist!!")
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

	//Logger 设置
	if g.LogFile != "" {
		zlog.SetLogFile(g.LogDir, g.LogFile)
	}
	if g.LogDebugClose == true {
		zlog.CloseDebug()
	}
}

//PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 初始化当前GlobalObject对象
func init() {

	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

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

		ConfFilePath:     pwd + "/conf/zinx.json",

		LogDir:           pwd + "/log",
		LogFile:          "",
		LogDebugClose:    false,
	}

	// 加载自定义参数
	GlobalObject.Reload()
}
