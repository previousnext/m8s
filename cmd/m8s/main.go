package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/m8s/cmd/m8s/purge"
)

func main() {
	app := kingpin.New("m8s", "Flexible, temporary environments for rapid iteration")

	purge.Commands(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
