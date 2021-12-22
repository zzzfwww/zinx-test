package main

import (
	"fmt"
	"zinx-test/zinx/ziface"
	"zinx-test/zinx/znet"
)

/*
基于Zinx框架来开发的服务端应用程序
*/

type PingRouter struct {
	znet.BaseRouter
}

// Test Hanle
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数，再回写ping...ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(),
		", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 创建server
	s := znet.NewServer("[zinx V0.4]")
	// 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 启动服务
	s.Serve()
}
