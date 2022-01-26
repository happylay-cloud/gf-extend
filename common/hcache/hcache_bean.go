package hcache

import (
	"errors"
	"github.com/xujiajun/nutsdb"
	"sync"
)

// nutsDbCache 单例缓存对象
var nutsDbCache *nutsdb.DB

// nutsDbCacheErr 单例错误对象
var nutsDbCacheErr error

// once 只操作一次对象
var once sync.Once

// nutsDbPath
var nutsDbPath = "/tmp/hcache/nutsdb/.cache/db"

// GetNutsDbCacheBean 获取实例化单例缓存对象
func GetNutsDbCacheBean() (*nutsdb.DB, error) {

	once.Do(func() {
		// 默认配置
		opt := nutsdb.DefaultOptions
		// 自动创建数据库
		opt.Dir = GetNutsDbGlobalPath()
		// 开启数据库
		db, err := nutsdb.Open(opt)
		if err != nil {
			nutsDbCacheErr = err
		}
		nutsDbCache = db
	})

	return nutsDbCache, nutsDbCacheErr
}

// InitSetNutsDbGlobalPath 初始化NutsDb数据库安装路径，仅初始化一次，默认安装位置：/tmp/hcache/nutsdb/.cache/db
func InitSetNutsDbGlobalPath(dir string) error {
	if nutsDbCache != nil {
		return errors.New("初始化NutsDb数据库安装路径失败，已初始化位置：" + nutsDbPath)
	}
	// 更改默认安装路径
	nutsDbPath = dir
	return nil
}

// GetNutsDbGlobalPath 获取NutsDb数据库安装路径
func GetNutsDbGlobalPath() string {
	return nutsDbPath
}
