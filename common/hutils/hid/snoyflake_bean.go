package hid

import (
	"time"

	"github.com/gogf/gf/frame/g"
)

// IdBean 全局分布式唯一id实例，IdBean.NextID()获取唯一id。
var IdBean *SonyFlake

func init() {
	g.Log().Line().Debug("初始化snowflake")
	parse, err := time.Parse("2006-01-02", "1970-01-01")
	if err != nil {
		g.Log().Line(false).Error(err)
		panic("初始化snowflake异常")
	}

	settings := Settings{
		StartTime: parse,
	}
	IdBean = NewSonyFlake(settings)
}
