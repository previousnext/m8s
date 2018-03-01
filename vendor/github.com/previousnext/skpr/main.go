package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/skpr/cmd"
)

func main() {
	app := kingpin.New("skpr", "Cloud-native production hosting proudly supported by PreviousNext.")

	// Core workflow.
	cmd.Package(app, nil)
	cmd.Deploy(app, nil)
	cmd.Version(app, nil)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
