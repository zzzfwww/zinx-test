package znet

import (
	"fmt"
	"net"
	"zinx-test/zinx/utils"
	"zinx-test/zinx/ziface"
)

// Iserver 服务的实现
type Server struct {
	Name      string
	IPVersion string
	IP        string
	// 服务监听的端口
	Port int
	// 当前server的消息管理模块，用来绑定msgid和对应的处理关系的业务api
	MsgHandler ziface.IMsgHandle
	// 该server的链接管理器
	ConnMgr ziface.IConnManager
	//该Server创建链接之后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该server销毁链接之前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s Listenner at IP: %s, Port:%d,is starting\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s MaxConn:%d MaxPackageSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		// 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ,", s.Name, "succ, Listening...")
		var cid uint32
		cid = 0
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 设置最大链接个数的判断，如果超过最大的链接，那么关闭此新链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端相应一个超出最大链接的错误包
				fmt.Println("---> Too Many Connection MaxConn", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//将处理新链接额业务方法和conn进行绑定，得到我们链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前链接
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	// 将一些服务器的资源，状态或者一些以及开辟的链接信息，进行停止或者回收
	fmt.Println("[STOP] Zinx server name=", s.Name)
	s.ConnMgr.ClearConn()
}
func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
}

// 注册OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart ----")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop ----")
		s.OnConnStop(conn)
	}
}
