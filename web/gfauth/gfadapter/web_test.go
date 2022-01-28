package gfadapter

import (
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/happylay-cloud/gf-extend/common/hresp"
)

func TestGf(t *testing.T) {

	s := g.Server()

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", hi)
	})

	s.Run()
}

func hi(r *ghttp.Request) {
	hresp.Ok(r)
}
