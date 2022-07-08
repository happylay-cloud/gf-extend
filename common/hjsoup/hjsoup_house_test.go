package hjsoup

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestHeFeiFangJiaPage(t *testing.T) {

	viewState1, totalPage, err := GetHeFeiFangJiaRecordViewState()
	if err == nil {
		page, _ := ListHeFeiFangJiaRecordPage(viewState1, 3)
		fmt.Println("总分页：", totalPage, "当前条数：", len(page))
		g.Dump(page)
	}

}

func TestHeFeiFangJiaDetail(t *testing.T) {

	hrefId := "4"

	detail, viewState, totalPage, err := GetHeFeiFangJiaDetail(hrefId)
	if err == nil {
		g.Dump("详情：", detail, "总分页：", totalPage, "状态：", viewState)

		page, err := ListHeFeiFangJiaHousePage(viewState, hrefId, 1)
		if err == nil {
			g.Dump(len(page), page)
		}
	} else {
		g.Dump(err.Error())
	}

}

func TestHeFeiFangJiaRelease(t *testing.T) {

	releaseTime, releaseSource, err := GetHeFeiFangJiaRecordReleaseInfo()
	if err == nil {
		fmt.Println(releaseTime, releaseSource)
	}

}
