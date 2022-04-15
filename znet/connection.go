package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"github.com/cheless/chex/utils"
	"github.com/cheless/chex/ziface"
)

/*
	封装连接模块
*/

type Connection struct {
	server     ziface.IServer     // 当前连接所属的 server
	conn       *net.TCPConn       // 当前连接的 TCP socket
	connID     uint32             // 连接的ID
	isClosed   bool               // 连接的状态
	ExitChan   chan bool          // 通知当前连接退出的 channel
	msgHandler ziface.IMsgHandler // 绑定的 Message 管理模块
	msgChan    chan []byte        // 读写 routine 间传递数据，无缓冲

	propertySet  map[string]interface{} // 用户配置的连接属性
	propertyLock sync.RWMutex           // 连接属性读写锁
}

// 初始化连接模块的方法
func newConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHanler ziface.IMsgHandler) *Connection {
	c := &Connection{
		server:       server,
		conn:         conn,
		connID:       connID,
		isClosed:     false,
		ExitChan:     make(chan bool, 1),
		msgHandler:   msgHanler,
		msgChan:      make(chan []byte),
		propertySet:  make(map[string]interface{}),
		propertyLock: sync.RWMutex{},
	}
	c.server.GetConnManager().Add(c) // 在连接管理模块中添加该连接
	return c
}

// 启动读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exist!]")
	defer c.Stop()

	for {
		// 读取 Head
		dp := DataPack{}
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.conn, headData)
		if err != nil {
			if err == io.EOF { // 如果 client 关闭 Connection
				fmt.Printf("client%d existed...\n", c.connID)
				break
			}
			fmt.Println("read Head failed, err:", err)
			continue
		}
		// 读取 Data
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack failed, err:", err)
			continue
		}
		if msg.GetDataLen() > 0 {
			data := make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read MsgData failed, err", err)
				break
			}
			msg.SetData(data)
			fmt.Printf("recived from client%d: %s\n", c.connID, msg.GetData())
		}
		// 将 Connection 和 Data 封装到 Request 中
		r := Request{
			conn: c,
			msg:  msg,
		}

		// 如果开启了 WorkerPool 模式
		if utils.Global.WorkerPoolSize > 0 {
			// 将请求发送给 WorkerPool 处理
			c.msgHandler.SendMsgToTaskQueue(&r)
		} else {
			// 直接调用当前 Message 绑定的 Handler
			go c.msgHandler.DoMsgHandler(&r)
		}
	}
}

// 启动写业务方法：将在客户端读取的数据返回给客户端
func (c *Connection) StartWriter() {
	fmt.Println("Writer goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exist!]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.conn.Write(data); err != nil {
				fmt.Println("send data to client failed:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 启动，让当前连接开始准备工作
func (c *Connection) Start() {
	fmt.Printf("conn Start()... connID: %v\n", c.connID)

	// 读写分离
	go c.StartReader()
	go c.StartWriter()

	// 创建连接之后的 hook function
	c.server.GetConnManager().CallOnConnStart(c)
}

// 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Printf("conn Stop()... connID: %v\n", c.connID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭连接之前的 hook function
	c.server.GetConnManager().CallOnConnStop(c)

	// 关闭 socket
	c.conn.Close()

	// 通知 Writer 连接已关闭
	c.ExitChan <- false

	// 将当前连接删除
	c.server.GetConnManager().Remove(c.connID)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前连接绑定的 socket Connection
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.connID
}

// 获取客户端的TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// 封包并发送给 Writer
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is closed while sending message")
	}
	// 封包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack message", msgID, "error")
		return errors.New("Pack message failed")
	}
	// 发送给 Writer
	c.msgChan <- msg

	return nil
}

func (c *Connection) SetProperty(key string, val interface{}) {
	c.propertyLock.Lock()
	c.propertySet[key] = val
	c.propertyLock.Unlock()
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	if val, ok := c.propertySet[key]; ok {
		return val, nil
	}
	c.propertyLock.RUnlock()
	return nil, fmt.Errorf("key:%s is not in propertySet", key)
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	delete(c.propertySet, key)
	c.propertyLock.Unlock()
}
