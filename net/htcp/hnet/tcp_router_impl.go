package hnet

import (
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
)

// 实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseTcpRouter struct{}

// 这里之所以BaseRouter的方法都为空，是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化
func (br *BaseTcpRouter) PreHandle(req hiface.ITcpRequest)  {}
func (br *BaseTcpRouter) Handle(req hiface.ITcpRequest)     {}
func (br *BaseTcpRouter) PostHandle(req hiface.ITcpRequest) {}
