package gfadapter

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gogf/gf/frame/g"
)

// 测试gf-casbin适配器
func testNewEnforcer(t *testing.T) {

	// 实例化casbin执行器
	//e, err := NewEnforcer(g.DB("pgsql"))
	//e, err := NewEnforcer(g.DB("mysql"))
	e, err := NewEnforcer(g.DB("sqlite"))

	if err != nil {
		g.Log().Error(err)
		return
	}
	e.EnableAutoSave(true)

	// 添加策略
	if ok, _ := e.AddPolicy("admin", "/api/v1/hello", "GET"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	// 添加策略
	if ok, _ := e.AddPolicy("happylay", "/api/v1/hello", "POST"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	// 添加组策略
	if ok, _ := e.AddNamedGroupingPolicy("g", "eat", "干饭人"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	if ok, _ := e.AddNamedGroupingPolicy("g", "worker", "打工人"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	// 指定字段删除策略
	if ok, _ := e.RemoveFilteredNamedGroupingPolicy("g", 0, "worker", "打工人"); !ok {
		fmt.Println("策略不存在")
	} else {
		fmt.Println("删除成功")
	}

	// 删除策略
	if ok, _ := e.RemovePolicy("admin", "/api/v1/hello", "GET"); !ok {
		fmt.Println("策略不存在")
	} else {
		fmt.Println("删除成功")
	}

	// 获取策略
	list := e.GetPolicy()
	for _, vList := range list {
		fmt.Print("策略：")
		for _, v := range vList {
			fmt.Printf("%s, ", v)
		}
		fmt.Print("\n")
	}

	// 检查权限
	if ok, _ := e.Enforce("admin", "/api/v1/hello", "GET"); ok {
		fmt.Println("权限正常")
	} else {
		fmt.Println("没有权限")
	}

	m := g.Map{
		"sub":   e.GetAllSubjects(),
		"obj":   e.GetAllObjects(),
		"act":   e.GetAllActions(),
		"role":  e.GetAllRoles(),
		"model": e.GetModel(),
	}

	g.Dump(m)

	// Output:
}

// 测试gf-casbin适配器（单例）
func testNewEnforcerBean(t *testing.T) {

	// 实例化casbin执行器
	//e, err := NewEnforcerBean(g.DB("pgsql"))
	//e, err := NewEnforcerBean(g.DB("mysql"))
	e, err := NewEnforcerBean(g.DB("sqlite"))

	//e, err := NewEnforcerBean()

	if err != nil {
		g.Log().Error(err)
		return
	}
	e.EnableAutoSave(true)

	// 添加策略
	if ok, _ := e.AddPolicy("admin", "/api/v1/hello", "GET"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	// 添加策略
	if ok, _ := e.AddPolicy("happylay", "/api/v1/hello", "POST"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	// 添加组策略
	if ok, _ := e.AddNamedGroupingPolicy("g", "eat", "干饭人"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	if ok, _ := e.AddNamedGroupingPolicy("g", "worker", "打工人"); !ok {
		fmt.Println("策略已经存在")
	} else {
		fmt.Println("增加成功")
	}

	// 指定字段删除策略
	if ok, _ := e.RemoveFilteredNamedGroupingPolicy("g", 0, "worker", "打工人"); !ok {
		fmt.Println("策略不存在")
	} else {
		fmt.Println("删除成功")
	}

	// 删除策略
	if ok, _ := e.RemovePolicy("admin", "/api/v1/hello", "GET"); !ok {
		fmt.Println("策略不存在")
	} else {
		fmt.Println("删除成功")
	}

	// 获取策略
	list := e.GetPolicy()
	for _, vList := range list {
		fmt.Print("策略：")
		for _, v := range vList {
			fmt.Printf("%s, ", v)
		}
		fmt.Print("\n")
	}

	// 检查权限
	if ok, _ := e.Enforce("admin", "/api/v1/hello", "GET"); ok {
		fmt.Println("权限正常")
	} else {
		fmt.Println("没有权限")
	}

	m := g.Map{
		"sub":   e.GetAllSubjects(),
		"obj":   e.GetAllObjects(),
		"act":   e.GetAllActions(),
		"role":  e.GetAllRoles(),
		"model": e.GetModel(),
	}

	g.Dump(m)

	// Output:
}

// 测试非单例模式
func TestNewEnforcer(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println("###################################🔥开始新NewEnforcer测试START###################################")
		testNewEnforcer(t)
		fmt.Println("###################################🚀结束新NewEnforcer测试END#####################################")
	}
}

// 测试单例模式
func TestTestNewEnforcerBean(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println("###################################🔥开始新NewEnforcerBean测试START###################################")
		testNewEnforcerBean(t)
		fmt.Println("###################################🚀结束新NewEnforcerBean测试END####################################")
	}
}
