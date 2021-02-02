package hredis

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/happylay-cloud/gf-extend/common/hdto"
)

// ZSetAdd 新增ZSet数据
//  @key   集合键
//  @value 值
//  @score 分数（小数）
func ZSetAdd(key interface{}, value interface{}, score interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("ZADD", key, score, value)
}

// ZSetCountAll 查询ZSet总数
//  @key 集合键
func ZSetCountAll(key interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("ZCARD", key)
}

// ZSetPage ZSet集合分页
//  @key      集合键
//  @pageNum  当前页码
//  @pageSize 每页记录数
//  返回分页对象，注意list集合是*gvar.Var类型
func ZSetPage(key interface{}, pageNum, pageSize int, needPageCount ...bool) (*hdto.Page, error) {

	// 创建分页对象
	page := hdto.NewPage(pageNum, pageSize)
	start, stop := page.RedisPage()

	// 获取分页数据
	list, err := g.Redis().DoVar("ZRANGE", key, start, stop)
	if err != nil {
		return nil, err
	}

	// 默认不计算总分页数
	if len(needPageCount) > 0 && needPageCount[0] {
		// 获取总记录数
		count, err := g.Redis().DoVar("ZCARD", key)
		if err != nil {
			return nil, err
		}
		// 计算总分页数
		pageCount := page.CalPageCount(count.Int())
		page.PageCount = pageCount
	}

	page.List = list
	return page, nil
}

// ZSetPage ZSet集合删除指定元素
//  @key   集合键
//  @value 要删除的值
func ZSetRemove(key interface{}, value interface{}) (*gvar.Var, error) {
	return g.Redis().DoVar("ZREM", key, value)
}
