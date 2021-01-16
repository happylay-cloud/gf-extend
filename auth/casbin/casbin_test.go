package casbin_test

import (
	"fmt"
	"testing"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gogf/gf/frame/g"
)

// TestCasbin 权限测试
//  文档：
//  https://casbin.org/docs/zh-CN/get-started
//  适配器：
//  https://casbin.org/docs/zh-CN/adapters
//  示例：
//  https://github.com/casbin/casbin/tree/master/examples
func TestCasbin(t *testing.T) {

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
	// enforcer, _ := casbin.NewEnforcer("./config/rbac_model.conf", "./config/rbac_policy.csv")

	// 从字符串中加载模型
	modelFromString, err := model.NewModelFromString(rbacModelText)

	if err != nil {
		fmt.Println(err)
	}
	g.Dump(modelFromString)

	// 从文件中创建一个适配器
	newAdapter := fileadapter.NewAdapter("./config/rbac_policy.csv")

	enforcer, _ := casbin.NewEnforcer(modelFromString, newAdapter)

	// 添加策略
	if ok, _ := enforcer.AddPolicy("admin", "/api/v1/hello", "GET"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	// 删除策略
	if ok, _ := enforcer.RemovePolicy("admin", "/api/v1/hello", "GET"); !ok {
		fmt.Println("策略不存在")
	} else {
		fmt.Println("删除成功")
	}

	// 获取策略
	list := enforcer.GetPolicy()
	for _, vList := range list {
		fmt.Print("策略：")
		for _, v := range vList {
			fmt.Printf("%s, ", v)
		}
		fmt.Print("\n")
	}

	// 检查权限
	if ok, _ := enforcer.Enforce("李四", "B数据", "write"); ok {
		fmt.Println("权限正常")
	} else {
		fmt.Println("没有权限")
	}

	g.Dump(enforcer.GetAllSubjects())
	g.Dump(enforcer.GetAllRoles())
	g.Dump(enforcer.GetAllObjects())
	g.Dump(enforcer.GetAllActions())

}
