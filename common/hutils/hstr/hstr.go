package hstr

import (
	"bytes"
	"github.com/gogf/gf/text/gstr"
	"strconv"
)

// TrimBlank 删除字符串首尾空白
func TrimBlank(str string) string {
	// 去除制表符
	trimStr := gstr.TrimStr(str, "\t")
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
