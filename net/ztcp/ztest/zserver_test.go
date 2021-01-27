package ztest

import (
	"testing"

	"github.com/happylay-cloud/gf-extend/net/ztcp/ziface"
	"github.com/happylay-cloud/gf-extend/net/ztcp/zlog"
	"github.com/happylay-cloud/gf-extend/net/ztcp/znet"
)

// 创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Debug("创建连接被执行... ")

	// 设置两个连接属性，在连接创建之后
	zlog.Debug("设置连接属性")
	conn.SetProperty("name", "gf-puls")

	err := conn.SendMsg(2, []byte("开始创建连接..."))
	if err != nil {
		zlog.Error(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	// 在连接销毁之前，查询conn的name属性
	if name, err := conn.GetProperty("name"); err == nil {
		zlog.Error("连接属性 name = ", name)
	}

	zlog.Debug("连接断开被执行... ")
}

// TestTcpServer 模拟服务端
func TestTcpServer(t *testing.T) {

	// 创建一个server句柄
	s := znet.NewServer()

	// 注册连接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 开启服务
	s.Serve()
}
