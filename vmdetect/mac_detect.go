/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * mac_detect.go
 * ---
 * Last Modified: 11/05/2024 01:15AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

//go:build darwin

package vmdetect

import (
	"github.com/Inspect-Element-Ltd/vm/internal/check"
	"github.com/Inspect-Element-Ltd/vm/internal/util"
	"strings"
)

func SIPDisabled() bool {
	sip, err := util.InvokeCMD("bash", "-c", "csrutil status")
	if err != nil {
		return false
	}

	return strings.TrimSpace(sip) != "System Integrity Protection status: enabled."
}

func detectVM() (bool, string, string) {

	if vm, vendor, why := check.HardwareModel(); vm {
		return vm, vendor, why
	}

	if vm, vendor, why := check.MemorySize(); vm {
		return vm, vendor, why
	}

	if vm, vendor, why := check.Registry(); vm {
		return vm, vendor, why
	}

	return false, "", ""
}
