package hctx

import (
	"context"
	"errors"
	"github.com/gogf/gf/frame/g"
	"sync"
	"time"
)

// DoTaskSuccessOne 并发执行任务，成功一次即返回
func DoTaskSuccessOne(count int, channelObj chan interface{}, timeout time.Duration, doBizFun func(ctx context.Context, wg *sync.WaitGroup, index int, params ...interface{})) (interface{}, error) {

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// 取消上下文
	defer cancel()

	// 创建同步等待组
	wg := &sync.WaitGroup{}
	wg.Add(count)

	// 自动关闭通道
	go deferCloseTaskChannel(channelObj, wg)

	// 执行多个协程
	for x := 0; x <= count; x++ {
		go doBizFun(ctx, wg, x, channelObj)
	}

	// 获取任务数据
	data, open := <-channelObj
	if !open {
		return nil, errors.New("任务执行失败")
	}

	return data, nil
}

// deferCloseTaskChannel 关闭通道，内部方法
func deferCloseTaskChannel(channel chan interface{}, wg *sync.WaitGroup) {
	// 关闭通道
	defer close(channel)
	g.Log().Line(false).Info("..........................监控channel通道...............................")
	// 子线程等待
	wg.Wait()
	g.Log().Line(false).Info("..........................关闭channel通道...............................")
}
