package main

import (
	"os"

	"github.com/previousnext/m8s/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	app := kingpin.New("M8s", "Short lived environments")

	cmd.Build(app)
	cmd.Step(app)
	cmd.API(app)
	cmd.Version(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
