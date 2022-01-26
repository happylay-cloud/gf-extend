package hstr

import (
	"fmt"
	"testing"
)

func TestTrimBlank(t *testing.T) {
	fmt.Println(TrimBlank(" \t这是一个去除首尾空白字符的测试。 \n "))
}
