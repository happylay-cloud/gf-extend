package hjsoup

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestHttpClient(t *testing.T) {

	productCode, err := SearchByProductCode("6917751470508", false)
	if err != nil {
		fmt.Println(err)
	}
	g.Dump(productCode)

}
