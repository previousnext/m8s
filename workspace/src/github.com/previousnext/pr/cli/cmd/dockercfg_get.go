package cmd

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdDockerCfgGet struct{}

func (cmd *cmdDockerCfgGet) run(c *kingpin.ParseContext) error {
	return nil
}

func DockerCfgGet(app *kingpin.Application) {
	c := new(cmdDockerCfgGet)

	app.Command("dockercfg-get", "Returns the Docker secret configuration").Action(c.run)
}
