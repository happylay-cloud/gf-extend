package hjsoup

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestHeFeiFangJiaPage(t *testing.T) {

	viewState1, err := GetHeFeiFangJiaRecordViewState()
	if err == nil {
		page, _ := ListHeFeiFangJiaRecordPage(viewState1, 3)
		fmt.Println(len(page))
		g.Dump(page)
	}

}

func TestHeFeiFangJiaDetail(t *testing.T) {

	detail, viewState, err := GetHeFeiFangJiaDetail("8324")
	if err == nil {
		g.Dump(detail, viewState)
	}

}
