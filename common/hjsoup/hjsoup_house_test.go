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

func TestGetJsCookie(t *testing.T) {

	item := map[string]string{
		"bts":   "1657277174.162|0|A79,quRWXdetTIBGl73hRfU3p4%3D",
		"chars": "mnyePObBVIzHBtPFWXlawL",
		"ct":    "6ea17e29e28dc2dab9fa7ac973a2832ce4110ef7",
		"ha":    "sha1",
		"tn":    "__jsl_clearance",
		"vt":    "3600",
		"wt":    "1500",
	}

	ha := item["ha"]
	ct := item["ct"]
	tn := item["tn"]
	bts := item["bts"]
	chars := item["chars"]

	cookie, err := GetJsCookie(chars, bts, ct, ha, tn)
	if err != nil {
		return
	}

	g.Dump(cookie)
}

func TestHeFeiFangJiaRelease(t *testing.T) {

	releaseTime, releaseSource, err := GetHeFeiFangJiaRecordReleaseInfo()
	if err == nil {
		fmt.Println(releaseTime, releaseSource)
	}

}
