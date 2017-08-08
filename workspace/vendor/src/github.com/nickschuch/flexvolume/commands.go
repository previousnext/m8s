package flexvolume

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli"
)

func Commands(fv FlexVolume) []cli.Command {
	return []cli.Command{
		{
			Name:  "init",
			Usage: "Initialize the driver",
			Action: func(c *cli.Context) error {
				return handle(fv.Init())
			},
		},
		{
			Name:  "attach",
			Usage: "Attach the volume",
			Action: func(c *cli.Context) error {
				var opts map[string]string

				if err := json.Unmarshal([]byte(c.Args().Get(1)), &opts); err != nil {
					return err
				}

				return handle(fv.Attach(opts))
			},
		},
		{
			Name:  "detach",
			Usage: "Detach the volume",
			Action: func(c *cli.Context) error {
				return handle(fv.Detach(c.Args().Get(1)))
			},
		},
		{
			Name:  "mount",
			Usage: "Mount the volume",
			Action: func(c *cli.Context) error {
				var opts map[string]string

				if err := json.Unmarshal([]byte(c.Args().Get(1)), &opts); err != nil {
					return err
				}

				return handle(fv.Mount(c.Args().Get(1), c.Args().Get(2), opts))
			},
		},
		{
			Name:  "umount",
			Usage: "Mount the volume",
			Action: func(c *cli.Context) error {
				return handle(fv.Unmount(c.Args().Get(1)))
			},
		},
	}
}

// The following handles:
//   * Output of the Response object.
//   * Sets an error so we can bubble up an error code.
func handle(resp Response) error {
	// Format the output as JSON.
	output, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
