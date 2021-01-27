package hiface

import "net"

type ITcpConnection interface {
	// 启动连接，让当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态
	Stop()
	// 从当前连接获取原始的Socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接ID（会话ID）
	GetConnID() int64
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 设置连接属性
	SetProperty(key string, value interface{})
	// 获取连接属性
	GetProperty(key string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(key string)
	// pkgBodyType 发送tcp数据包
	//  @pkgBodyType 数据类型
	//  @pkg     数据内容（自定义结构体对象）
	SendTcpPkg(pkgBodyType byte, pkg interface{}) error
}
