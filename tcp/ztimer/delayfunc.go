package ztimer

import (
	"fmt"
	"reflect"

	"github.com/happylay-cloud/gf-extend/tcp/zlog"
)

// DelayFunc 定义一个延迟调用函数
//  延迟调用函数就是 时间定时器超时的时候，触发的事先注册好的回调函数
type DelayFunc struct {
	f    func(...interface{}) // f : 延迟函数调用原型
	args []interface{}        // args: 延迟调用函数传递的形参
}

// NewDelayFunc 创建一个延迟调用函数
func NewDelayFunc(f func(v ...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:    f,
		args: args,
	}
}

// String 打印当前延迟函数的信息，用于日志记录
func (df *DelayFunc) String() string {
	return fmt.Sprintf("{延迟函数：%s，参数：%v}", reflect.TypeOf(df.f).Name(), df.args)
}

// Call 执行延迟函数，如果执行失败，抛出异常
func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Error(df.String(), "调用错误：", err)
		}
	}()

	// 调用定时器超时函数
	df.f(df.args...)
}
