package github

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/m8s/cmd/m8s/purge/github/drafts"
)

// Commands initializes the github commands.
func Commands(app *kingpin.CmdClause) {
	cmd := app.Command("github", "Purge environments using Github as metadata")
	drafts.Command(cmd)
}
