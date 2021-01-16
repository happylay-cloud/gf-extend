package gfmiddleware

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/happylay-cloud/gf-extend/common/gfresponse"
)

// MiddlewareErrorHandler 全局后置中间件捕获异常
func MiddlewareErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		// 记录到自定义错误日志文件
		g.Log("exception").Error(err)
		// 清除系统异常响应
		r.Response.ClearBuffer()
		// 返回自定义异常响应
		gfresponse.FailWithEx(r, err.Error())
	}
}
