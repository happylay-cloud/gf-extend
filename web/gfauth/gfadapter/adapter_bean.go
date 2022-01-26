package gfadapter

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// Enforcer 单例casbin同步执行器
var Enforcer *casbin.SyncedEnforcer

// EnforcerErr 单例casbin执行器错误对象
var EnforcerErr error

// 只操作一次对象
var once sync.Once

// NewEnforcerBean 实例化casbin执行器bean
//  1、支持自动注册，自动寻找gf框架default分组数据源（gdb.DB），无需关心数据源种类。
//  2、支持自定义分组数据源注册。
//
//  目前支持sqlite3、mysql5.7、postgresql数据库。
//
//  备注：1.sqlite3、mysql5.7数据库表新增主键自增，postgresql数据库无主键。
//       2.sqlite3、pgsql需要添加额外驱动
//
//	sqlite3驱动：	_ "github.com/lib/pq"
//
//	pgsql驱动：		_ "github.com/mattn/go-sqlite3"
//
//  示例：
//  e, err := gfadapter.NewEnforcerBean()
//  e, err := gfadapter.NewEnforcerBean(g.DB())
//  e, err := gfadapter.NewEnforcerBean(g.DB("sqlite"))
//  e, err := gfadapter.NewEnforcerBean(g.DB("mysql"))
//  e, err := gfadapter.NewEnforcerBean(g.DB("pgsql"))
func NewEnforcerBean(customDb ...gdb.DB) (*casbin.SyncedEnforcer, error) {

	once.Do(func() {
		// 定义数据源
		var db gdb.DB
		// 判断是否传参数
		num := len(customDb)
		if num == 0 {
			db = g.DB()
		} else {
			db = customDb[0]
		}
		// 创建同步casbin执行器
		_, EnforcerErr = newSyncedEnforcer(db)
	})
	return Enforcer, EnforcerErr
}

// newSyncedEnforcer 创建cacbin同步执行器
func newSyncedEnforcer(db gdb.DB) (*casbin.SyncedEnforcer, error) {

	// 打印日志
	g.Log().Line().Debug("实例化NewEnforcerBean")
	// TODO 需要添加动态模型配置，从配置文件中读取，获取不到则使用默认配置
	// 从字符串中加载模型
	modelFromString, _ := getDefaultNewModel()
	// 创建gf默认数据源适配器
	if adapter, err := NewAdapterByGdb(db); err == nil {
		// 调用已有连接的适配器中的构造器
		if options, err := NewAdapterFromOptions(adapter); err == nil {
			// 返回casbin执行器
			enforcer, err := casbin.NewSyncedEnforcer(modelFromString, options)
			// 生成单例执行器
			Enforcer = enforcer
			return enforcer, err
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

}
