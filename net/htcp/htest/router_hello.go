package htest

import (
	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
	"github.com/happylay-cloud/gf-extend/net/htcp/hnet"
)

type HelloRouter struct {
	hnet.BaseTcpRouter
}

func (h *HelloRouter) Handle(request hiface.ITcpRequest) {
	g.Log().Line(false).Debug("接收客户端消息：", request.GetHandlerRouter())
	err := request.GetConnection().SendTcpPkg(1, []byte("Hello gf-plus"))
	if err != nil {
		g.Log().Line(false).Error(err)
	}
}
