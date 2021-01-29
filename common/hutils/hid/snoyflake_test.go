package hid

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
)

func TestSonyFlake(t *testing.T) {
	parse, _ := time.Parse(time.RFC3339, "2021-01-01")

	settings := Settings{
		StartTime: parse,
	}

	sf := NewSonyFlake(settings)

	// 3395,6479,3034,2449,89
	// 339564902606177149
	// 339564939314791293
	// 339564980184089469
	// 339565231053800317
	// 339565309868049277
	// 339565349785502589
	go func() {
		for i := 0; i < 100000; i++ {
			id, err := sf.NextID()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(id)
			decompose := Decompose(id)
			g.Dump(decompose)
		}
	}()

	select {}
}
