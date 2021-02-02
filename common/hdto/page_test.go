package hdto

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/frame/g"
)

func TestPage(t *testing.T) {

	// newPage := NewPage(-4, -3, false)

	newPage := NewPage(2, 10)

	g.Dump(newPage)

	g.Dump("分页总数", newPage.CalPageCount(10))

	fmt.Println("mysql分页")
	g.Dump(newPage.LimitPage())

	fmt.Println("redis分页")
	g.Dump(newPage.RedisPage())

}
