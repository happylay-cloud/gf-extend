package hredis

import (
	"errors"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

// https://www.redis.net.cn

// Set 新增string类型
//  @key   键
//  @value 值
//  @expire 过期时间，单位秒（默认不设置，不推荐使用，建议使用SetAndExpire方法）
func Set(key interface{}, value interface{}, expire ...uint32) (*gvar.Var, error) {

	// 设置键值
	doVar, err := g.Redis().DoVar("SET", key, value)
	// 设置过期时间（默认不设置）
	if len(expire) > 0 {
		// 设置键值成功后，设置过期时间
		if err == nil {
			doVar, err = g.Redis().DoVar("EXPIRE", key, expire[0])
			if err != nil {
				return nil, err
			}
		}
	}
	return doVar, err
}

// SetAndExpire 设置key-value并设置过期时间，限定为字符串
//  @key     键
//  @value   值
//  @expire 过期时间，单位秒
func SetAndExpire(key interface{}, value interface{}, expire uint32) (*gvar.Var, error) {
	// 设置键值
	return g.Redis().DoVar("SETEX", key, expire, value)
}

// Setnx 只有在key不存在时设置key的值（常用于分布式锁）
//  设置成功，返回1，设置失败，返回0。
//  @key   键
//  @value 值
func Setnx(key interface{}, value interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("SETNX", key, value)
}

// Get 获取string值
//  @key   键
func Get(key interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("SET", key)
}

// BatchValues 模糊查询key对应所有值（注意，结果集要判空处理）
//  @likeKey   键
func BatchValues(likeKey interface{}) (*gvar.Var, error) {
	// 查询所有键
	keys, err := g.Redis().DoVar("KEYS", likeKey)
	if err != nil {
		return nil, err
	}

	// 判断键是否存在
	if keys.IsEmpty() {
		return nil, errors.New("键" + gconv.String(likeKey) + "不存在")
	}

	// 查询所有键对应的值
	return g.Redis().DoVar("MGET", keys.Interfaces()...)

}

// DeleteKeyValue 删除key及value
//  @key 键
func DeleteKeyValue(key interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("DEL", key)
}
