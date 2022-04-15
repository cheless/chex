package ziface

// 服务器接口
type IServer interface {
	Start()                                 // 启动
	Stop()                                  // 停止
	Serve()                                 // 运行
	AddRouter(msgID uint32, router IRouter) // 添加 Router
	GetConnManager() IConnManager
}
