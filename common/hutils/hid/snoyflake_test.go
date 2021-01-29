package hid

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
)

func TestSonyFlake(t *testing.T) {
	parse, _ := time.Parse(time.RFC3339, "2000-01-01")

	settings := Settings{
		StartTime: parse,
	}

	sf := NewSonyFlake(settings)
	go func() {
		for i := 0; i < 10; i++ {
			id, err := sf.NextID()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(id)
		}
	}()

	decompose := Decompose(339556772367303549)

	fmt.Println(1 / 0.000001)
	g.Dump(decompose)

	select {}
}
