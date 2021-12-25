package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx-test/zinx/ziface"
)

/*
链接管理模块
*/
type ConnManager struct {
	// 管理链接的集合
	connections map[uint32]ziface.IConnection
	// 保护链接的读写锁
	connLock sync.RWMutex
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}

}

// 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入到connManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID=", conn.GetConnID(), " add to ConnManager successfully: conn num=", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID=", conn.GetConnID(), " remove from ConnManager successfully: conn num=", connMgr.Len())
}

// 根据connID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}
}

// 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清楚并终止所有的链接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear Connection")
}
