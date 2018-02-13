package cmd

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/skpr/cmd/client"
)

type cmdVersion struct {
	client client.Client
	params client.VersionParams
}

func (v *cmdVersion) run(c *kingpin.ParseContext) error {
	return v.client.Version(os.Stdout, v.params)
}

// Version declares the "version" sub command.
func Version(app *kingpin.Application, client client.Client) {
	v := new(cmdVersion)
	app.Command("version", fmt.Sprintf("Prints %s version", app.Name)).Action(v.run)
}
