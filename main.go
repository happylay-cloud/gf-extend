package main

import (
	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/happylay-cloud/gf-extend/common/gfres"
	"github.com/happylay-cloud/gf-extend/web/auth/gfadapter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	s := g.Server()

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", func(r *ghttp.Request) {

			sql := gfadapter.CreateSqlite3Table("sacbin_rule")
			exec, _ := g.DB().Exec(sql)

			a, _ := gfadapter.
				NewAdapterFromOptions(gfadapter.NewAdapterByGdb(g.DB()))

			e, _ := casbin.NewEnforcer("./web/auth/gfadapter/config/rbac_model.conf", a)

			g.Dump(exec, e)
			g.Dump(e.GetModel())

			_, _ = e.AddPolicy("张三", "/web/spi/v1", "只读")

			g.Dump(e.GetAllSubjects())
			g.Dump(e.GetAllRoles())
			g.Dump(e.GetAllObjects())
			g.Dump(e.GetAllActions())

			gfres.OkWithData(r, sql)
		})
	})

	s.Run()
}
