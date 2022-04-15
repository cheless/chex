package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	msgHandler  ziface.IMsgHandler
	connManager ziface.IConnManager
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[START] Server name: \"%s\",listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.Global.Version,
		utils.Global.MaxConn,
		utils.Global.MaxPacketSize)

	// 判断是否启动 WorkerPool 模式
	if utils.Global.WorkerPoolSize > 0 {
		s.msgHandler.StartWorkerPool()
	}

	// 获取 server 的 TCP Address
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp address failed:", err)
		return
	}

	// 监听该 TCP Address
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("Listen TCP failed:", err)
		return
	}
	fmt.Printf("start [%s] sucess\n", s.Name)

	// 阻塞等待客户端连接，处理连接
	var cid uint32 = 0
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept failed:", err)
		}

		// 如果当前 server 的 tcp 连接已达上限，则不再处理新的连接
		if s.connManager.Len() >= utils.Global.MaxConn {
			fmt.Println("====> Too Many Connections MaxConn = ", utils.Global.MaxConn)
			conn.Close()
			continue
		}

		// 将连接信息封装到 Connection 中，并绑定 Connection 和 Router
		c := newConnection(s, conn, cid, s.msgHandler)
		c.Start()
		cid++
	}
}

func (s *Server) Stop() {
	// TODO 将服务器用到的资源进行停止和回收
	fmt.Println("[STOP] zinx server, name:", s.Name)
	s.connManager.ClearAllConn()
}

// 启动服务
func (s *Server) Serve() {
	// 启动
	go s.Start()

	// TODO 启动服务之外的额外业务

	// 阻塞，因为 Start 是异步，防止 server直接结束
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.connManager
}

// NewServer 返回一个 TCP server engine
func NewServer(name string) *Server {
	return &Server{
		Name:        utils.Global.Name, //从全局参数获取
		IPVersion:   "tcp4",
		IP:          utils.Global.Host,    //从全局参数获取
		Port:        utils.Global.TcpPort, //从全局参数获取
		msgHandler:  newMsgHandler(),
		connManager: newConnManager(),
	}
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.connManager.SetOnConnStart(hookFunc)
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.connManager.SetOnConnStop(hookFunc)
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	s.connManager.CallOnConnStart(conn)
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	s.connManager.CallOnConnStop(conn)
}
