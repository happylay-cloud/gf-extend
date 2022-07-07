package hjsoup

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestHeFeiFangJia(t *testing.T) {

	viewState, err := GetHeFeiFangJiaRecordViewState()
	if err == nil {
		page, _ := ListHeFeiFangJiaRecordPage(viewState, 3)
		fmt.Println(len(page))
		g.Dump(page)
	}

}
