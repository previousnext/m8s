Kubernetes: FlexVolume
======================

Simple wrapper library for quick FlexVolume prototypes.

## Example

This is a simple implementation which creates and deletes folders, nothing crazy.

The following example should demonstate:

* Setting up the cli
* Implementing the flexvolume interface
* Returning the device output required

```go
package main

import (
    "os"
    "os/exec"

    "github.com/nickschuch/flexvolume"
    "github.com/urfave/cli"
)

func main() {
    app := cli.NewApp()
    app.Commands = flexvolume.Commands(Mock{})
    app.Run(os.Args)
}

type Mock struct {}

func (m Mock) Init() flexvolume.Response {
    return flexvolume.Response{
        Status:  flexvolume.StatusSuccess,
        Message: "Mock is available",
    }
}

func (m Mock) Attach(options map[string]string) flexvolume.Response {
    return flexvolume.Response{
        Status:  flexvolume.StatusSuccess,
        Message: "Successfully attached the mock volume",
    }
}

func (m Mock) Detach(device string) flexvolume.Response {
    return flexvolume.Response{
        Status:  flexvolume.StatusSuccess,
        Message: "Successfully detached the mock volume",
    }
}

func (m Mock) Mount(target, device string, options map[string]string) flexvolume.Response {
    device := "/dev/sdb1"
    
    err := os.MkdirAll(target, 0755)
    if err != nil {
        return flexvolume.Response{
            Status:  flexvolume.StatusFailure,
            Message: err.Error(),
            Device:  device,
        }
    }

    return flexvolume.Response{
        Status:  flexvolume.StatusSuccess,
        Message: "Successfully mounted the mock volume",
        Device:  device,
    }
}

func (m Mock) Unmount(mount string) flexvolume.Response {
    err := os.RemoveAll(mount)
    if err != nil {
        return flexvolume.Response{
            Status:  flexvolume.StatusFailure,
            Message: err.Error(),
        }
    }

    return flexvolume.Response{
        Status:  flexvolume.StatusSuccess,
        Message: "Successfully unmounted the mock volume",
    }
}
```
