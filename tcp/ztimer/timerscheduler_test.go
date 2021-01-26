package ztimer

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/happylay-cloud/gf-extend/tcp/zlog"
)

// 触发函数
func foo(args ...interface{}) {
	fmt.Printf("我是No. %d函数，延迟%dms\n", args[0].(int), args[1].(int))
}

// 手动创建调度器运转时间轮
func TestNewTimerScheduler(t *testing.T) {
	timerScheduler := NewTimerScheduler()
	timerScheduler.Start()

	// 在scheduler中添加timer
	for i := 1; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tid, err := timerScheduler.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			zlog.Error("创建timer错误", tid, err)
			break
		}
	}

	// 执行调度器触发函数
	go func() {
		delayFuncChan := timerScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			df.Call()
		}
	}()

	// 阻塞等待
	select {}
}

// 采用自动调度器运转时间轮
func TestNewAutoExecTimerScheduler(t *testing.T) {
	autoTS := NewAutoExecTimerScheduler()

	// 给调度器添加Timer
	for i := 0; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tid, err := autoTS.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			zlog.Error("创建timer错误", tid, err)
			break
		}
	}

	// 阻塞等待
	select {}
}

// 时间轮定时器调度器单元测试
// TestCancelTimerScheduler 测试取消一个定时器
func TestCancelTimerScheduler(t *testing.T) {
	Scheduler := NewAutoExecTimerScheduler()
	f1 := NewDelayFunc(foo, []interface{}{3, 3})
	f2 := NewDelayFunc(foo, []interface{}{5, 5})
	timerId1, _ := Scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)
	timerId2, _ := Scheduler.CreateTimerAfter(f2, time.Duration(5)*time.Second)
	log.Printf("timerId1=%d，timerId2=%d\n", timerId1, timerId2)
	// 删除timerId1
	Scheduler.CancelTimer(timerId1)

	// 阻塞等待
	select {}
}
