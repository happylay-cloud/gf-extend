package hstr

import "github.com/gogf/gf/text/gstr"

// TrimBlank 删除字符串首尾空白
func TrimBlank(str string) string {
	// 去除制表符
	trimStr := gstr.TrimStr(str, "\t")
	// 去除换行符
	trimStr = gstr.TrimStr(trimStr, "\n")
	// 去除首尾空白字符
	return gstr.Trim(trimStr)
}
