package hjsoup

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestHttpClient(t *testing.T) {

	productCode, err := SearchByProductCode("6921168509256", false)
	if err != nil {
		fmt.Println(err)
	}
	g.Dump(productCode)
}
