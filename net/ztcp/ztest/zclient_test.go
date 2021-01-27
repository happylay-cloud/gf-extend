package ztest

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/happylay-cloud/gf-extend/net/ztcp/znet"
)

// TestTcpClient 模拟客户端
func TestTcpClient(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:2021")
	if err != nil {
		fmt.Println("客户端启动错误，退出！")
		return
	}

	for {
		// 发封包message消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("客户端 MsgID=0，[Ping]")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("写错误，异常：", err)
			return
		}

		// 先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		// ReadFull会把msg填充满为止
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("读头信息错误")
			break
		}
		// 将headData字节流拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("服务器解包错误：", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// msg是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			// 根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("服务器解包数据错误：", err)
				return
			}

			fmt.Println("==> 测试路由：[Ping]接收消息：ID=", msg.Id, "，len=", msg.DataLen, "，data=", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
