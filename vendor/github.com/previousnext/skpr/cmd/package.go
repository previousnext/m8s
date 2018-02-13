package cmd

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/skpr/cmd/client"
)

type cmdPackage struct {
	client client.Client
	params client.PackageParams
}

func (p *cmdPackage) run(c *kingpin.ParseContext) error {
	return p.client.Package(os.Stdout, p.params)
}

// Package declares the "package" sub command.
func Package(app *kingpin.Application, cli client.Client) {
	p := new(cmdPackage)
	p.client = cli

	command := app.Command("package", "Packages an app container image. The packaging steps are defined in the project `Dockerfile`.").Action(p.run)
	command.Arg("version", "The version to tag the built image.").Required().StringVar(&p.params.Version)
	command.Arg("directory", "The project root directory.").Default(".").StringVar(&p.params.Directory)
}
