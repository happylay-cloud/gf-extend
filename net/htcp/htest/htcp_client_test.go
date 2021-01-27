package htest

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hnet"
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
		dp := hnet.NewPackage(0, []byte("哈哈哈哈"))
		pack, _ := dp.Pack()
		_, err := conn.Write(pack)
		if err != nil {
			fmt.Println("写错误，异常：", err)
			return
		}

		unpack, err := hnet.NewAcceptPackage().Unpack(conn)
		g.Dump(unpack)

		time.Sleep(1 * time.Second)
	}
}
