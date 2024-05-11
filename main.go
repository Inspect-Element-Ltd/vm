/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * main.go
 * ---
 * Last Modified: 11/05/2024 01:26AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

package main

import (
	"fmt"
	"github.com/Inspect-Element-Ltd/vm/vmdetect"
)

func main() {
	vm, vendor, reason := vmdetect.Check()
	if vm {
		fmt.Printf("Detected VM (Vendor: %s, Reason: %s)", vendor, reason)
	}
}
