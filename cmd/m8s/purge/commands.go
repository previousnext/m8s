package purge

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/m8s/cmd/m8s/purge/github"
)

// Commands initializes the github commands.
func Commands(app *kingpin.Application) {
	cmd := app.Command("purge", "Purge environments")
	github.Commands(cmd)
}
