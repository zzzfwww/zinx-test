package main

import "zinx-test/zinx/znet"

func main() {
	// 创建server
	s := znet.NewServer("[zinx V0.2]")
	// 启动服务
	s.Serve()
}
