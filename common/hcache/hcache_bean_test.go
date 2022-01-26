package hcache

import (
	"github.com/gogf/gf/frame/g"
	"github.com/xujiajun/nutsdb"

	"fmt"
	"log"
	"testing"
)

func TestNutsDbCacheBean(t *testing.T) {

	// 初始化缓存安装路径（覆盖默认安装路径）
	err := InitSetNutsDbGlobalPath("./.cache/nutsdb")
	if err != nil {
		return
	}

	// 获取单例缓存bean
	db, err := GetNutsDbCacheBean()
	if err != nil {
		log.Println(err)
		return
	}

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

func TestNutsDbSetCache(t *testing.T) {
	err := SetCache("b1", []byte("k1"), []byte("测试"), 20)
	if err != nil {
		fmt.Println(err)
	}
}

func TestNutsDbGetCache(t *testing.T) {
	entry, err := GetCache("b1", []byte("k1"))
	if err != nil {
		return
	}
	g.Dump(string(entry.Value))
}
