package gfcasbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

// 测试casbin适配器
func CabinAdapterTest(e *casbin.Enforcer) {

	e.EnableAutoSave(true)

	// 添加策略
	if ok, _ := e.AddPolicy("admin", "/api/v1/hello", "GET"); !ok {
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

}

// 测试casbin适配器（单例）
func CabinAdapterBeanTest(e *casbin.SyncedEnforcer) {

	e.EnableAutoSave(true)

	// 添加策略
	if ok, _ := e.AddPolicy("admin", "/api/v1/hello", "GET"); !ok {
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

}
