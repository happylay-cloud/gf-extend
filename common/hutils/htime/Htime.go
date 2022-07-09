package htime

import "time"

// Sleep 程序休眠
func Sleep(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}

// TimestampFormat 时间戳（毫秒）格式化
func TimestampFormat(timestamp int64) string {
	return time.Unix(timestamp/1e3, 0).Format("2006-01-02 15:04:05")
}
