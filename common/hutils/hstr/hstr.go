package hstr

import (
	"bytes"
	"github.com/gogf/gf/text/gstr"
	"strconv"
	"strings"
)

// TrimBlank 删除字符串首尾空白
func TrimBlank(str string) string {
	// 去除空白字符及特殊空白
	trimStr := strings.TrimSpace(str)
	// 去除制表符
	trimStr = gstr.TrimStr(trimStr, "\t")
	// 去除换行符
	trimStr = gstr.TrimStr(trimStr, "\n")
	// 去除首尾空白字符
	return gstr.Trim(trimStr)
}

// ZhToUnicode 中文转unicode
func ZhToUnicode(zhContent string) string {
	textQuoted := strconv.QuoteToASCII(zhContent)
	// 去除引号
	return textQuoted[1 : len(textQuoted)-1]
}

// SpecialEmptyToNormalEmpty 特殊空白转换成普通空白
func SpecialEmptyToNormalEmpty(content string) string {
	// 转换成unicode编码
	unicodeText := ZhToUnicode(content)
	// 转换全角空格
	unicodeText = gstr.Replace(unicodeText, "\\u00a0", " ")
	// 转换半角空格
	unicodeText = gstr.Replace(unicodeText, "\\u2003", " ")
	return UnicodeToZh(unicodeText)
}

// UnicodeToZh unicode转中文
func UnicodeToZh(unicodeUnquoted string) string {
	buf := bytes.NewBuffer(nil)
	i, j := 0, len(unicodeUnquoted)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(unicodeUnquoted[i:])
			break
		}
		if unicodeUnquoted[i] == '\\' && unicodeUnquoted[i+1] == 'u' {
			hex := unicodeUnquoted[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(unicodeUnquoted[i:x])
			}
			i = x
		} else {
			buf.WriteByte(unicodeUnquoted[i])
			i++
		}
	}
	return buf.String()
}
