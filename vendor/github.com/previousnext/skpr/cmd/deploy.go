package cmd

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/skpr/cmd/client"
)

type cmdDeploy struct {
	client client.Client
	params client.DeployParams
}

func (d *cmdDeploy) run(c *kingpin.ParseContext) error {
	return d.client.Deploy(os.Stdout, d.params)
}

// Deploy declares the "deploy" sub command.
func Deploy(app *kingpin.Application, cli client.Client) {
	d := new(cmdDeploy)
	d.client = cli

	command := app.Command("deploy", "Triggers deployment of a packaged version to the specified environment.").Action(d.run)
	command.Flag("wait", "Waits for the deployment to complete before exiting.").Default("true").BoolVar(&d.params.Wait)
	command.Flag("dry-run", "Return to the user what would happen if dry-run is omitted.").BoolVar(&d.params.DryRun)
	command.Arg("env", "Environment as configured. Usually one of: prod/staging/dev.").Required().StringVar(&d.params.Environment)
	command.Arg("version", "Version to deploy. Corresponds to the tag on a packaged container image.").Required().StringVar(&d.params.Version)
}
