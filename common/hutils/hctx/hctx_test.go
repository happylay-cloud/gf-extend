package hctx

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext1(t *testing.T) {
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("监控退出，停止了...")
				return
			default:
				fmt.Println("goroutine监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	stop <- true
	// 为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

//  使用close关闭通道后，如果该通道是无缓冲的，则它会从原来的阻塞变成非阻塞，也就是可读的，
//  只不过读到的会一直是零值，因此根据这个特性就可以判断拥有该通道的goroutine是否要关闭。
//
// DEPRECATED: 测试用例
func monitor1(ch chan bool, number int) {
	for {
		select {
		case v := <-ch:
			// 仅当ch通道被close，或者有数据发过来(无论是true还是false)才会走到这个分支
			fmt.Printf("监控器%v，接收到通道值为：%v，监控结束。\n", number, v)
			return
		default:
			fmt.Printf("监控器%v，正在监控中...\n", number)
			time.Sleep(2 * time.Second)
		}
	}
}

func TestContext2(t *testing.T) {

	stopSingal := make(chan bool)

	for i := 1; i <= 5; i++ {
		go monitor1(stopSingal, i)
	}

	time.Sleep(1 * time.Second)
	// 关闭所有goroutine
	close(stopSingal)

	// 等待5s，若此时屏幕没有输出<正在监控中>就说明所有的goroutine都已经关闭
	time.Sleep(5 * time.Second)

	fmt.Println("主程序退出！！")
}

func monitor2(ctx context.Context, number int) {
	for {
		select {
		// 写成case <- ctx.Done()
		// 仅是看到Done返回的内容
		case v := <-ctx.Done():
			fmt.Printf("监控器%v，接收到通道值为：%v，监控结束。\n", number, v)
			return
		default:
			fmt.Printf("监控器%v，正在监控中...\n", number)
			time.Sleep(2 * time.Second)
		}
	}
}

func TestContext3(t *testing.T) {
	// 以context.Background()为parent context定义一个可取消的context
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 5; i++ {
		go monitor2(ctx, i)
	}

	time.Sleep(1 * time.Second)
	// 关闭所有goroutine
	cancel()

	// 等待5s，若此时屏幕没有输出<正在监控中>就说明所有的goroutine都已经关闭
	time.Sleep(5 * time.Second)

	fmt.Println("主程序退出！！")
}

// WithDeadline 传入的第二个参数是 time.Time 类型，它是一个绝对的时间，
//  意思是在什么时间点超时取消。
func TestContext4(t *testing.T) {
	ctx01, cancel := context.WithCancel(context.Background())
	// 绝对的时间
	ctx02, cancel := context.WithDeadline(ctx01, time.Now().Add(1*time.Second))

	defer cancel()

	for i := 1; i <= 5; i++ {
		go monitor2(ctx02, i)
	}

	time.Sleep(5 * time.Second)
	if ctx02.Err() != nil {
		fmt.Println("监控器取消的原因: ", ctx02.Err())
	}

	fmt.Println("主程序退出！！")
}

//  WithTimeout 传入的第二个参数是 time.Duration 类型，
//   它是一个相对的时间，意思是多长时间后超时取消。
func TestContext5(t *testing.T) {
	ctx01, cancel := context.WithCancel(context.Background())
	// 相对的时间
	ctx02, cancel := context.WithTimeout(ctx01, 1*time.Second)

	defer cancel()

	for i := 1; i <= 5; i++ {
		go monitor2(ctx02, i)
	}

	time.Sleep(5 * time.Second)
	if ctx02.Err() != nil {
		fmt.Println("监控器取消的原因: ", ctx02.Err())
	}

	fmt.Println("主程序退出！！")
}

func monitor3(ctx context.Context, number int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("监控器%v，监控结束。\n", number)
			return
		default:
			// 获取 item 的值
			value := ctx.Value("item")
			fmt.Printf("监控器%v，正在监控 %v \n", number, value)
			time.Sleep(2 * time.Second)
		}
	}
}

// 以ctx02为父context，再创建一个能携带value的ctx03，
//  由于他的父context是ctx02，所以ctx03也具备超时自动取消的功能。
func TestContext6(t *testing.T) {
	ctx01, cancel := context.WithCancel(context.Background())
	ctx02, cancel := context.WithTimeout(ctx01, 1*time.Second)
	ctx03 := context.WithValue(ctx02, "item", "CPU")

	defer cancel()
	for i := 1; i <= 5; i++ {
		go monitor3(ctx03, i)
	}

	time.Sleep(5 * time.Second)
	if ctx02.Err() != nil {
		fmt.Println("监控器取消的原因: ", ctx02.Err())
	}

	fmt.Println("主程序退出！！")
}
