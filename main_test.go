package gfextendtest

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"github.com/happylay-cloud/gf-extend/common/gfres"
	"github.com/happylay-cloud/gf-extend/web/auth/gfadapter"
	"github.com/happylay-cloud/gf-extend/web/auth/gfcasbin"
)

func TestGf(t *testing.T) {

	s := g.Server()

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/v1", testNewCasbin)
		group.ALL("/v2", testNewCasbinBean)
	})

	s.Run()
}

func testNewCasbin(r *ghttp.Request) {

	// 实例化casbin执行器
	//e, err := gfadapter.NewEnforcer(g.DB("pgsql"))
	//e, err := gfadapter.NewEnforcer(g.DB("mysql"))
	e, err := gfadapter.NewEnforcer(g.DB("sqlite"))

	if err != nil {
		gfres.FailWithEx(r, err.Error())
		return
	}

	gfcasbin.CabinAdapterTest(e)

	gfres.OkWithData(r, g.Map{
		"sub":   e.GetAllSubjects(),
		"obj":   e.GetAllObjects(),
		"act":   e.GetAllActions(),
		"role":  e.GetAllRoles(),
		"model": e.GetModel(),
	})

}

func testNewCasbinBean(r *ghttp.Request) {

	// 实例化casbin执行器
	//e, err := gfadapter.NewEnforcerBean(g.DB("pgsql"))
	//e, err := gfadapter.NewEnforcerBean(g.DB("mysql"))
	//e, err := gfadapter.NewEnforcerBean(g.DB("sqlite"))

	e, err := gfadapter.NewEnforcerBean()

	if err != nil {
		gfres.FailWithEx(r, err.Error())
		return
	}

	gfcasbin.CabinAdapterBeanTest(e)

	gfres.OkWithData(r, g.Map{
		"sub":   e.GetAllSubjects(),
		"obj":   e.GetAllObjects(),
		"act":   e.GetAllActions(),
		"role":  e.GetAllRoles(),
		"model": e.GetModel(),
	})

}
