package gfresponse

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 自定义响应枚举
const (
	SUCCESS        = 200 // 操作成功
	FAIL           = 201 // 操作失败
	SERVER_ERROR   = 500 // 服务器内部错误
	BUSINESS_ERROR = 501 // 通用业务异常
)

// 数据返回通用JSON数据结构
type JsonResponse struct {
	Code      int         `json:"code"`      // 状态码（200：操作成功，201：操作失败，500：服务器内部错误，501：通用业务异常）
	Message   string      `json:"message"`   // 提示信息
	Data      interface{} `json:"data"`      // 响应数据
	Exception interface{} `json:"exception"` // 异常信息
}

// 内部调用
func result(r *ghttp.Request, code int, message string, data interface{}, exception interface{}) {
	if err := r.Response.WriteJson(JsonResponse{
		code,
		message,
		data,
		exception,
	}); err != nil {
		g.Log().Error(err)
	}
	// 仅退出当前执行的逻辑方法
	r.Exit()
}

// 操作成功
func Ok(r *ghttp.Request) {
	result(r, SUCCESS, "操作成功", map[string]interface{}{}, nil)
}

// 操作成功（自定义消息）
func OkWithMsg(r *ghttp.Request, message string) {
	result(r, SUCCESS, message, map[string]interface{}{}, nil)
}

// 操作成功（自定义数据）
func OkWithData(r *ghttp.Request, data interface{}) {
	result(r, SUCCESS, "操作成功", data, nil)
}

// 操作成功（自定义消息和数据）
func OkWithDataAndMsg(r *ghttp.Request, message string, data interface{}) {
	result(r, SUCCESS, message, data, nil)
}

// 操作失败
func Fail(r *ghttp.Request) {
	result(r, FAIL, "操作失败", map[string]interface{}{}, nil)
}

// 操作失败（自定义消息）
func FailWithMsg(r *ghttp.Request, message string) {
	result(r, FAIL, message, map[string]interface{}{}, nil)
}

// 操作失败（自定义异常信息）
func FailWithEx(r *ghttp.Request, exception interface{}) {
	result(r, FAIL, "操作失败", nil, exception)
}

// 操作失败（自定义错误码和消息）
func FailWithCodeAndMsg(r *ghttp.Request, code int, message string) {
	result(r, code, message, nil, nil)
}

// 操作失败（自定义错误码、消息、异常信息）
func FailWithCodeMsgAndEx(r *ghttp.Request, code int, message string, exception interface{}) {
	result(r, code, message, nil, exception)
}

// 标准返回结果数据结构封装
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	err := r.Response.WriteJson(JsonResponse{
		Code:    code,
		Message: message,
		Data:    responseData,
	})

	if err != nil {
		g.Log().Error(err)
	}
}

// 返回JSON数据并退出当前HTTP执行函数
func JsonExit(r *ghttp.Request, err int, msg string, data ...interface{}) {
	Json(r, err, msg, data...)
	r.Exit()
}
