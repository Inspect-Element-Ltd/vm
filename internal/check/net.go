/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * net.go
 * ---
 * Last Modified: 11/05/2024 12:22AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

package check

import (
	"net"
	"strings"
)

var (
	ouiByVendor = map[string][]string{
		"Parallels": {
			"00:1C:42",
		},
		"VMware": {
			"00:05:69", "00:0C:29", "00:1C:14", "00:50:56",
		},
	}
)

func MACAddress() (bool, string, string) {

	if ifaces, err := net.Interfaces(); err == nil && ifaces != nil {
		for _, iface := range ifaces {
			for vendor, ouis := range ouiByVendor {
				for _, oui := range ouis {
					if strings.HasPrefix(iface.HardwareAddr.String(), oui) {
						return true, vendor, "OUI Prefix matches " + vendor
					}
				}
			}
		}
	}

	return false, "", ""
}
