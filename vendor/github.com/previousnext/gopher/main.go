package main

import (
	"os"

	"github.com/previousnext/gopher/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("Gopher", "Bootstrap a go utility")

	cmd.Version(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
