package hjsoup

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/frame/g"
)

func TestHttpClient(t *testing.T) {

	productCode, err := SearchByProductCode("6921168509256", false)
	if err != nil {
		fmt.Println(err)
	}
	g.Dump(productCode)
}

func TestSearchByProductCodeCache(t *testing.T) {
	// 商品条码
	productCode := "6928804011142"

	productCodeInfo, err := SearchByProductCodeCache(productCode)
	if err != nil {
		fmt.Println(err)
	}

	g.Dump(productCodeInfo)
}
