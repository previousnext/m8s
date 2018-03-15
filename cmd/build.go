package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/previousnext/compose"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/previousnext/k8s-black-death/retention"
	"github.com/previousnext/m8s/client"
	"github.com/previousnext/m8s/client/types"
	"github.com/previousnext/m8s/config"
	"github.com/previousnext/m8s/utils/environ"
	"github.com/previousnext/m8s/utils/metadata"
)

type cmdBuild struct {
	Client           string
	Config           string
	Name             string
	Domain           string
	Retention        time.Duration
	Repository       string
	Revision         string
	DockerCompose    string
	Master           string
	KubeConfig       string
	ExtraAnnotations string
}

func (cmd *cmdBuild) run(c *kingpin.ParseContext) error {
	dc, err := compose.Load(cmd.DockerCompose)
	if err != nil {
		return errors.Wrap(err, "failed to load Docker Compose file")
	}

	cfg, err := config.LoadWithDefaults(cmd.Config)
	if err != nil {
		return errors.Wrap(err, "failed to load steps")
	}

	// These are additional environment variables that have been provided outside of this build, with the intent
	// for them to be injected into our running containers.
	//   eg. M8S_ENV_FOO=bar, will inject FOO=bar into the containers.
	envs := environ.Get()

	for name, service := range dc.Services {
		service.Environment = append(service.Environment, envs...)
		dc.Services[name] = service
	}

	annotations, err := getAnnotations(cfg.Retention)
	if err != nil {
		return errors.Wrap(err, "failed to build annotations")
	}

	for key, value := range getExtraAnnotations(cmd.ExtraAnnotations) {
		annotations[key] = value
	}

	cli, err := client.New(cmd.Client, types.ClientParams{
		Master:     cmd.Master,
		KubeConfig: cmd.KubeConfig,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	err = cli.Build(os.Stdout, types.BuildParams{
		Name:          strings.ToLower(cmd.Name),
		Domain:        strings.ToLower(cmd.Domain),
		Annotations:   annotations,
		Repository:    cmd.Repository,
		Revision:      cmd.Revision,
		Config:        cfg,
		DockerCompose: dc,
	})
	if err != nil {
		return errors.Wrap(err, "failed to build environment")
	}

	for _, step := range cfg.Build.Steps {
		err = cli.Step(os.Stdout, types.StepParams{
			Namespace: cfg.Namespace,
			Name:      strings.ToLower(cmd.Name),
			Container: cfg.Build.Container,
			Command:   step,
		})
		if err != nil {
			return errors.Wrap(err, "failed to run step")
		}
	}

	return nil
}

// Build declares the "build" sub command.
func Build(app *kingpin.Application) {
	c := new(cmdBuild)

	cmd := app.Command("build", "Build the environment").Action(c.run)
	cmd.Flag("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
	cmd.Flag("domain", "Domain for this environment").Required().StringVar(&c.Domain)
	cmd.Flag("repository", "Git repository to clone from").Default("").Envar("M8S_REPOSITORY").StringVar(&c.Repository)
	cmd.Flag("revision", "Git revision to checkout during clone").Required().StringVar(&c.Revision)
	cmd.Flag("client", "Client to use for building an environment").Default("k8s").Envar("M8S_CLIENT").StringVar(&c.Client)
	cmd.Flag("config", "Build configuration").Default("m8s.yml").Envar("M8S_CONFIG").StringVar(&c.Config)
	cmd.Flag("docker-compose", "Docker Compose file").Default("docker-compose.yml").Envar("M8S_DOCKER_COMPOSE").StringVar(&c.DockerCompose)
	cmd.Flag("master", "Kubernetes master URL").Default().StringVar(&c.Master)
	cmd.Flag("kubeconfig", "Kubernetes config file").Default("~/.kube/config").StringVar(&c.KubeConfig)
	cmd.Flag("extra-annotations", "Add extra annotations to the environment").StringVar(&c.ExtraAnnotations)
}

func getAnnotations(ret time.Duration) (map[string]string, error) {
	annotations, err := metadata.Annotations(os.Environ())
	if err != nil {
		return annotations, errors.Wrap(err, "failed to get annotations from metadata")
	}

	unix, err := retention.Unix(ret)
	if err != nil {
		return annotations, errors.Wrap(err, "failed to convert to unix timestamp")
	}
	annotations[retention.Annotation] = unix

	// This tells admins where the environment came from.
	annotations["author"] = "m8s"

	return annotations, nil
}

// Helper function to extra additional annotations from the cmd flag "ExtraAnnotations".
func getExtraAnnotations(annotations string) map[string]string {
	list := make(map[string]string)

	for _, value := range strings.Split(annotations, ",") {
		sl := strings.Split(value, "=")

		if len(sl) != 2 {
			continue
		}

		var (
			key = sl[0]
			val = sl[1]
		)

		list[key] = val
	}

	return list
}
