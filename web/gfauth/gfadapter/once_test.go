package gfadapter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// once只执行一次，无论是否更换once.Do()方法，这个sync.Once块只会执行一次
var testOnce sync.Once

// TestOnce 测试Once对象
func TestOnce(t *testing.T) {

	for i, v := range make([]string, 1) {
		testOnce.Do(tomato)
		fmt.Println("计数:", v, "-", i)
	}
	for i := 0; i < 1; i++ {
		go func() {
			testOnce.Do(banana)
			fmt.Println("异步函数执行完毕")
		}()
	}
	time.Sleep(1000)

	//Output:
	//🍅
	//计数:  - 0
	//异步函数执行完毕
}

func tomato() {
	fmt.Println("🍅")
}
func banana() {
	fmt.Println("🍌")
}
