package hid

import (
	"time"

	"github.com/gogf/gf/frame/g"
)

// IdBean 全局分布式唯一id实例，IdBean.NextID()获取唯一id。
var IdBean *SonyFlake

// 初始化snowflake
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

// GetIdGenTime 获取Id生成时间
//  @id SnowFlake生成的唯一Id
func GetIdGenTime(id uint64) string {
	sTime := id >> (BitLenSequence + BitLenMachineID)
	return time.Unix(int64(sTime*10/1000), 0).Format("2006-01-02 15:04:05")
}
