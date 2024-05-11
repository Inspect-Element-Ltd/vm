/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * mac_sysctl.go
 * ---
 * Last Modified: 11/05/2024 01:15AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

//go:build darwin

package check

import (
	"github.com/Inspect-Element-Ltd/vm/internal/util"
	"strconv"
	"strings"
)

// HardwareModel checks the hw.model is missing the word 'Mac'.
func HardwareModel() (bool, string, string) {
	hwModel, err := util.InvokeCMD("sysctl", "-n", "hw.model")
	if err != nil {
		return false, "", ""
	}

	if !strings.Contains(hwModel, "Mac") {
		return true, hwModel, "hw.modal doesn't contain 'Mac'"
	}

	return false, "", ""
}

// MemorySize checks the hw.memsize to see if it's less than 4GB.
func MemorySize() (bool, string, string) {
	memSize, err := util.InvokeCMD("sysctl", "-n", "hw.memsize")
	if err != nil {
		return false, "", ""
	}

	memBytes, err := strconv.ParseInt(strings.TrimSpace(memSize), 10, 64)
	if err != nil {
		return false, "", ""
	}

	if memBytes < 4294967296 {
		return true, "Generic", "hw.memsize is less than 4GB"
	}

	return false, "", ""
}
