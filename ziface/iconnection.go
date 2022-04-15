package ziface

import "net"

// 定义链接模块的抽象层
type IConnection interface {
	Start()                                  // 启动链接，让当前链接开始准备工作
	Stop()                                   // 停止链接，结束当前连接的工作
	GetTCPConnection() *net.TCPConn          // 获取当前链接绑定的 socket connection
	GetConnID() uint32                       // 获取当前链接模块的链接ID
	RemoteAddr() net.Addr                    // 获取客户端的 TCP状态 IP Port
	SendMsg(msgID uint32, data []byte) error // 封包并发送数据给 Writer

	SetProperty(key string, val interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}

// 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
