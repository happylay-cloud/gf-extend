package hrand

import (
	"fmt"
	"testing"
	"time"

	"github.com/happylay-cloud/gf-extend/common/hutils/htime"
)

func TestHrand(t *testing.T) {
	now := time.Now()
	htime.Sleep(2)
	elapsed := time.Since(now)

	fmt.Println(elapsed)
	fmt.Println(GetRandInt(100))
}

func TestHrandAndRange(t *testing.T) {
	fmt.Println("测试", GetRandIntAndRange(1000, 9999))
}
