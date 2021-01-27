package znet

import (
	"fmt"
	"net"

	"github.com/happylay-cloud/gf-extend/net/ztcp/ziface"
	"github.com/happylay-cloud/gf-extend/net/ztcp/zutils"
)

var gfPlusLogo = `                                        
 ██████╗ ███████╗    ██████╗ ██╗     ██╗   ██╗███████╗
██╔════╝ ██╔════╝    ██╔══██╗██║     ██║   ██║██╔════╝
██║  ███╗█████╗█████╗██████╔╝██║     ██║   ██║███████╗
██║   ██║██╔══╝╚════╝██╔═══╝ ██║     ██║   ██║╚════██║
╚██████╔╝██║         ██║     ███████╗╚██████╔╝███████║
 ╚═════╝ ╚═╝         ╚═╝     ╚══════╝ ╚═════╝ ╚══════╝`
var topLine = `┌───────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└───────────────────────────────────────────────────┘`

// Server iServer接口实现，定义一个Server服务类
type Server struct {
	// 服务器的名称
	Name string
	// tcp4或other
	IPVersion string
	// 服务绑定的IP地址
	IP string
	// 服务绑定的端口
	Port int
	// 当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler ziface.IMsgHandle
	// 当前Server的连接管理器
	ConnMgr ziface.IConnManager
	// 该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	printLogo()

	s := &Server{
		Name:       zutils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         zutils.GlobalObject.Host,
		Port:       zutils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// ============== 实现ziface.IServer里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[启动] TCP服务名称：%s\n", s.Name)

	// 开启一个go去做服务端Listener业务
	go func() {
		// 0、启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		// 1、获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("解析tcp地址错误：", err)
			return
		}

		// 2、监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("监听：", s.IPVersion, "错误：", err)
			return
		}

		// 已经监听成功
		fmt.Printf("[运行] 成功开启TCP服务，监听地址：%s:%d\n", s.IP, s.Port)

		// TODO server.go应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		// 3、启动server网络连接业务
		for {
			// 3.1、阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("接收客户端连接错误：", err)
				continue
			}
			fmt.Println("获取连接远程地址 addr =", conn.RemoteAddr().String())

			// 3.2、设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= zutils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			// 3.3、处理该新连接请求的业务方法，此时应该有handler和conn是绑定的
			dealConn := NewConntion(s, conn, cid, s.msgHandler)
			cid++

			// 3.4、启动当前连接的处理业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[停止] TCP服务，服务名称：", s.Name)

	// 将其他需要清理的连接信息或者其他信息，也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	// TODO Server.Serve()是否在启动服务的时候，还要处理其他的事情，可以在这里添加

	// 阻塞，否则主Go退出，监听器listener的go将会退出
	select {}
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端连接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

// GetConnMgr 得到连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("[创建新连接] 连接创建时调用...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("[关闭旧连接] 连接关闭时调用...")
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
