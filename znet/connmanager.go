package znet

import (
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	conns       map[uint32]ziface.IConnection
	connRWLock  sync.RWMutex
	onConnStart func(ziface.IConnection)
	onConnStop  func(ziface.IConnection)
}

func newConnManager() *ConnManager {
	return &ConnManager{
		conns:      make(map[uint32]ziface.IConnection),
		connRWLock: sync.RWMutex{},
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connRWLock.Lock()
	cm.conns[conn.GetConnID()] = conn
	cm.connRWLock.Unlock()
	fmt.Println("connection add to ConnManager successfully: conn num =", cm.Len())
}

func (cm *ConnManager) Remove(connID uint32) {
	cm.connRWLock.Lock()
	delete(cm.conns, connID)
	cm.connRWLock.Unlock()
	fmt.Println("connection Remove ConnID =", connID, " successfully: conn num =", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connRWLock.RLock()
	if conn, ok := cm.conns[connID]; ok {
		return conn, nil
	}
	cm.connRWLock.RUnlock()
	return nil, fmt.Errorf("connID:%d not found", connID)
}

func (cm *ConnManager) Len() int {
	return len(cm.conns)
}

// 清除所有连接，通常在服务下线时使用，用于释放资源
func (cm *ConnManager) ClearAllConn() {
	cm.connRWLock.Lock()
	for connID, conn := range cm.conns {
		conn.Stop()
		delete(cm.conns, connID)
	}
	cm.connRWLock.Unlock()
	fmt.Println("Clear All Connections successfully: conn num =", cm.Len())
}

func (cm *ConnManager) SetOnConnStart(f func(ziface.IConnection)) {
	cm.onConnStart = f
}

func (cm *ConnManager) SetOnConnStop(f func(ziface.IConnection)) {
	cm.onConnStop = f
}

func (cm *ConnManager) CallOnConnStart(conn ziface.IConnection) {
	if cm.onConnStart != nil {
		cm.onConnStart(conn)
	}
}

func (cm *ConnManager) CallOnConnStop(conn ziface.IConnection) {
	if cm.onConnStart != nil {
		cm.onConnStop(conn)
	}
}
