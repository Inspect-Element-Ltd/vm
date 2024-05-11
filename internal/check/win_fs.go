/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * win_fs.go
 * ---
 * Last Modified: 11/05/2024 01:15AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

//go:build windows

package check

import (
	"fmt"
	"os"
)

var (
	filesByVendor = map[string][]string{
		"Parallels": {
			`c:\windows\system32\drivers\prleth.sys`,
			`c:\windows\system32\drivers\prlfs.sys`,
			`c:\windows\system32\drivers\prlmouse.sys`,
			`c:\windows\system32\drivers\prlvideo.sys`,
			`c:\windows\system32\drivers\prltime.sys`,
			`c:\windows\system32\drivers\prl_pv32.sys`,
			`c:\windows\system32\drivers\prl_paravirt_32.sys`,
		},
		"VirtualBox": {
			`c:\windows\system32\drivers\VBoxMouse.sys`,
			`c:\windows\system32\drivers\VBoxGuest.sys`,
			`c:\windows\system32\drivers\VBoxSF.sys`,
			`c:\windows\system32\drivers\VBoxVideo.sys`,
			`c:\windows\system32\vboxdisp.dll`,
			`c:\windows\system32\vboxhook.dll`,
			`c:\windows\system32\vboxmrxnp.dll`,
			`c:\windows\system32\vboxogl.dll`,
			`c:\windows\system32\vboxoglarrayspu.dll`,
			`c:\windows\system32\vboxoglcrutil.dll`,
			`c:\windows\system32\vboxoglerrorspu.dll`,
			`c:\windows\system32\vboxoglfeedbackspu.dll`,
			`c:\windows\system32\vboxoglpackspu.dll`,
			`c:\windows\system32\vboxoglpassthroughspu.dll`,
			`c:\windows\system32\vboxservice.exe`,
			`c:\windows\system32\vboxtray.exe`,
			`c:\windows\system32\VBoxControl.exe`,
		},
		"VirtualPC": {
			`c:\windows\system32\drivers\vmsrvc.sys`,
			`c:\windows\system32\drivers\vpc-s3.sys`,
		},
		"VMware": {
			`c:\windows\system32\drivers\vmmouse.sys`,
			`c:\windows\system32\drivers\vmnet.sys`,
			`c:\windows\system32\drivers\vmxnet.sys`,
			`c:\windows\system32\drivers\vmhgfs.sys`,
			`c:\windows\system32\drivers\vmx86.sys`,
			`c:\windows\system32\drivers\hgfs.sys`,
		},
	}
)

func FileSystem() (bool, string, string) {
	for vendor, files := range filesByVendor {
		for _, file := range files {
			if _, err := os.Stat(file); err == nil {
				return true, vendor, fmt.Sprintf("%s file exists", file)
			}
		}
	}

	return false, "", ""
}
