package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/happylay-cloud/gf-extend/net/tcp/ziface"
)

// ClientTest 模拟客户端
//  go test -v ./znet -run=TestServer
func ClientTest(i uint32) {

	fmt.Println("启动客户端测试...")
	// 3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:2021")
	if err != nil {
		fmt.Println("启动客户端错误，退出！")
		return
	}

	for {
		dp := NewDataPack()
		msg, _ := dp.Pack(NewMsgPackage(i, []byte("客户端测试消息")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("客户端写消息错误：", err)
			return
		}

		// 先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("客户端读head消息错误：", err)
			return
		}

		// 将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("客户端head拆包错误：", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// msg是有data数据的，需要再次读取data数据
			msg := msgHead.(*Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// 根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("客户端data拆包错误")
				return
			}

			fmt.Printf("==> 客户端接收消息：Id = %d, len = %d , data = %s\n", msg.Id, msg.DataLen, msg.Data)
		}

		time.Sleep(time.Second)
	}
}

// PingRouter 模拟服务器端
//  测试自定义路由
type PingRouter struct {
	BaseRouter
}

// PreHandle 测试PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("调用路由PreHandle方法")
	err := request.GetConnection().SendMsg(1, []byte("ping之前...\n"))
	if err != nil {
		fmt.Println("PreHandle发送消息错误：", err)
	}
}

// Handle 测试Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("调用PingRouter处理器")
	// 先读取客户端的数据，再回写ping
	fmt.Println("接收客户消息：msgId=", request.GetMsgID(), "，data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping\n"))
	if err != nil {
		fmt.Println("处理发送消息错误：", err)
	}
}

// PostHandle 测试PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("调用路由PostHandle")
	err := request.GetConnection().SendMsg(1, []byte("ping之后...\n"))
	if err != nil {
		fmt.Println("PostHandle发送消息错误：", err)
	}
}

type HelloRouter struct {
	BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("调用helloRouter处理器")
	fmt.Printf("接收客户端消息：msgId=%d，data=%s\n", request.GetMsgID(), string(request.GetData()))

	err := request.GetConnection().SendMsg(2, []byte("gf-plus"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("创建连接时被调用...")
	err := conn.SendMsg(2, []byte("创建连接"))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("断开连接时被调用")
}

func TestServer(t *testing.T) {
	// 创建一个server句柄
	s := NewServer()

	// 注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 多路由
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 客户端测试
	go ClientTest(1)
	go ClientTest(2)

	// 开启服务
	go s.Serve()

	select {
	case <-time.After(time.Second * 10):
		return
	}
}
