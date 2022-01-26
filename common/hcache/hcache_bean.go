package hcache

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/xujiajun/nutsdb"

	"errors"
	"os"
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

		// 解析安装真正路径
		var realPath string

		if gstr.HasPrefix(nutsDbPath, "./") {
			pwd, err := os.Getwd()
			if err != nil {
				nutsDbCacheErr = err
			}
			realPath = pwd + gstr.TrimLeftStr(nutsDbPath, ".")
		} else {
			realPath = nutsDbPath
		}

		g.Log().Line(false).Info("完成NutsDb缓存数据库实例化，当前安装路径：" + realPath)
	})

	return nutsDbCache, nutsDbCacheErr
}

// InitSetNutsDbGlobalPath 初始化NutsDb数据库安装路径，仅初始化一次，默认安装位置：/tmp/hcache/nutsdb/.cache/db
func InitSetNutsDbGlobalPath(dir string) error {
	if nutsDbCache != nil {
		g.Log().Line(false).Error("初始化NutsDb数据库安装路径失败，已初始化位置：" + nutsDbPath)
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
