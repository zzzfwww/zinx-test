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

// Test PreHandle
func (b *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test Hanle
func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ... ping ... ping...\n"))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

// Test PostHandle
func (b *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping ...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}

}

func main() {
	// 创建server
	s := znet.NewServer("[zinx V0.3]")
	// 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 启动服务
	s.Serve()
}
