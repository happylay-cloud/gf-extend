package hnet

import (
	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
	"github.com/happylay-cloud/gf-extend/net/htcp/hutils"
)

type PkgHandle struct {
	Apis           map[string]hiface.ITcpRouter // 存放每个HandlerRouter所对应的处理器
	WorkerPoolSize int64                        // 业务工作Worker池的数量
	TaskQueue      []chan hiface.ITcpRequest    // Worker负责取任务的消息队列，一个Worker对应一个Queue
}

func NewPkgHandle() *PkgHandle {
	return &PkgHandle{
		Apis:           make(map[string]hiface.ITcpRouter),
		WorkerPoolSize: hutils.GlobalHTcpObject.WorkerPoolSize,
		TaskQueue:      make([]chan hiface.ITcpRequest, hutils.GlobalHTcpObject.WorkerPoolSize),
	}
}

// SendPkgToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (p *PkgHandle) SendPkgToTaskQueue(request hiface.ITcpRequest) {
	// 根据ConnID来分配当前的连接应该由哪个worker负责处理，
	// 轮询的平均分配法则，得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % p.WorkerPoolSize
	// 记录日志
	g.Log().Line(false).Debug(
		"添加ConnID：", request.GetConnection().GetConnID(),
		"请求HandlerRouter：", request.GetHandlerRouter(),
		"分配到workerID：", workerID)

	// 将请求消息发送给任务队列
	p.TaskQueue[workerID] <- request
}

// DoPkgHandler 以非阻塞方式处理消息
func (p *PkgHandle) DoPkgHandler(request hiface.ITcpRequest) {
	handler, ok := p.Apis[request.GetHandlerRouter()]
	if !ok {
		// 记录错误日志
		g.Log().Line(false).Error("路由HandlerRouter：", request.GetHandlerRouter(), "没有发现！")
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (p *PkgHandle) AddRouter(handlerRouter string, router hiface.ITcpRouter) {
	// 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := p.Apis[handlerRouter]; ok {
		panic("路由重复handlerRouter：" + handlerRouter)
	}
	// 添加msg与api的绑定关系
	p.Apis[handlerRouter] = router
	// 记录日志
	g.Log().Line(false).Info("[Gf-Plus] 添加接口路由 handlerRouter = ", handlerRouter)
}

// StartOneWorker 启动一个Worker工作协程
func (p *PkgHandle) StartOneWorker(workerID int, taskQueue chan hiface.ITcpRequest) {
	// 记录日志
	g.Log().Line(false).Info("[启动] 工作协程 workerID = ", workerID)
	// 不断的等待队列中的消息
	for {
		select {
		// 有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			p.DoPkgHandler(request)
		}
	}
}

// StartWorkerPool 启动worker工作池
func (p *PkgHandle) StartWorkerPool() {
	// 遍历需要启动worker的数量，依此启动
	for i := 0; i < int(p.WorkerPoolSize); i++ {
		// 一个worker被启动，给当前worker对应的任务队列开辟空间
		p.TaskQueue[i] = make(chan hiface.ITcpRequest, hutils.GlobalHTcpObject.MaxWorkerTaskLen)
		// 启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go p.StartOneWorker(i, p.TaskQueue[i])
	}
}
