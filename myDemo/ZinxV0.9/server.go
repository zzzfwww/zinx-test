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
	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的数，再回写ping...ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(),
		", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

// hello
func (b *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	fmt.Println("recv from client: msgID=", request.GetMsgID(),
		", data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("hello Welcome"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建链接之后执行的钩子函数
func DoConnectionBegin(con ziface.IConnection) {
	fmt.Println("===> Do ConnectionBegin is Called ...")
	if err := con.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
}

// 链接端口之前需要执行的函数
func DoConnectionLost(con ziface.IConnection) {
	fmt.Println("===> Do DoConnectionLost is Called ...")
	fmt.Println("conn ID = ", con.GetConnID(), "is Lost...")
}

func main() {
	// 创建server
	s := znet.NewServer("[zinx V0.9]")
	// 注册链接hook钩子函数

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	// 给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// 启动服务
	s.Serve()
}
