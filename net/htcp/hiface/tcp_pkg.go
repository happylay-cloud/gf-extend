package hiface

type ITcpPkg interface {

	// 获取数据包消息头长度
	GetPkgHeadLen() int16
	// 获取数据包消息体内容
	GetPkgBody() []byte
	// 获取操作
	GetHandlerRouter()
}
