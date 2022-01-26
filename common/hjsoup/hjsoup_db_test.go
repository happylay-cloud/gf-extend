package hjsoup

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"github.com/gogf/gf/frame/g"
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	// 商品条码
	productCode := "6921168509256"

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
	dm, err := db.NewDMap("ProductCodeInfoMap")
	if err != nil {
		panic(err)
	}

	value, err := dm.Get(productCode)
	if err == nil {
		// 接口断言
		v, ok := value.(string)
		if ok {
			productCodeDto := ProductCodeDto{}
			err := json.Unmarshal([]byte(v), &productCodeDto)
			if err != nil {
				return
			}
			g.Dump("缓存中读取：", productCodeDto)
		}

		return
	}

	// 查询数据
	productCodeInfo, err := SearchByProductCode(productCode, false)
	if err != nil {
		fmt.Println(err)
	}
	jsonByte, err := json.Marshal(productCodeInfo)
	if err != nil {
		return
	}

	g.Dump("存储缓存值：", string(jsonByte))
	err = dm.Put(productCodeInfo.ProductCode, string(jsonByte))

	value, err = dm.Get(productCodeInfo.ProductCode)
	if err != nil {
		panic(err)
	}

	// 接口断言
	v, ok := value.(string)
	if ok {
		productCodeDto := ProductCodeDto{}
		err := json.Unmarshal([]byte(v), &productCodeDto)
		if err != nil {
			return
		}
		g.Dump("缓存中读取：", productCodeDto)
	}

	// 关闭缓存数据库
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = db.Shutdown(ctx)

}
