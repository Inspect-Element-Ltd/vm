/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * mac_reg.go
 * ---
 * Last Modified: 11/05/2024 01:15AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

//go:build darwin

package check

import (
	"fmt"
	"github.com/Inspect-Element-Ltd/vm/internal/util"
	"strings"
)

var vendors = []string{
	"VMWare",
	"Oracle",
	"Parallels",
	"VirtualBox",
}

func Registry() (bool, string, string) {

	// Most VM software like VMWare, VirtualBox, etc. will have a serial number of "0".
	serialNumber, err := util.InvokeCMD("bash", "-c", "ioreg -rd1 -c IOPlatformExpertDevice | grep 'IOPlatformSerialNumber'")
	if err == nil {
		serialNumber = strings.TrimSpace(serialNumber)
		serialNumber = strings.ReplaceAll(serialNumber, `"`, "")
		if len(strings.Split(serialNumber, " = ")) == 2 {
			serialNumber = strings.TrimSpace(strings.Split(serialNumber, " = ")[1])
			if serialNumber == "0" {
				return true, "Generic", "Serial Number is 0"
			}
		}
	}

	// If the board manufacturer doesn't contain "Apple" then it's likely a VM.
	manufacturer, err := util.InvokeCMD("bash", "-c", "ioreg -rd1 -c IOPlatformExpertDevice | grep 'manufacturer'")
	if err == nil {
		manufacturer = strings.TrimSpace(manufacturer)
		manufacturer = strings.ReplaceAll(manufacturer, `"`, "")
		if len(strings.Split(manufacturer, " = ")) == 2 {
			manufacturer = strings.TrimSpace(strings.Split(manufacturer, " = ")[1])
			manufacturer = strings.ReplaceAll(manufacturer, `<`, "")
			manufacturer = strings.ReplaceAll(manufacturer, `>`, "")
			if !strings.Contains(manufacturer, "Apple") {
				return true, "Generic", fmt.Sprintf("Manufacturer is %s not Apple Inc.", manufacturer)
			}
		}
	}

	vendorNames, err := util.InvokeCMD("bash", "-c", "ioreg -l | grep -e Manufacturer -e 'Vendor Name'")
	if err == nil {
		for _, vendorName := range strings.Split(vendorNames, "\n") {
			vendorName = strings.ReplaceAll(vendorName, `|`, "")
			vendorName = strings.ReplaceAll(vendorName, `"`, "")
			vendorName = strings.TrimSpace(vendorName)
			if vendorName == "" {
				continue
			}

			if len(strings.Split(vendorName, " = ")) == 2 {
				vendorName = strings.TrimSpace(strings.Split(vendorName, " = ")[1])
				for _, vendor := range vendors {
					if strings.Contains(strings.ToLower(vendorName), strings.ToLower(vendor)) {
						return true, vendor, fmt.Sprintf("Vendor Name contains %s", vendor)
					}
				}
			}
		}
	}

	return false, "", ""
}
