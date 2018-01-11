package main

import (
	"os"

	"github.com/previousnext/m8s/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

//go:generate go run scripts/generate-version.go

func main() {
	app := kingpin.New("M8s", "Short lived environments")

	// Core workflow.
	cmd.Server(app)
	cmd.Build(app)
	cmd.Step(app)

	// API for the M8s UI.
	cmd.API(app)

	// Utility for installing M8s components on a K8s stack.
	cmd.Install(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
