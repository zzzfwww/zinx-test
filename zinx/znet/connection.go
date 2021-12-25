package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx-test/zinx/utils"
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
	// 无缓冲用于读写之间消息的传递
	msgChan chan []byte
	// 消息的管理msgId和对应的处理业务的api
	MsgHandler ziface.IMsgHandle
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running ...]")
	defer fmt.Println("[Reader is exit] connID = ", c.ConnID, "remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		// 创建一个拆包解包的对象
		dp := NewDataPack()
		//读取客户端的Msg Head 二进制流8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err", err)
			break
		}
		//拆包，得到msgId 和 msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据dataLen 再次读取Data，放在msg Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 根据绑定好的MsgID找到对应处理的api业务
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return err
	}
	c.msgChan <- binaryMsg
	return nil
}

// 写消息的goroutine 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running ...]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())
	// 不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

//启动链接，让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID=", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	// 启动从当前链接写数据业务
	go c.StartWriter()
}

// 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnId = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	//告知writer关闭
	c.ExitChan <- true
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
