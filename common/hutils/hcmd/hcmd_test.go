package hcmd

import (
	"fmt"
	"testing"
)

func TestOpenBrowser(t *testing.T) {
	err := OpenBrowser("http://localhost:8080/")
	if err != nil {
		fmt.Println(err)
	}
}
