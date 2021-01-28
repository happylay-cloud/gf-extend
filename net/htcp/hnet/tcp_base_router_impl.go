package hnet

import (
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
)

// 实现ITcpRouter接口，自定义接口只需继承即可，或者继承后重写
type BaseTcpRouter struct{}

func (br *BaseTcpRouter) PreHandle(req hiface.ITcpRequest)  {}
func (br *BaseTcpRouter) Handle(req hiface.ITcpRequest)     {}
func (br *BaseTcpRouter) PostHandle(req hiface.ITcpRequest) {}
