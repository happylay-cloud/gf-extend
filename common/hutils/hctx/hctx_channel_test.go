package hctx

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/google/uuid"

	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannelCtx(t *testing.T) {

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 上下文传值
	ctx = context.WithValue(ctx, "desc", "关闭子协程")

	// 取消上下文
	defer cancel()

	// 创建验证码通道
	doorCode := make(chan interface{})

	// 设置计数器
	count := 300

	// 创建同步等待组
	wg := &sync.WaitGroup{}
	wg.Add(count)

	go deferCloseChannel(doorCode, wg)

	// 执行多个协程
	for x := 1; x <= count; x++ {
		go doSomething(x, doorCode, ctx, wg)
	}

	// 读取数据
	data, open := <-doorCode
	if !open {
		fmt.Println("通道已关闭，任务执行失败！！！！！！！！！！！！！！！！！！")
		return
	}

	fmt.Println("读取业务返回数据：", data)

}

func deferCloseChannel(doorCode chan interface{}, wg *sync.WaitGroup) {
	// 关闭通道
	defer close(doorCode)

	fmt.Println("..........................监控channel通道...............................")
	// 子线程等待
	wg.Wait()
	fmt.Println("..........................关闭channel通道...............................")
}

func doSomething(x int, channel chan interface{}, ctx context.Context, wg *sync.WaitGroup) {
	fmt.Println("任务执行中...，序号：", x)
	fmt.Println("上下文传值：", ctx.Value("desc"))
	// 计数器减一
	defer wg.Done()
	// 查询数据
	data := uuid.New().String()
	// 处理业务数据
	if gstr.LenRune(data) > 0 {
		select {
		case <-ctx.Done(): // 取消执行
			fmt.Println("关闭任务，序号", x)
			break
		case channel <- data: // 传递数据
			fmt.Println("***********************任务执行成功，序号", x)
		}
	}
}

func TestTaskContext(t *testing.T) {

	// 定义任务
	doManyTask := DoManyTask{
		Count:      300,
		ChannelObj: make(chan interface{}),
		Timeout:    20 * time.Second,
		Debug:      true,
	}

	// 执行任务
	successOne, err := doManyTask.DoTaskSuccessOne(nil, func(do *DoManyTask, ctx context.Context, wg *sync.WaitGroup, index int, params interface{}) {
		fmt.Println("任务执行中...，序号：", index)
		// 计数器减一
		defer wg.Done()
		// 业务处理
		data := uuid.New().String()
		// 返回数据
		if gstr.LenRune(data) > 0 {
			// 获取返回结果
			do.WaitDataReturn(index, ctx, data)
		}
	})

	if err != nil {
		fmt.Println("任务执行失败：", err)
		return
	}

	// 获取返回值
	g.Dump("任务返回值：", successOne)

}
