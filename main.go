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
		group.ALL("/testSqlite3", testSqlite3)
		group.ALL("/", testNewCasbin)
	})

	s.Run()
}

func testSqlite3(r *ghttp.Request) {
	sql := gfadapter.CreateSqlite3Table("sacbin_rule")
	g.Dump(g.Cfg())
	exec, _ := g.DB().Exec(sql)
	g.Dump(exec)
}

func testNewCasbin(r *ghttp.Request) {

	e, err := gfadapter.NewEnforcer(g.DB())
	if err != nil {
		gfres.FailWithEx(r, err.Error())
	}
	_, _ = e.AddPolicy("张三", "/web/spi/v1", "只读")
	gfres.OkWithData(r, e.GetAllObjects())
}
