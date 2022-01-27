package hjsoup

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/google/uuid"

	"fmt"
	"sync"
	"testing"
)

func TestHttpClient(t *testing.T) {

	productCode, err := SearchByProductCode("6921168509256", false)
	if err != nil {
		fmt.Println(err)
	}
	g.Dump(productCode)
}

func TestSearchByProductCodeCache(t *testing.T) {
	// 商品条码
	productCode := "6928804011142"

	productCodeInfo, err := SearchByProductCodeCache(productCode)
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(productCodeInfo)
}

func TestChannel(t *testing.T) {

	// 创建验证码通道
	doorCode := make(chan string)

	// 设置计数器
	count := 300

	// 创建同步等待组
	wg := &sync.WaitGroup{}
	wg.Add(count)

	go deferCloseChannel(doorCode, wg)

	// 执行多个协程
	for x := 1; x <= count; x++ {
		go doSomething(x, doorCode, wg)
	}

	for data := range doorCode {
		fmt.Println("读取业务返回数据：", data)
	}

}

func doSomething(x int, doorCode chan string, wg *sync.WaitGroup) {
	// 查询数据
	data := uuid.New().String()
	// 处理业务数据
	if gstr.LenRune(data) > 0 {
		doorCode <- data
	}
	fmt.Println("执行协程业务：", x)
	// 计数器减一
	wg.Done()
}

func deferCloseChannel(doorCode chan string, wg *sync.WaitGroup) {
	// 子线程等待
	wg.Wait()
	// 关闭通道
	close(doorCode)
	fmt.Println("关闭channel通道")
}
