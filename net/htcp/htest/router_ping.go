package htest

import (
	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
	"github.com/happylay-cloud/gf-extend/net/htcp/hnet"
)

type PingRouter struct {
	hnet.BaseTcpRouter
}

func (p *PingRouter) Handle(request hiface.ITcpRequest) {
	g.Log().Line(false).Debug("接收客户端消息：", request.GetHandlerRouter())

	err := request.GetConnection().SendTcpPkg(1, []byte("ping"))
	err = request.GetConnection().SendTcpPkg(1, []byte("ping"), true)
	if err != nil {
		g.Log().Line(false).Error(err)
	}
}
