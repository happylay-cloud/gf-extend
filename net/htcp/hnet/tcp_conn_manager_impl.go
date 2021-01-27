package hnet

import (
	"errors"
	"sync"

	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
)

// ConnTcpManager 连接管理模块
type ConnTcpManager struct {
	connections map[int64]hiface.ITcpConnection // 管理的连接信息
	connLock    sync.RWMutex                    // 读写连接的读写锁
}

// NewTcpConnManager 创建一个连接管理
func NewTcpConnManager() *ConnTcpManager {
	return &ConnTcpManager{
		connections: make(map[int64]hiface.ITcpConnection),
	}
}

// Add 添加连接
func (connMgr *ConnTcpManager) Add(conn hiface.ITcpConnection) {
	// 保护共享资源Map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	// 将conn连接添加到ConnTcpManager中
	connMgr.connections[conn.GetConnID()] = conn
	// 记录日志
	g.Log().Line(false).Info("添加连接到连接管理器成功：当前连接数", connMgr.Len())

}

// Remove 删除连接
func (connMgr *ConnTcpManager) Remove(conn hiface.ITcpConnection) {
	// 保护共享资源Map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	// 记录日志
	g.Log().Line(false).Info("删除连接 ConnID=", conn.GetConnID(), "成功：当前连接数", connMgr.Len())

}

// Get 利用ConnID获取连接
func (connMgr *ConnTcpManager) Get(connID int64) (hiface.ITcpConnection, error) {
	// 保护共享资源Map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("未发现Conn连接")
	}
}

// Len 获取当前连接
func (connMgr *ConnTcpManager) Len() int {
	return len(connMgr.connections)
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnTcpManager) ClearConn() {
	// 保护共享资源Map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}

	// 记录日志
	g.Log().Line().Info("成功清除所有连接：当前连接数", connMgr.Len())

}
