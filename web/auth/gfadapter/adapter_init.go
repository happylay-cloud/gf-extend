package gfadapter

import (
	"runtime"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

const (
	sqlite = "sqlite"
	mysql  = "mysql"
)

// 从自定义数据库连接中创建适配器
func NewAdapterByGdb(customDb gdb.DB) (*Adapter, error) {

	// 获取当前数据库类型
	dbType := g.DB().GetConfig().Type

	//  TODO 需要对不同类型数据库进行处理

	// 判断当前数据库类型
	if gstr.Equal(sqlite, dbType) { // 自动创建sqlite3数据库casbin_rule表
		sql := CreateSqlite3Table("casbin_rule")
		if _, err := g.DB().Exec(sql); err != nil {
			return nil, err
		}
	}
	// TODO ------------------------------------------------------------
	// 1.添加事务操作
	// 2.修改自动保存逻辑
	// 3.完善不同数据库操作
	// -----------------------------------------------------------------

	// 构造适配器对象
	a := &Adapter{
		db: customDb,
	}

	if customDb == nil {
		return nil, gerror.New("数据库默认连接不存在，无法实例化casbin执行器")
	}

	// 释放对象时调用析构函数。
	runtime.SetFinalizer(a, finalizer)

	return a, nil
}

// NewEnforcer 实例化gdb.DB默认数据源casbin执行器
func NewEnforcer(customDb ...gdb.DB) (*casbin.Enforcer, error) {

	// rbac_model.conf配置字符串
	rbacModelText :=
		`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	// 从字符串中加载模型
	modelFromString, _ := model.NewModelFromString(rbacModelText)

	// 定义数据源
	var db gdb.DB

	// 判断是否传参数
	num := len(customDb)

	if num == 0 {
		db = g.DB()
	} else {
		db = customDb[0]
	}

	// 创建gf默认数据源适配器
	if adapter, err := NewAdapterByGdb(db); err == nil {

		// 调用已有连接的适配器中的构造器
		if options, err := NewAdapterFromOptions(adapter); err == nil {
			// 返回casbin执行器
			return casbin.NewEnforcer(modelFromString, options)
		} else {
			return nil, err
		}

	} else {
		return nil, err
	}

}
