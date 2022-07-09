package hid

import (
	"github.com/gogf/gf/frame/g"
	"sync"
)

var once sync.Once

// SnowWorker 全局唯一雪花算法实例对象
var snowWorker *SnowFlake

// GetSnowWorker 获取雪花算法实例，警告：golang版本的雪花算法生成的ID有序递增
func GetSnowWorker() *SnowFlake {
	once.Do(func() {
		g.Log().Line(false).Info("初始化雪花算法")
		snowFlake, err := NewSnowFlake(0, 0)
		if err != nil {
			g.Log().Line(false).Error("初始化雪花算法失败", err)
		}
		snowWorker = snowFlake
	})

	return snowWorker
}
