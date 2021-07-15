package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// iserver 接口实现
type Server struct {
	// 服务器的名称
	Name string
	// 服务器的IP版本
	IPVersion string
	// 服务器的IP地址
	IP string
	// 服务器的端口
	Port int
	// 绑定MsgID和对应的处理业务的api
	MsgHandler ziface.IMsgHandle

	// server的链接管理器
	ConnMgr ziface.IConnManager

	OnConnStart func(conn ziface.IConn)
	OnConnStop  func(conn ziface.IConn)
}

func (s Server) Start() {
	fmt.Printf("[Zinx] Server Name：%s，listenner at IP：%s, Port：%d is strting\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version：%s， MaxConn：%d , MaxPackeetSize：%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	go func() {
		// 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

		// 获取一个tcp Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err: ", err)
			return
		}
		fmt.Println("Start Zinx server is Success, ", s.Name, " Success, Listenning...")

		var cid uint32
		cid = 0

		// 阻塞的等待
		for {
			// 请求链接进来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err: ", err)
				continue
			}

			// 判断链接数是否超出最大连接数
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应链接过多导致此次链接失败
				fmt.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConn(&s, conn, cid, s.MsgHandler)
			cid++

			// 启动链接模块
			go dealConn.Start()

		}
	}()

	// 处理客户端读写任务
}

func (s Server) Stop() {
	// 将服务器资源、状态、或者一些已经已经开辟的链接信息，进行停止或者回收
	fmt.Println("[SERVER] IS STOP, resource clear!")
	s.ConnMgr.ClearConn()
}

func (s Server) Server() {
	// 启动服务器
	s.Start()

	// TODO 做一些启动服务器之后的操作

	// 阻塞
	select {}
}

// 添加路由方法
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Scuu!")
}

// 获取conn方法
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 初始化 server 模块
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConn)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConn)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConn) {
	if s.OnConnStart != nil {
		fmt.Println("----> call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConn) {
	if s.OnConnStop != nil {
		fmt.Println("----> call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
