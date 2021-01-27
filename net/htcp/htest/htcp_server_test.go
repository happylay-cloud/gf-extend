package htest

import (
	"fmt"
	"testing"

	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
	"github.com/happylay-cloud/gf-extend/net/htcp/hnet"
)

// 创建连接的时候执行
func DoConnectionBegin(conn hiface.ITcpConnection) {
	fmt.Println("创建连接被执行... ")

	// 设置两个连接属性，在连接创建之后
	fmt.Println("设置连接属性")
	conn.SetProperty("name", "gf-puls")

	err := conn.SendTcpPkg(2, []byte("开始创建连接..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn hiface.ITcpConnection) {
	// 在连接销毁之前，查询conn的name属性
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("连接属性 name = ", name)
	}

	fmt.Println("连接断开被执行...")
}

// TestTcpServer 模拟服务端
func TestTcpServer(t *testing.T) {

	// 创建一个server句柄
	s := hnet.NewTcpServer()

	// 注册连接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 配置路由
	s.AddRouter("1", &PingRouter{})
	s.AddRouter("2", &HelloRouter{})

	// 开启服务
	s.Serve()
}
