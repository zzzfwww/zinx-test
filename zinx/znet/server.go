package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx-test/zinx/ziface"
)

// Iserver 服务的实现
type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

// 定义当前客户端链接的所绑定handler api（目前这个handler是写死的，以后优化应该由用户
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显的业务
	fmt.Println("[Conn Handle] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s, Port:%d,is starting\n", s.IP, s.Port)
	go func() {
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
			//将处理新链接额业务方法和conn进行绑定，得到我们链接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			// 启动当前链接
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {

}
func (s *Server) Serve() {
	s.Start()
	select {}
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}
