package znet

import (
	"fmt"
	"net"
	"zinx-test/zinx/ziface"
)

// 链接模块
type Connection struct {
	// 当前链接的socketTCP 套接字
	Conn *net.TCPConn
	//链接id
	ConnID uint32
	//当前链接的状态
	isClosed bool
	// 告知当前链接已经退出的/停止 channel
	ExitChan chan bool
	//该链接处理的方法Router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exitremote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据到buf中，最大512个字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}
		//得到当前conn数据的request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		go func(request ziface.IRequest) {
			// 从路由中，找到注册绑定的Conn对应的router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

//启动链接，让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID=", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	// TODO 启动从当前链接写数据业务
}

// 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnId = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()

	close(c.ExitChan)
}

// 获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
