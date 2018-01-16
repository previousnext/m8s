package cmd

import (
	"fmt"
	"runtime"

	"github.com/ryanuber/columnize"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/gopher/version"
)

type cmdVersion struct {
	APICompatibility int
	BuildDate        string
	BuildVersion     string
	GOARCH           string
	GOOS             string
}

func (cmd *cmdVersion) run(c *kingpin.ParseContext) error {
	fmt.Println(renderVersionOutput(cmd))
	return nil
}

// Version declares the "version" sub command.
func Version(app *kingpin.Application, apiCompatibility int) {
	cmd := cmdVersion{
		APICompatibility: apiCompatibility,
		BuildDate:        version.BuildDate,
		BuildVersion:     version.BuildVersion,
		GOARCH:           runtime.GOARCH,
		GOOS:             runtime.GOOS,
	}

	app.Command("version", fmt.Sprintf("Prints %s version", app.Name)).Action(cmd.run)
}

// RenderVersionOutput is responsible for producing the rendered version info string.
func renderVersionOutput(cmd *cmdVersion) string {
	output := []string{
		fmt.Sprintf("Version | %s", cmd.BuildVersion),
		fmt.Sprintf("Date | %s", cmd.BuildDate),
		fmt.Sprintf("API | v%d", cmd.APICompatibility),
		fmt.Sprintf("OS | %s", cmd.GOOS),
		fmt.Sprintf("Arch | %s", cmd.GOARCH),
	}
	return columnize.SimpleFormat(output)
}
