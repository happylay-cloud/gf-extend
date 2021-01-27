package hiface

// ITcpConnManager Tcp连接管理抽象层
type ITcpConnManager interface {
	// 添加链接
	Add(conn ITcpConnection)
	// 删除连接
	Remove(conn ITcpConnection)
	// 利用ConnID获取链接
	Get(connID int64) (ITcpConnection, error)
	// 获取当前连接
	Len() int
	// 删除并停止所有连接
	ClearConn()
}
