package hcache

import (
	"errors"
	"os"
	"sync"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/xujiajun/nutsdb"
)

// nutsDbCache 单例缓存对象
var nutsDbCache *nutsdb.DB

// nutsDbCacheErr 单例错误对象
var nutsDbCacheErr error

// once 只操作一次对象
var once sync.Once

// defaultNutsDbPath 默认安装路径
var defaultNutsDbPath = "/tmp/hcache/nutsdb/.cache/db"

// nutsDbPath 真正安装路径
var nutsDbPath string

// GetNutsDbCacheBean 获取实例化单例缓存对象
func GetNutsDbCacheBean() (*nutsdb.DB, error) {

	once.Do(func() {

		if gstr.LenRune(nutsDbPath) == 0 {
			// 启用默认安装路径
			nutsDbPath = defaultNutsDbPath
			g.Log().Line(false).Warning("启用NutsDb缓存数据库默认安装路径，安装路径：" + nutsDbPath +
				"，警告：默认的NutsDb安装路径，在同时启动多个项目下可能会产生缓存共享，建议使用自定义安装路径，调用方法：hcache.InitSetNutsDbGlobalPath()。")
		} else {
			// 解析安装真正安装路径
			var realPath string
			// 路径处理
			if gstr.HasPrefix(nutsDbPath, "./") {
				pwd, err := os.Getwd()
				if err != nil {
					nutsDbCacheErr = err
				}
				realPath = pwd + gstr.TrimLeftStr(nutsDbPath, ".")
			} else {
				realPath = nutsDbPath
			}
			g.Log().Line(false).Info("启用自定义NutsDb缓存数据库安装路径，安装路径：" + realPath)
		}

		// 默认配置
		opt := nutsdb.DefaultOptions

		// 自动创建数据库
		opt.Dir = nutsDbPath
		// 开启数据库
		db, err := nutsdb.Open(opt)
		if err != nil {
			nutsDbCacheErr = err
		}
		nutsDbCache = db

		g.Log().Line(false).Info("完成NutsDb缓存数据库实例化。")
	})

	return nutsDbCache, nutsDbCacheErr
}

// InitSetNutsDbGlobalPath 初始化NutsDb数据库安装路径，仅初始化一次，默认安装位置：/tmp/hcache/nutsdb/.cache/db，建议：在init方法下初始化NutsDb安装路径
func InitSetNutsDbGlobalPath(dir string) error {
	if nutsDbCache != nil {
		g.Log().Line(false).Error("初始化NutsDb数据库安装路径失败，已初始化位置：" + nutsDbPath)
		return errors.New("初始化NutsDb缓存数据库安装路径失败，已初始化位置：" + nutsDbPath)
	}
	// 更改默认安装路径
	nutsDbPath = dir
	return nil
}

// GetInitNutsDbGlobalPath 获取NutsDb数据库安装路径，当且仅当数据库实例化才可获取安装路径
func GetInitNutsDbGlobalPath() (string, error) {
	if nutsDbCache == nil {
		return "", errors.New("NutsDb缓存数据库还未初始化")
	}
	return nutsDbPath, nil
}

// SetCache 添加或更新缓存数据
func SetCache(bucket string, key []byte, value []byte, ttl uint32) error {
	// 获取实例单例缓存
	db, err := GetNutsDbCacheBean()
	if err != nil {
		return err
	}
	// 添加或更新缓存数据
	if err = db.Update(
		func(tx *nutsdb.Tx) error {
			if err = tx.Put(bucket, key, value, ttl); err != nil {
				return err
			}
			return nil
		}); err != nil {
		return err
	}
	return err
}

// GetCache 获取缓存数据
func GetCache(bucket string, key []byte) (entry *nutsdb.Entry, err error) {
	// 获取实例单例缓存
	db, err := GetNutsDbCacheBean()
	if err != nil {
		return nil, err
	}
	// 获取缓存数据
	if err = db.View(
		func(tx *nutsdb.Tx) error {
			if entry, err = tx.Get(bucket, key); err != nil {
				return err
			}
			return nil

		}); err != nil {
		return nil, err
	}

	return entry, err
}

// DelCache 删除缓存数据
func DelCache(bucket string, key []byte) (err error) {
	// 获取实例单例缓存
	db, err := GetNutsDbCacheBean()
	if err != nil {
		return err
	}
	// 删除缓存
	if err = db.Update(
		func(tx *nutsdb.Tx) error {
			if err = tx.Delete(bucket, key); err != nil {
				return err
			}
			return nil
		}); err != nil {
		return err
	}

	return err
}

// ScanByKeyPrefix 前缀扫描缓存
func ScanByKeyPrefix(bucket string, keyPrefix string, offset int, rows int) (entries nutsdb.Entries, off int, err error) {

	// 获取单例缓存数据库实例
	db, err := GetNutsDbCacheBean()
	if err != nil {
		return nil, offset, err
	}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			// 从指定偏移量开始查询缓存数据
			if entries, off, err = tx.PrefixScan(bucket, []byte(keyPrefix), offset, rows); err != nil {
				return err
			}
			return nil
		}); err != nil {
		return nil, offset, err
	}

	return entries, off, err
}
