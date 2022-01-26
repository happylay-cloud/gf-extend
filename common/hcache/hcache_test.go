package hcache

import (
	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"github.com/gogf/gf/frame/g"
	"github.com/xujiajun/nutsdb"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

// testOlricCacheValue 测试用缓存结构体
type testOlricCacheValue struct {
	Value string
}

func TestOlricCache(t *testing.T) {

	// 配置本地缓存
	c := config.New("local")

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	c.Started = func() {
		// 关闭上下文
		defer cancel()
		fmt.Println("[INFO] 开启缓存数据库")
	}

	db, err := olric.New(c)
	if err != nil {
		panic(err)
	}

	go func() {
		// 启动缓存数据库
		err = db.Start()
		if err != nil {
			panic(err)
		}
	}()

	// 等待协程结束
	<-ctx.Done()

	// 创建map缓存类型
	dm, err := db.NewDMap("cacheMap")
	if err != nil {
		panic(err)
	}

	key := "k1"

	jsonByte, err := json.Marshal(&testOlricCacheValue{
		Value: "测试缓存值",
	})
	if err != nil {
		return
	}

	g.Dump("存储缓存值：", string(jsonByte))
	err = dm.Put(key, string(jsonByte))

	value, err := dm.Get(key)
	if err != nil {
		panic(err)
	}

	// 接口断言
	v1, ok := value.(string)
	if ok {
		cacheV1 := testOlricCacheValue{}
		err := json.Unmarshal([]byte(v1), &cacheV1)
		if err != nil {
			return
		}
		g.Dump("缓存中读取：", cacheV1)
	}

	// 关闭缓存数据库
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = db.Shutdown(ctx)

}

// 相关文档：https://www.bookstack.cn/read/NutsDB-zh/spilt.8.spilt.3.NutsDB.md
func TestNutsDbCache(t *testing.T) {

	// 默认配置
	opt := nutsdb.DefaultOptions
	// 自动创建数据库
	opt.Dir = "./.cache/nutsdb"
	// 开启数据库
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *nutsdb.DB) {
		err := db.Close()
		if err != nil {
			g.Log().Line(false).Error("数据库关闭异常：", err.Error())
		}
	}(db)

	// 添加或更新数据
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte("k1")
			val := []byte("v1")
			bucket := "bucket_test"
			if err := tx.Put(bucket, key, val, 0); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Fatal(err)
	}

	// 获取数据
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte("k1")
			bucket := "bucket_test"
			if e, err := tx.Get(bucket, key); err != nil {
				return err
			} else {
				fmt.Println(string(e.Value))
			}
			return nil
		}); err != nil {
		log.Println(err)
	}

	// 删除数据
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte("k1")
			bucket := "bucket_test"
			if err := tx.Delete(bucket, key); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Fatal(err)
	}

}

func TestNutsDbCacheUpdate(t *testing.T) {
	// 默认配置
	opt := nutsdb.DefaultOptions
	// 自动创建数据库
	opt.Dir = "./.cache/nutsdb"
	// 开启数据库
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *nutsdb.DB) {
		err := db.Close()
		if err != nil {
			g.Log().Line(false).Error("数据库关闭异常：", err.Error())
		}
	}(db)

	// 添加或更新数据
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte("ttl_k1")
			val := []byte("过期值")
			bucket := "bucket_test"
			if err := tx.Put(bucket, key, val, 30); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Fatal(err)
	}
}

func TestNutsDbCacheQuery(t *testing.T) {
	// 默认配置
	opt := nutsdb.DefaultOptions
	// 自动创建数据库
	opt.Dir = "./.cache/nutsdb"
	// 开启数据库
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *nutsdb.DB) {
		err := db.Close()
		if err != nil {
			g.Log().Line(false).Error("数据库关闭异常：", err.Error())
		}
	}(db)

	// 获取数据
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte("ttl_k1")
			bucket := "bucket_test"
			if e, err := tx.Get(bucket, key); err != nil {
				return err
			} else {
				fmt.Println(string(e.Value))
			}
			return nil
		}); err != nil {
		log.Println(err)
	}
}

func TestNutsDbCacheMerge(t *testing.T) {
	// 默认配置
	opt := nutsdb.DefaultOptions
	// 自动创建数据库
	opt.Dir = "./.cache/nutsdb"
	// 开启数据库
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *nutsdb.DB) {
		err := db.Close()
		if err != nil {
			g.Log().Line(false).Error("数据库关闭异常：", err.Error())
		}
	}(db)

	err = db.Merge()
	if err != nil {
		log.Fatal(err)
	}
}

func TestNutsDbCacheBackup(t *testing.T) {
	// 默认配置
	opt := nutsdb.DefaultOptions
	// 自动创建数据库
	opt.Dir = "./.cache/nutsdb"
	// 开启数据库
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *nutsdb.DB) {
		err := db.Close()
		if err != nil {
			g.Log().Line(false).Error("数据库关闭异常：", err.Error())
		}
	}(db)

	err = db.Backup(".cache/backup")
	if err != nil {
		log.Fatal(err)
	}
}
