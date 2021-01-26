package ztest

import (
	"github.com/happylay-cloud/gf-extend/tcp/ziface"
	"github.com/happylay-cloud/gf-extend/tcp/zlog"
	"github.com/happylay-cloud/gf-extend/tcp/znet"
)

type HelloZinxRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter 路由
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	zlog.Debug("执行路由")

	// 先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("接收客户端消息：msgId=", request.GetMsgID(), "，data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello gf-plus"))
	if err != nil {
		zlog.Error(err)
	}
}
