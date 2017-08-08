package cmd

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdSSHSet struct{}

func (cmd *cmdSSHSet) run(c *kingpin.ParseContext) error {
	return nil
}

func SSHSet(app *kingpin.Application) {
	c := new(cmdSSHSet)

	app.Command("ssh-set", "Updates the SSH secret configuration").Action(c.run)
}
