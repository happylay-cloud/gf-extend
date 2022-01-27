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
	Count   int           // 任务数
	Timeout time.Duration // 超时时间
	Debug   bool          // 开启调试模式
}

// DoTaskSuccessOne 并发执行任务，成功一次即返回
func (do *DoManyTask) DoTaskSuccessOne(params interface{}, doBizFun func(do *DoManyTask, ctx context.Context, channel chan interface{}, wg *sync.WaitGroup, index int, params interface{})) (interface{}, error) {

	// 创建通道
	channelObj := make(chan interface{})

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), do.Timeout)

	// 取消上下文
	defer cancel()

	// 创建同步等待组
	wg := &sync.WaitGroup{}
	wg.Add(do.Count)

	// 自动关闭通道
	go do.deferCloseTaskChannel(channelObj, wg)

	// 执行多个协程
	for index := 0; index < do.Count; index++ {
		go doBizFun(do, ctx, channelObj, wg, index, params)
	}

	// 获取任务数据
	data, open := <-channelObj
	if !open {
		g.Log().Line(false).Error("任务执行失败")
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
func (do *DoManyTask) WaitDataReturn(isNeedReturn bool, ctx context.Context, channel chan interface{}, wg *sync.WaitGroup, index int, data interface{}) {
	// 计数器减一
	defer wg.Done()

	// 处理返回值
	if isNeedReturn {
		select {
		case <-ctx.Done(): // 取消执行
			if do.Debug {
				g.Log().Line(false).Debug("关闭任务，任务序号：", index)
			}
			break
		case channel <- data: // 传递数据
			if do.Debug {
				g.Log().Line(false).Debug("任务执行成功，任务序号：", index)
			}
		}
	}

}
