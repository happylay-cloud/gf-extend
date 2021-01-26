package ztimer

import (
	"fmt"
	"testing"
)

func SayHello(message ...interface{}) {
	fmt.Println(message[0].(string), " ", message[1].(string))
}

// TestDelayfunc 针对delayFunc.go做单元测试，主要测试延迟函数结构体是否正常使用
func TestDelayfunc(t *testing.T) {
	df := NewDelayFunc(SayHello, []interface{}{"hello", "gf-plus"})
	fmt.Println("df.String() = ", df.String())
	df.Call()
}
