package utils

import (
	"fmt"
)

func RunWithRecovery(function func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("RECOVERED ERROR %v", r)
		}
	}()
	function()
}
