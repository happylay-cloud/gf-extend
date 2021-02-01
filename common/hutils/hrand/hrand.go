package hrand

import (
	"math/rand"
	"time"
)

// GetRandInt 获取随机数，返回值区间[0,num)
func GetRandInt(num int) int {
	// 设置种子数（时间戳）
	rand.Seed(time.Now().UnixNano())
	// 生成随机数
	return rand.Intn(num)
}

// GetRandIntAndRange 获取指定范围随机数，返回值区间[min,max)
func GetRandIntAndRange(min, max int) int {
	// 设置种子数（时间戳）
	rand.Seed(time.Now().UnixNano())
	// 生成随机数
	return rand.Intn(max-min) + min
}
