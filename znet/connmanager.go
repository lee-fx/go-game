package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	conns    map[uint32]ziface.IConn // 管理的连接信息
	connLock sync.RWMutex            // 读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		conns: make(map[uint32]ziface.IConn),
	}
}

func (cm *ConnManager) Add(conn ziface.IConn) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.conns[conn.GetConnID()] = conn
}

func (cm *ConnManager) Remove(conn ziface.IConn) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.conns, conn.GetConnID())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConn, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.conns[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("conn is not found!")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.conns)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除并停止conn
	for connID, conn := range cm.conns{
		conn.Stop()
		delete(cm.conns, connID)
	}

	fmt.Println("Clear All conn is succ! real conn num is :", cm.Len())
}
