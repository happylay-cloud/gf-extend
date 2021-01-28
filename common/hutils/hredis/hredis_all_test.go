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

func TestHRedisString(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	_, err := Set("hredis-string", "k", 100)
	if err != nil {
		fmt.Println(err)
	}

}

func TestHRedisBatchAll(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	list, err := BatchValues("hredis-*")
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(list)
}

func TestHRedisValueDel(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	result, err := DeleteKeyValue("hredis-string")
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(result)
}

func TestHRedisKeyExp(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	result, err := SetAndExpire("hredis-exp", "exp", 100)
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(result)
}

func TestHRedisSetnx(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	result, err := Setnx("hredis-setnx", "exp")
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(result)
}

func TestHRedisHashSet(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	result, err := HashSet("hredis-hash", "key", "value")
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(result)
}

func TestHRedisHashGet(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	result, err := HashGet("hredis-hash", "key")
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(result)
}

func TestHRedisHashDel(t *testing.T) {

	config := gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}

	gredis.SetConfig(config)

	result, err := HashDel("hredis-hash", "key")
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(result)
}
