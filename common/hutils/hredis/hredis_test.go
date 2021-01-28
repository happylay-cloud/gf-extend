package hredis

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
)

func TestHRedis(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	_, err := ZSetAdd("hredis-zset", "k2", 123)
	if err != nil {
		fmt.Println(err)
	}

}

func TestHRedisPage(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	page, err := ZSetPage("hredis-zset", 2, 2, true)
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(page)
}

func TestHRedisRemove(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	_, err := ZSetRemove("hredis-zset", "k1")
	if err != nil {
		fmt.Println(err)
	}

}
