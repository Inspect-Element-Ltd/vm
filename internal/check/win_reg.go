/*
 * Copyright 2024, Inspect Element Ltd <https://echo.ac>.
 *
 * See LICENCE <https://github.com/Inspect-Element-Ltd/vm/blob/master/LICENCE>
 *
 * win_reg.go
 * ---
 * Last Modified: 11/05/2024 01:15AM (BST)
 * Modified By: Gianluca Oliva <hello@gian.sh>
 */

//go:build windows

package check

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	// https://github.com/CheckPointSW/Evasions/blob/master/_src/Evasions/techniques/registry.md

	hyperVKeys = []string{
		`HKLM\SOFTWARE\Microsoft\Hyper-V`,
		`HKLM\SOFTWARE\Microsoft\VirtualMachine`,
		`HKLM\SOFTWARE\Microsoft\Virtual Machine\Guest\Parameters`,
		// False flags? They show up on any system with Hyper-V enabled?
		//`HKLM\SYSTEM\ControlSet001\Services\vmicheartbeat`,
		//`HKLM\SYSTEM\ControlSet001\Services\vmicvss`,
		//`HKLM\SYSTEM\ControlSet001\Services\vmicshutdown`,
		//`HKLM\SYSTEM\ControlSet001\Services\vmicexchange`,
	}

	parallelsKeys = []string{
		`HKLM\SYSTEM\CurrentControlSet\Enum\PCI\VEN_1AB8*`,
	}

	virtualBoxKeys = []string{
		`HKLM\SYSTEM\CurrentControlSet\Enum\PCI\VEN_80EE*`,
		`HKLM\HARDWARE\ACPI\DSDT\VBOX__`,
		`HKLM\HARDWARE\ACPI\FADT\VBOX__`,
		`HKLM\HARDWARE\ACPI\RSDT\VBOX__`,
		`HKLM\SOFTWARE\Oracle\VirtualBox Guest Additions`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxGuest`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxMouse`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxService`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxSF`,
		`HKLM\SYSTEM\ControlSet001\Services\VBoxVideo`,
	}

	vmwareKeys = []string{
		`HKLM\SYSTEM\CurrentControlSet\Enum\PCI\VEN_15AD*`,
		`HKCU\SOFTWARE\VMware, Inc.\VMware Tools`,
		`HKLM\SOFTWARE\VMware, Inc.\VMware Tools`,
		`HKLM\SYSTEM\ControlSet001\Services\vmdebug`,
		`HKLM\SYSTEM\ControlSet001\Services\vmmouse`,
		`HKLM\SYSTEM\ControlSet001\Services\VMTools`,
		`HKLM\SYSTEM\ControlSet001\Services\VMMEMCTL`,
		`HKLM\SYSTEM\ControlSet001\Services\vmware`,
		`HKLM\SYSTEM\ControlSet001\Services\vmci`,
		`HKLM\SYSTEM\ControlSet001\Services\vmx86`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\CdRomNECVMWar_VMware_IDE_CD*`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\CdRomNECVMWar_VMware_SATA_CD*`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\DiskVMware_Virtual_IDE_Hard_Drive*`,
		`HKLM\SYSTEM\CurrentControlSet\Enum\IDE\DiskVMware_Virtual_SATA_Hard_Drive*`,
	}

	wineKeys = []string{
		`HKCU\SOFTWARE\Wine`,
		`HKLM\SOFTWARE\Wine`,
	}

	xenKeys = []string{
		`HKLM\HARDWARE\ACPI\DSDT\xen`,
		`HKLM\HARDWARE\ACPI\FADT\xen`,
		`HKLM\HARDWARE\ACPI\RSDT\xen`,
		`HKLM\SYSTEM\ControlSet001\Services\xenevtchn`,
		`HKLM\SYSTEM\ControlSet001\Services\xennet`,
		`HKLM\SYSTEM\ControlSet001\Services\xennet6`,
		`HKLM\SYSTEM\ControlSet001\Services\xensvc`,
		`HKLM\SYSTEM\ControlSet001\Services\xenvdb`,
	}

	vendorValues = map[string]map[string][]string{
		"Generic": {
			`HKLM\HARDWARE\Description\System\SystemBiosDate`:         []string{"06/23/99"},
			`HKLM\HARDWARE\Description\System\BIOS\SystemProductName`: []string{"A M I"},
		},
		"Parallels": {
			`HKLM\HARDWARE\Description\System\SystemBiosVersion`: []string{"PARALLELS"},
			`HKLM\HARDWARE\Description\System\VideoBiosVersion`:  []string{"PARALLELS"},
		},
		"QEMU": {
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 0\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: []string{"QEMU"},
			`HKLM\HARDWARE\Description\System\SystemBiosVersion`:                                           []string{"QEMU"},
			`HKLM\HARDWARE\Description\System\VideoBiosVersion`:                                            []string{"QEMU"},
			`HKLM\HARDWARE\Description\System\BIOS\SystemManufacturer`:                                     []string{"QEMU"},
		},
		"VirtualBox": {
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 0\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: []string{"VBOX"},
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 1\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: []string{"VBOX"},
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 2\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`: []string{"VBOX"},
			`HKLM\HARDWARE\Description\System\SystemBiosVersion`:                                           []string{"VBOX"},
			`HKLM\HARDWARE\Description\System\VideoBiosVersion`:                                            []string{"VIRTUALBOX"},
			`HKLM\HARDWARE\Description\System\BIOS\SystemProductName`:                                      []string{"VIRTUAL"},
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\DeviceDesc`:                                      []string{"VBOX"},
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\FriendlyName`:                                    []string{"VBOX"},
			`HKLM\SYSTEM\ControlSet002\Services\Disk\Enum\DeviceDesc`:                                      []string{"VBOX"},
			`HKLM\SYSTEM\ControlSet002\Services\Disk\Enum\FriendlyName`:                                    []string{"VBOX"},
			`HKLM\SYSTEM\ControlSet003\Services\Disk\Enum\DeviceDesc`:                                      []string{"VBOX"},
			`HKLM\SYSTEM\ControlSet003\Services\Disk\Enum\FriendlyName`:                                    []string{"VBOX"},
			`HKLM\SYSTEM\CurrentControlSet\Control\SystemInformation\SystemProductName`:                    []string{"VIRTUAL", "VIRTUALBOX"},
		},
		"VMware": {
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 0\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`:                    []string{"VMWARE"},
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 1\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`:                    []string{"VMWARE"},
			`HKLM\HARDWARE\DEVICEMAP\Scsi\Scsi Port 2\Scsi Bus 0\Target Id 0\Logical Unit Id 0\Identifier`:                    []string{"VMWARE"},
			`HKLM\HARDWARE\Description\System\SystemBiosVersion`:                                                              []string{"VMWARE", "INTEL - 6040000"},
			`HKLM\HARDWARE\Description\System\VideoBiosVersion`:                                                               []string{"VMWARE"},
			`HKLM\HARDWARE\Description\System\BIOS\SystemProductName`:                                                         []string{"VMware"},
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\0`:                                                                  []string{"VMware"},
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\1`:                                                                  []string{"VMware"},
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\DeviceDesc`:                                                         []string{"VMware"},
			`HKLM\SYSTEM\ControlSet001\Services\Disk\Enum\FriendlyName`:                                                       []string{"VMware"},
			`HKLM\SYSTEM\ControlSet002\Services\Disk\Enum\DeviceDesc`:                                                         []string{"VMware"},
			`HKLM\SYSTEM\ControlSet002\Services\Disk\Enum\FriendlyName`:                                                       []string{"VMware"},
			`HKLM\SYSTEM\ControlSet003\Services\Disk\Enum\DeviceDesc`:                                                         []string{"VMware"},
			`HKLM\SYSTEM\ControlSet003\Services\Disk\Enum\FriendlyName`:                                                       []string{"VMware"},
			`HKCR\Installer\Products\ProductName`:                                                                             []string{"vmware tools"},
			`HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\DisplayName`:                                            []string{"vmware tools"},
			`HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\DisplayName`:                                            []string{"vmware tools"},
			`HKLM\SYSTEM\ControlSet001\Control\Class\{4D36E968-E325-11CE-BFC1-08002BE10318}\0000\CoInstallers32`:              []string{"*vmx*"},
			`HKLM\SYSTEM\ControlSet001\Control\Class\{4D36E968-E325-11CE-BFC1-08002BE10318}\0000\DriverDesc`:                  []string{"VMWare*"},
			`HKLM\SYSTEM\ControlSet001\Control\Class\{4D36E968-E325-11CE-BFC1-08002BE10318}\0000\InfSection`:                  []string{"vmx*"},
			`HKLM\SYSTEM\ControlSet001\Control\Class\{4D36E968-E325-11CE-BFC1-08002BE10318}\0000\ProviderName`:                []string{"VMware*"},
			`HKLM\SYSTEM\ControlSet001\Control\Class\{4D36E968-E325-11CE-BFC1-08002BE10318}\0000\Settings\Device Description`: []string{"VMware*"},
			`HKLM\SYSTEM\CurrentControlSet\Control\SystemInformation\SystemProductName`:                                       []string{"VMWARE"},
			`HKLM\SYSTEM\CurrentControlSet\Control\Video\{GUID}\Video\Service`:                                                []string{"vm3dmp", "vmx_svga"},
			`HKLM\SYSTEM\CurrentControlSet\Control\Video\{GUID}\0000\Device Description`:                                      []string{"VMware SVGA*"},
		},
		"Xen": {
			`HKLM\HARDWARE\Description\System\BIOS\SystemProductName`: []string{"Xen"},
		},
	}
)

// https://github.com/josheyr/VM-Detection/blob/74d0e106ec7dd0f6cce49c4fc0e9ba682d4dc657/vmdetect/windows.go#L15C1-L42C2
func extractKeyTypeFrom(registryKey string) (registry.Key, string, error) {
	firstSeparatorIndex := strings.Index(registryKey, string(os.PathSeparator))
	if firstSeparatorIndex == -1 {
		return 0, "", errors.New("invalid registry key")
	}

	keyTypeStr := registryKey[:firstSeparatorIndex]
	keyPath := registryKey[firstSeparatorIndex+1:]

	var keyType registry.Key
	switch keyTypeStr {
	case "HKLM":
		keyType = registry.LOCAL_MACHINE
		break
	case "HKCR":
		keyType = registry.CLASSES_ROOT
		break
	case "HKCU":
		keyType = registry.CURRENT_USER
		break
	case "HKU":
		keyType = registry.USERS
		break
	case "HKCC":
		keyType = registry.CURRENT_CONFIG
		break
	default:
		return keyType, "", errors.New(fmt.Sprintf("Invalid keytype (%v)", keyTypeStr))
	}

	return keyType, keyPath, nil
}

// https://github.com/josheyr/VM-Detection/blob/74d0e106ec7dd0f6cce49c4fc0e9ba682d4dc657/vmdetect/windows.go#L44C1-L71C2
func doesRegistryKeyContain(registryKey string, expectedSubString string) bool {
	keyType, keyPath, err := extractKeyTypeFrom(registryKey)

	if err != nil {
		return false
	}

	keyPath, keyName := filepath.Split(keyPath)

	keyHandle, err := registry.OpenKey(keyType, keyPath, registry.QUERY_VALUE)

	if err != nil {
		return false
	}
	defer keyHandle.Close()

	valueFound, _, err := keyHandle.GetStringValue(keyName)
	if err != nil {
		return false
	}

	return strings.Contains(valueFound, expectedSubString)
}

// https://github.com/josheyr/VM-Detection/blob/74d0e106ec7dd0f6cce49c4fc0e9ba682d4dc657/vmdetect/windows.go#L73
func doesRegistryKeyExist(registryKey string) bool {
	subkeyPrefix := ""

	// Handle trailing wildcard
	if registryKey[len(registryKey)-1:] == "*" {
		registryKey, subkeyPrefix = path.Split(registryKey)
		subkeyPrefix = subkeyPrefix[:len(subkeyPrefix)-1] // remove *
	}

	keyType, keyPath, err := extractKeyTypeFrom(registryKey)

	if err != nil {
		return false
	}

	keyHandle, err := registry.OpenKey(keyType, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}

	defer keyHandle.Close()

	// If a wildcard has been specified...
	if subkeyPrefix != "" {
		// ... we look for sub-keys to see if one exists
		subKeys, err := keyHandle.ReadSubKeyNames(0xFFFF)

		if err != nil {
			return false
		}

		for _, subKeyName := range subKeys {
			if strings.HasPrefix(subKeyName, subkeyPrefix) {
				return true
			}
		}

		return false
	}

	return true
}

func Registry() (bool, string, string) {
	for _, key := range hyperVKeys {
		if doesRegistryKeyExist(key) {
			return true, "Hyper-V", fmt.Sprintf("%s found in Registry", key)
		}
	}

	for _, key := range parallelsKeys {
		if doesRegistryKeyExist(key) {
			return true, "Parallels", fmt.Sprintf("%s found in Registry", key)
		}
	}

	for _, key := range virtualBoxKeys {
		if doesRegistryKeyExist(key) {
			return true, "VirtualBox", fmt.Sprintf("%s found in Registry", key)
		}
	}

	for _, key := range vmwareKeys {
		if doesRegistryKeyExist(key) {
			return true, "VMware", fmt.Sprintf("%s found in Registry", key)
		}
	}

	for _, key := range wineKeys {
		if doesRegistryKeyExist(key) {
			return true, "Wine", fmt.Sprintf("%s found in Registry", key)
		}
	}

	for _, key := range xenKeys {
		if doesRegistryKeyExist(key) {
			return true, "Xen", fmt.Sprintf("%s found in Registry", key)
		}
	}

	for vendor, registryValues := range vendorValues {
		for registryPath, values := range registryValues {
			for _, value := range values {
				if doesRegistryKeyContain(registryPath, value) {
					return true, vendor, fmt.Sprintf("Registry Path %s contains %s", registryPath, value)
				}
			}
		}
	}

	return false, "", ""
}
