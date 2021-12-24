package znet

import (
	"fmt"
	"strconv"
	"zinx-test/zinx/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

//初始化创建的MsgHandl方法

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {

	// 1 从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "NOT FOUND!! need register")
		return
	}
	// 根据msgID调度router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id已经注册了
		panic("repeat api, msgId=" + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgId = ", msgID, "succ!!")
}
