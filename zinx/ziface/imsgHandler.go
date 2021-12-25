package ziface

/*
消息管理抽象层
*/

type IMsgHandle interface {
	// 调用/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// 开启工作池
	StartWorkerPool()
	// 启动Worker工作池
	StartOneWorker(workerID int, taskQueue chan IRequest)
	// 将消息发送给消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
}
