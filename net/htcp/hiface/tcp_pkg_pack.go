package hiface

// ITcpPakPack 封包数据、拆包数据、获取路由主题或方法接口抽象
type ITcpPakPack interface {
	// 获取路由（即消息处理器）
	GetHandlerRouter() string
	// 封包方法
	Pack()
	// 拆包方法
	Unpack()
}
