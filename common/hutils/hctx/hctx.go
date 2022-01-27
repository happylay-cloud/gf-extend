package hctx

import (
	"context"
	"errors"
	"github.com/gogf/gf/frame/g"
	"sync"
	"time"
)

// DoManyTask 多任务结构体
type DoManyTask struct {
	Count      int
	ChannelObj chan interface{}
	Timeout    time.Duration
	Debug      bool
}

// DoTaskSuccessOne 并发执行任务，成功一次即返回
func (do *DoManyTask) DoTaskSuccessOne(params interface{}, doBizFun func(do *DoManyTask, ctx context.Context, wg *sync.WaitGroup, index int, params interface{})) (interface{}, error) {

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), do.Timeout)

	// 取消上下文
	defer cancel()

	// 创建同步等待组
	wg := &sync.WaitGroup{}
	wg.Add(do.Count)

	// 自动关闭通道
	go do.deferCloseTaskChannel(do.ChannelObj, wg)

	// 执行多个协程
	for index := 0; index < do.Count; index++ {
		go doBizFun(do, ctx, wg, index, params)
	}

	// 获取任务数据
	data, open := <-do.ChannelObj
	if !open {
		return nil, errors.New("任务执行失败")
	}

	return data, nil
}

// deferCloseTaskChannel 关闭通道，内部方法
func (do *DoManyTask) deferCloseTaskChannel(channel chan interface{}, wg *sync.WaitGroup) {
	// 关闭通道
	defer close(channel)
	g.Log().Line(false).Info("..........................监控channel通道...............................")
	// 子线程等待
	wg.Wait()
	g.Log().Line(false).Info("..........................关闭channel通道...............................")
}

// WaitDataReturn 获取返回数据，自定义回调函数中必须执行此方法，以获取返回值
func (do *DoManyTask) WaitDataReturn(index int, ctx context.Context, data interface{}) {
	select {
	case <-ctx.Done(): // 取消执行
		if do.Debug {
			g.Log().Line(false).Debug("关闭任务，任务序号：", index)
		}
		break
	case do.ChannelObj <- data: // 传递数据
		if do.Debug {
			g.Log().Line(false).Debug("任务执行成功，任务序号：", index)
		}
	}
}
