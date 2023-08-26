package znet

import (
	"Zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnectionManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (connMgr *ConnectionManager) Add(conn ziface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

// Remove 删除连接
func (connMgr *ConnectionManager) Remove(conn ziface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

// Get 利用ConnID获取链接
func (connMgr *ConnectionManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// Len 获取当前连接数
func (connMgr *ConnectionManager) Len() int {
	return len(connMgr.connections)
}

// ClearConnection 清除并停止所有连接
func (connMgr *ConnectionManager) ClearConnection() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}
