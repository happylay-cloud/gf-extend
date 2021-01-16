package gfmiddleware

import (
	"github.com/gogf/gf/net/ghttp"
)

// MiddlewareCORS 跨域中间件
//  示例：
// 	s.Group("/", func(group *ghttp.RouterGroup) {
//		group.Middleware(gfmiddleware.MiddlewareCORS)
//	})
func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}


