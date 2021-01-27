package hnet

import (
	"fmt"
	"net"

	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
	"github.com/happylay-cloud/gf-extend/net/htcp/hutils"
	"github.com/happylay-cloud/gf-extend/net/ztcp/zutils"
)

var gfPlusLogo = `                                        
 ██████╗ ███████╗    ██████╗ ██╗     ██╗   ██╗███████╗      ██╗  ██╗████████╗ ██████╗██████╗ 
██╔════╝ ██╔════╝    ██╔══██╗██║     ██║   ██║██╔════╝      ██║  ██║╚══██╔══╝██╔════╝██╔══██╗
██║  ███╗█████╗█████╗██████╔╝██║     ██║   ██║███████╗█████╗███████║   ██║   ██║     ██████╔╝
██║   ██║██╔══╝╚════╝██╔═══╝ ██║     ██║   ██║╚════██║╚════╝██╔══██║   ██║   ██║     ██╔═══╝ 
╚██████╔╝██║         ██║     ███████╗╚██████╔╝███████║      ██║  ██║   ██║   ╚██████╗██║     
 ╚═════╝ ╚═╝         ╚═╝     ╚══════╝ ╚═════╝ ╚══════╝      ╚═╝  ╚═╝   ╚═╝    ╚═════╝╚═╝`
var topLine = `┌───────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└───────────────────────────────────────────────────┘`

// TcpServer ITcpServer接口实现，定义一个TcpServer服务类
type TcpServer struct {
	Name        string                           // 服务器的名称
	IPVersion   string                           // tcp4或other
	IP          string                           // 服务绑定的IP地址
	Port        int                              // 服务绑定的端口
	pkgHandler  hiface.ITcpPkgHandle             // 当前服务的消息管理模块，用来绑定路由和对应的处理器
	ConnMgr     hiface.ITcpConnManager           // 当前服务的连接管理器
	OnConnStart func(conn hiface.ITcpConnection) // 服务创建连接时Hook函数
	OnConnStop  func(conn hiface.ITcpConnection) // 服务断开连接时的Hook函数
}

// NewTcpServer 创建一个服务器句柄
func NewTcpServer() hiface.ITcpServer {
	printLogo()

	s := &TcpServer{
		Name:       zutils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         zutils.GlobalObject.Host,
		Port:       zutils.GlobalObject.TcpPort,
		pkgHandler: NewPkgHandle(),
		ConnMgr:    NewTcpConnManager(),
	}
	return s
}

// Start 开启网络服务
func (s *TcpServer) Start() {
	g.Log().Line(false).Info("[启动] TCP服务名称：", s.Name)

	// 开启一个go去做服务端Listener业务
	go func() {
		// 启动worker工作池机制
		s.pkgHandler.StartWorkerPool()
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			g.Log().Line(false).Error("解析tcp地址错误：", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			g.Log().Line(false).Error("监听：", s.IPVersion, "错误：", err)
			return
		}

		// 已经监听成功
		g.Log().Line(false).Info("[运行] 成功开启TCP服务，监听地址：", s.IP, ":", s.Port)

		// TODO 缺失一个生成连接会话id功能
		var cid int64
		cid = 0

		// 启动server网络连接业务
		for {
			// 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				g.Log().Line(false).Error("接收客户端连接错误：", err)
				continue
			}
			g.Log().Line(false).Debug("获取连接远程地址 addr =", conn.RemoteAddr().String())

			// 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= hutils.GlobalHTcpObject.MaxConn {
				err := conn.Close()
				if err != nil {
					g.Log().Line(false).Error("关闭连接异常：", err)
				}
				continue
			}

			// 处理该新连接请求的业务方法，此时应该有handler和conn是绑定的
			dealConn := NewTcpConnection(s, conn, cid, s.pkgHandler)
			cid++

			// 3.4、启动当前连接的处理业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务
func (s *TcpServer) Stop() {
	g.Log().Line(false).Info("[停止] TCP服务，服务名称：", s.Name)

	// 将其他需要清理的连接信息或者其他信息，也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

// Serve 运行服务
func (s *TcpServer) Serve() {
	s.Start()

	// TODO Server.Serve()是否在启动服务的时候，还要处理其他的事情，可以在这里添加

	// 阻塞，否则主Go退出，监听器listener的go将会退出
	select {}
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端连接处理使用
func (s *TcpServer) AddRouter(handlerRouter string, router hiface.ITcpRouter) {
	s.pkgHandler.AddRouter(handlerRouter, router)
}

// GetConnMgr 获取连接管理
func (s *TcpServer) GetConnMgr() hiface.ITcpConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *TcpServer) SetOnConnStart(hookFunc func(hiface.ITcpConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *TcpServer) SetOnConnStop(hookFunc func(hiface.ITcpConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *TcpServer) CallOnConnStart(conn hiface.ITcpConnection) {
	if s.OnConnStart != nil {
		g.Log().Line().Debug("[创建新连接] 连接创建时调用...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *TcpServer) CallOnConnStop(conn hiface.ITcpConnection) {
	if s.OnConnStop != nil {
		g.Log().Line().Debug("[关闭旧连接] 连接关闭时调用...")
		s.OnConnStop(conn)
	}
}

func printLogo() {
	fmt.Println(gfPlusLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s 1、只做增强不做改变                                  %s", borderLine, borderLine))
	fmt.Println(fmt.Sprintf("%s 2、简化开发、提高效率而生                             %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Gf-Plus] 版本：%s 最大连接：%d 可接受最大数据包：%d\n",
		zutils.GlobalObject.Version,
		zutils.GlobalObject.MaxConn,
		zutils.GlobalObject.MaxPacketSize)
}

func init() {
}
