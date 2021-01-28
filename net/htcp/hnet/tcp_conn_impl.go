package hnet

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync"

	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
	"github.com/happylay-cloud/gf-extend/net/htcp/hutils"
)

type TcpConnection struct {
	// 当前Conn属于哪个Server
	TcpServer hiface.ITcpServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID int64
	// 消息管理MsgId和对应处理方法的消息管理模块
	PkgHandler hiface.ITcpPkgHandle
	// 告知该连接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	// 无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	// 有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	sync.RWMutex
	// 连接属性
	property map[string]interface{}
	// 保护当前property的锁
	propertyLock sync.Mutex
	// 当前连接的关闭状态
	isClosed bool
}

// 创建连接的方法
func NewTcpConnection(server hiface.ITcpServer, conn *net.TCPConn, connID int64, pkgHandler hiface.ITcpPkgHandle) *TcpConnection {
	// 初始化Conn属性
	c := &TcpConnection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		PkgHandler:  pkgHandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, hutils.GlobalHTcpObject.MaxPkgChanLen),
		property:    make(map[string]interface{}),
	}

	// 将新创建的Conn添加到连接管理中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// StartWriter 写消息协程，用户将数据发送给客户端
func (c *TcpConnection) StartWriter() {
	g.Log().Line(false).Info("[写协程正在运行]")

	defer g.Log().Line(false).Info(c.RemoteAddr().String(), "[写连接退出！]")

	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				g.Log().Line(false).Error("发送数据错误，", err, "写连接退出！")
				return
			}
			g.Log().Line().Debug("发送非缓冲数据成功！")

		case data, ok := <-c.msgBuffChan:
			if ok {
				// 有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					g.Log().Line(false).Error("发送缓冲数据错误：", err, " 写连接退出！")
					return
				}
				g.Log().Line().Debug("发送缓冲数据成功！")
			} else {
				g.Log().Line(false).Info("消息缓冲通道被关闭")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息协程，用于从客户端中读取数据
func (c *TcpConnection) StartReader() {
	g.Log().Line(false).Info("[读协程正在运行]")

	defer g.Log().Line(false).Info(c.RemoteAddr().String(), "[读连接退出！]")
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 创建拆包的对象
			dp := NewAcceptPackage()
			// 数据包拆包
			unpack, err := dp.Unpack(c.GetTCPConnection())
			if err != nil {
				g.Log().Line().Error("数据包拆包异常", err)
				return
			}
			// TODO 数据解码

			// 得到当前客户端请求的Request数据
			req := TcpRequest{
				conn: c,
				pkg:  unpack,
			}

			if hutils.GlobalHTcpObject.WorkerPoolSize > 0 {
				// 已经启动工作池机制，将消息交给Worker处理
				c.PkgHandler.SendPkgToTaskQueue(&req)
			} else {
				// 从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.PkgHandler.DoPkgHandler(&req)
			}
		}
	}
}

// SendTcpPkg 发送tcp数据包
//  @pkgBodyType 数据类型
//  @pkg         数据内容（自定义结构体对象）
//  @userBuf     是否启动缓冲（默认关闭）
func (c *TcpConnection) SendTcpPkg(pkgBodyType byte, pkg interface{}, userBuf ...bool) error {
	// 数据序列化
	bytes, err := json.Marshal(pkg)
	if err != nil {
		g.Log().Line(false).Error("数据包序列化异常", err)
		return err
	}

	c.RLock()
	if c.isClosed == true {
		c.RUnlock()
		return errors.New("发送消息时连接关闭")
	}
	c.RUnlock()

	// 构建数据包封包对象
	dp := NewPackage(pkgBodyType, bytes)

	// 数据封包
	msg, err := dp.Pack()
	if err != nil {
		g.Log().Line(false).Error("数据包封包失败，异常信息：", err)
		return err
	}

	if len(userBuf) > 0 && userBuf[0] {
		// 写回带缓冲客户端
		c.msgBuffChan <- msg
		return nil
	}

	// 写回非缓冲客户端
	c.msgChan <- msg
	return nil
}

// Start 启动连接，让当前连接开始工作
func (c *TcpConnection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	// 开启用户从客户端读取数据流程的协程
	go c.StartReader()
	// 2、开启用于写回客户端数据流程的协程
	go c.StartWriter()
	// 按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止连接，结束当前连接状态
func (c *TcpConnection) Stop() {
	g.Log().Line(false).Info("关闭连接 ConnID = ", c.ConnID)

	c.Lock()
	defer c.Unlock()

	// 如果用户注册了该连接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		g.Log().Line(false).Error("关闭连接异常，", err)
	}

	// 关闭Writer
	c.cancel()

	// 将连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭该连接全部管道
	close(c.msgBuffChan)
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *TcpConnection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *TcpConnection) GetConnID() int64 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *TcpConnection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SetProperty 设置连接属性
func (c *TcpConnection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *TcpConnection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("未发现连接属性")
	}
}

// RemoveProperty 移除连接属性
func (c *TcpConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
