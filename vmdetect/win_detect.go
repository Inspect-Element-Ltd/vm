/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * win_detect.go
 * ---
 * Last Modified: 11/05/2024 01:15AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

//go:build windows
 
package vmdetect

import (
	"github.com/Inspect-Element-Ltd/vm/internal/check"
)

func detectVM() (bool, string, string) {
	if vm, vendor, why := check.Registry(); vm {
		return vm, vendor, why
	}

	if vm, vendor, why := check.FileSystem(); vm {
		return vm, vendor, why
	}

	return false, "", ""
}
