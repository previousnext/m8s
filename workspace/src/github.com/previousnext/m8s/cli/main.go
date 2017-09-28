package main

import (
	"os"

	"github.com/previousnext/m8s/cli/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("M8s", "Short lived environments")

	// Setup all the subcommands.
	cmd.Build(app)
	cmd.Exec(app)
	cmd.List(app)
	cmd.Describe(app)
	cmd.Slack(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
