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
	cmd.Step(app)
	cmd.Notify(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
