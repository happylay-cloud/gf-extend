package zlog

import (
	"testing"
)

func TestStdZLog(t *testing.T) {

	// 测试默认debug输出
	Debug("debug 内容1")
	Debug("debug 内容2")

	Debugf("debug a = %d\n", 10)

	// 设置log标记位，加上长文件名称 和 微秒 标记
	ResetFlags(BitDate | BitLongFile | BitLevel)
	Info("info 内容")

	// 设置日志前缀，主要标记当前日志模块
	SetPrefix("MODULE")
	Error("error 内容")

	// 添加标记位
	AddFlag(BitShortFile | BitTime)
	Stack("Stack！")

	// 设置日志写入文件
	SetLogFile("./log", "testfile.log")
	Debug("===> debug 内容 ~~666")
	Debug("===> debug 内容 ~~888")
	Error("===> Error ！！！~~~555~~~")

	// 关闭debug调试
	CloseDebug()
	Debug("===> 我不应该出现~！")
	Debug("===> 我不应该出现~！")
	Error("===> Error debug后关闭！！！")

}
