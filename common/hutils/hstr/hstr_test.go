package hstr

import (
	"fmt"
	"strings"
	"testing"
)

func TestTrimBlank(t *testing.T) {
	fmt.Println(TrimBlank(" \t这是一个去除首尾空白字符的测试。 \n "))
}

func TestZhBlank(t *testing.T) {
	// 开始 占位符 占位符 结束
	text := " 开始 占位符 占位符 结束 "

	// 转换普通字符
	normalEmpty := SpecialEmptyToNormalEmpty(text)
	fmt.Println(normalEmpty)
	fmt.Println(ZhToUnicode(normalEmpty))

	// 去除首尾空白
	emptyText := strings.TrimSpace(text)
	fmt.Println(emptyText)
	// 查看编码
	fmt.Println(ZhToUnicode(emptyText))

}
