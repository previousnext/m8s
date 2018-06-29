package cmd

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/gopher/pkg/version"
)

var (
	// GitVersion is overridden at build time.
	GitVersion string
	// GitCommit is overridden at build time.
	GitCommit string
)

type cmdVersion struct{}

func (v *cmdVersion) run(c *kingpin.ParseContext) error {
	return version.Print(os.Stdout, version.PrintParams{
		Version: GitVersion,
		Commit:  GitCommit,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	})
}

// Version declares the "version" sub command.
func Version(app *kingpin.Application) {
	v := new(cmdVersion)
	app.Command("version", fmt.Sprintf("Prints %s version", app.Name)).Action(v.run)
}
