package hid

import (
	"fmt"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	nextId, _ := GetSnowWorker().NextId()
	fmt.Println(nextId)

	nextIdStr := GetSnowWorker().NextIdStr()
	fmt.Println(nextIdStr)

	timestamp := GetSnowWorker().IdToTimestamp(1545749517502894080)
	fmt.Println(timestamp)

}
