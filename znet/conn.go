package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

/*
	链接模块
*/
type Conn struct {
	// 当前conn隶属与哪个server
	TcpServer ziface.IServer

	// 当前链接的socket Tcp
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经退出的/停止chan (由reader告诉writer)
	ExitChan chan bool

	// 无缓冲读写goroutine通信
	msgChan chan []byte

	// 该链接处理的方法
	MsgHandler ziface.IMsgHandle

	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块方法
func NewConn(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Conn {
	c := &Conn{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandle,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}

	// 将conn加入ConnMgr中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 客户端链接的读业务
func (c *Conn) StartReader() {
	fmt.Println(" Reader Goroutine is runing...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remoteAddr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()

		// 读取客户端tlv格式的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				fmt.Println("read msg data error: ", err)
			}
		}
		msg.SetMsgData(data)

		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		// 判断是否开启工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			// 根据绑定好的MsgID 找到对应处理api业务执行
			go c.MsgHandler.DoMsgHandler(req)
		}

	}
}

func (c *Conn) StartWrite() {
	fmt.Println("[Writer Goretine is running]")
	defer fmt.Println(c.GetRemoteAddr(), "[conn Writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			// 有数据写入
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			// 关闭管道
			return
		}
	}
}

func (c *Conn) Start() {
	fmt.Println("conn start.. connID = ", c.ConnID)
	// 启动从当前链接的读数据业务
	go c.StartReader()
	// 写业务
	go c.StartWrite()

	// 调用创建conn之后钩子
	c.TcpServer.CallOnConnStart(c)
}

func (c *Conn) Stop() {
	fmt.Println("conn stop... connID = ", c.ConnID)
	// 当前链接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 调用销毁conn之前钩子
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前链接从connMgr中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭管道
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接绑定的socket conn
func (c *Conn) GetTCPConn() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Conn) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP状态 IP port
func (c *Conn) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 返回数据的封包方法
func (c *Conn) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id =", msgId)
		return errors.New("Pack error msg")
	}

	// 发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

// 设置链接属性
func (c *Conn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取链接属性
func (c *Conn) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found!")
}

// 移除链接属性
func (c *Conn) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
