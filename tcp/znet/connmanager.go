package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/happylay-cloud/gf-extend/tcp/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息
	connLock    sync.RWMutex                  // 读写连接的读写锁
}

// NewConnManager 创建一个连接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("添加连接到连接管理器成功：当前连接数", connMgr.Len())
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("删除连接 ConnID=", conn.GetConnID(), "成功：当前连接数", connMgr.Len())
}

// Get 利用ConnID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("未发现连接")
	}
}

// Len 获取当前连接
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("成功清除所有连接：当前连接数", connMgr.Len())
}
