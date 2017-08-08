package main

import (
	"os"

	"github.com/previousnext/pr/cli/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("PR Environment", "PreviousNext Pull Request environments")

	// Setup all the subcommands.
	cmd.Build(app)
	cmd.DockerCfgGet(app)
	cmd.DockerCfgSet(app)
	cmd.SSHGet(app)
	cmd.SSHSet(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
