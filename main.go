package main

import (
	"os"

	"github.com/previousnext/m8s/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("M8s", "Short lived environments")

	cmd.Server(app)
	cmd.Build(app)
	cmd.Step(app)
	cmd.Notify(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
