# VM Detect
Simple Go library to detect virtual machines.<br>

## Usage

```bash
go get github.com/Inspect-Element-Ltd/vm
```

```go
import (
    "fmt"
    "github.com/Inspect-Element-Ltd/vm/vmdetect"
)

vm, vendor, reason := vmdetect.Check()
if vm {
    fmt.Printf("Detected VM (Vendor: %s, Reason: %s)", vendor, reason)
}
```

### TODO
- [ ] Linux support
- [ ] Clean up the horrible code in `mac_reg.go`

### Credits
Heavily inspired by [VM-Detection by ShellCode33](https://github.com/ShellCode33/VM-Detection).
Most, if not all, of the Windows code is from VM-Detection but with some detections removed, our use-case doesn't want to flag a system with Sandbox tools installed as a VM. 