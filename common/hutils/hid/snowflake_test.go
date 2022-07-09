package hid

import (
	"github.com/happylay-cloud/gf-extend/common/hutils/htime"

	"fmt"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	nextId, _ := GetSnowWorker().NextId()
	fmt.Println(nextId)

	nextIdStr := GetSnowWorker().NextIdStr()
	fmt.Println(nextIdStr)

	timestamp := GetSnowWorker().IdToTimestamp(nextId)
	fmt.Println(timestamp)
	fmt.Println(htime.TimestampFormat(timestamp))
}
