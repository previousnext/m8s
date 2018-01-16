package main

import (
	"github.com/previousnext/gopher/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const (
	// APICompatibility allows client-server tools to avoid breaking changes.
	APICompatibility = 1
)

//go:generate go run scripts/generate-version.go

func main() {
	app := kingpin.New("Gopher", "Bootstrap a go utility")

	cmd.Version(app, APICompatibility)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
