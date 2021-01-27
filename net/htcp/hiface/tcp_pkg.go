package hiface

import "io"

type ITcpPkg interface {

	// 获取数据包消息头长度
	GetPkgHeadLen() int16
	// 获取请求数据包
	GetPkg() []byte
	// 获取数据包消息体内容
	GetPkgBody() []byte
	// 获取路由（即消息处理器）
	GetHandlerRouter() string
	// 封包方法
	Pack() ([]byte, error)
	// 拆包方法
	Unpack(conn io.Reader) (ITcpPkg, error)
}
