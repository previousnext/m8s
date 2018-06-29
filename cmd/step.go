package cmd

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/previousnext/m8s/client"
	"github.com/previousnext/m8s/config"
	"github.com/previousnext/m8s/client/types"
	"strings"
)

type cmdStep struct {
	Client    string
	Config    string
	Name      string
	Container string
	Command   string
	Master string
	KubeConfig string
}

func (cmd *cmdStep) run(c *kingpin.ParseContext) error {
	cli, err := client.New(cmd.Client, types.ClientParams{
		Master: cmd.Master,
		KubeConfig: cmd.KubeConfig,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	cfg, err := config.Load(cmd.Config)
	if err != nil {
		return errors.Wrap(err, "failed to load steps")
	}

	return cli.Step(os.Stdout, types.StepParams{
		Namespace: cfg.Namespace,
		Name: strings.ToLower(cmd.Name),
		Container: cmd.Container,
		Command: cmd.Command,
	})
}

// Step declares the "step" sub command.
func Step(app *kingpin.Application) {
	c := new(cmdStep)

	cmd := app.Command("step", "Step to run against the environment").Action(c.run)
	cmd.Flag("config", "Build configuration").Default("m8s.yml").Envar("M8S_CONFIG").StringVar(&c.Config)
	cmd.Flag("client", "Client to use for building an environment").Default("k8s").Envar("M8S_CLIENT").StringVar(&c.Client)
	cmd.Flag("master", "Kubernetes master URL").Default().StringVar(&c.Master)
	cmd.Flag("kubeconfig", "Kubernetes config file").Default("~/.kube/config").StringVar(&c.KubeConfig)
	cmd.Arg("name", "Name of the environment").Required().StringVar(&c.Name)
	cmd.Arg("container", "Container to execute the step inside").Required().StringVar(&c.Container)
	cmd.Arg("command", "Command to execute").Required().StringVar(&c.Command)
}
