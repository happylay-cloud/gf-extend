package main

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"github.com/happylay-cloud/gf-extend/common/gfres"
	"github.com/happylay-cloud/gf-extend/web/auth/gfadapter"
)

func main() {

	s := g.Server()

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", testNewCasbin)
	})

	s.Run()
}

func testNewCasbin(r *ghttp.Request) {

	e, err := gfadapter.NewEnforcer()
	if err != nil {
		gfres.FailWithEx(r, err.Error())
		return
	}
	// 添加策略
	_, _ = e.AddPolicy("happylay", "/api", "GET")

	gfres.OkWithData(r, g.Map{
		"sub":  e.GetAllSubjects(),
		"obj":  e.GetAllObjects(),
		"act":  e.GetAllActions(),
		"role": e.GetAllRoles(),
	})

}
