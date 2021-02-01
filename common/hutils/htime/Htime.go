package htime

import "time"

// Sleep 程序休眠
func Sleep(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}
