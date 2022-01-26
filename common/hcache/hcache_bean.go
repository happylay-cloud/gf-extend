package hcache

import (
	"github.com/xujiajun/nutsdb"
	"sync"
)

// NutsDbCache 单例缓存对象
var NutsDbCache *nutsdb.DB

// 只操作一次对象
var once sync.Once

// NutsDbCacheErr 单例错误对象
var NutsDbCacheErr error

// NewNutsDbCacheBean 实例化单例缓存对象
func NewNutsDbCacheBean() (*nutsdb.DB, error) {

	once.Do(func() {
		// 默认配置
		opt := nutsdb.DefaultOptions
		// 自动创建数据库
		opt.Dir = "./.cache/nutsdb"
		// 开启数据库
		db, err := nutsdb.Open(opt)
		if err != nil {
			NutsDbCacheErr = err
		}

		NutsDbCache = db
	})

	return NutsDbCache, nil
}
