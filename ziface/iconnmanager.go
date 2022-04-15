package ziface

type IConnManager interface {
	Add(conn IConnection)
	Remove(connID uint32)
	Get(connID uint32) (IConnection, error)
	Len() int
	ClearAllConn()

	// hook function
	CallOnConnStart(IConnection)
	CallOnConnStop(IConnection)
	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(IConnection))
}
