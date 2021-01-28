package huuid

import (
	"fmt"
	"testing"
)

func TestHUUiID(t *testing.T) {
	uuid := UUID()
	fmt.Println(uuid)
}
