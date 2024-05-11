/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * detect.go
 * ---
 * Last Modified: 11/05/2024 01:13AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

package vmdetect

import (
	"github.com/Inspect-Element-Ltd/vm/internal/check"
	"github.com/klauspost/cpuid/v2"
)

// IsVM attempts to figure out if the current system is a virtual machine.
//
// Calls Check but only returns a boolean value,
// the vendor and detection explanation are ignored
func IsVM() bool {
	vm, _, _ := Check()
	return vm
}

// Check attempts to figure out if the current system is a virtual machine.
//
// If a VM is detected the Vendor and why it was detected is also returned,
// these values will be empty if the machine is not detected as being virtualised.
func Check() (bool, string, string) {
	switch cpuid.CPU.VendorID {
	case cpuid.MSVM, cpuid.KVM, cpuid.VMware, cpuid.XenHVM, cpuid.Bhyve:
		return true, cpuid.CPU.VendorString, "CPUID"
	default:
		break
	}

	if vm, vendor, why := check.MACAddress(); vm {
		return vm, vendor, why
	}

	return detectVM()
}
