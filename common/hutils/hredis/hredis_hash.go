package hredis

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
)

// HashSet 将哈希表key中的字段field的值设为value
//  @key   哈希表key
//  @field 字段
//  @value 值
func HashSet(key interface{}, field interface{}, value interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("HSET", key, field, value)
}

// HashGet 获取存储在哈希表中指定字段的值
//  @key   哈希表key
//  @field 字段
func HashGet(key interface{}, field interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("HGET", key, field)
}

// HashDel 删除哈希表key中的一个指定字段，不存在的字段将被忽略
//  @key   哈希表key
//  @field 字段
func HashDel(key interface{}, field interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("HDEL", key, field)
}
