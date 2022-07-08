package hjs

import (
	"github.com/gogf/gf/frame/g"
	"testing"
)

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
