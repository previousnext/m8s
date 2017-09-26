package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/gosexy/to"
	"github.com/previousnext/m8s/cli/cmd/environ"
	"github.com/previousnext/m8s/cli/compose"
	pb "github.com/previousnext/m8s/pb"
	"github.com/smallfish/simpleyaml"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdBuild struct {
	API              string
	Token            string
	Name             string
	Domains          string
	BasicAuthUser    string
	BasicAuthPass    string
	GitRepository    string
	GitRevision      string
	DockerCompose    string
	DockerRepository string
	ExecFile         string
	ExecStep         string
	ExecInside       string
	Timeout          time.Duration
}

func (cmd *cmdBuild) run(c *kingpin.ParseContext) error {
	// Load the Docker Compose file, we are going to use alot of its
	// configuration for this build.
	dc, err := compose.Load(cmd.DockerCompose)
	if err != nil {
		return fmt.Errorf("failed to load Docker Compose file: %s", err)
	}

	// Load the steps required to run the build, these are bespoke steps used
	// for bootstrapping and testing the application.
	steps, err := loadSteps(cmd.ExecFile, cmd.ExecStep)
	if err != nil {
		return fmt.Errorf("failed to load steps: %s", err)
	}

	client, err := buildClient(cmd.API)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
	defer cancel()

	// Query the API for the Docker configuration.
	dockercfg, err := client.DockerCfg(ctx, &pb.DockerCfgRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to request Docker configuration to pushing built images: %s", err)
	}

	// These are additional environment variables that have been provided outside of this build, with the intent
	// for them to be injected into our running containers.
	//   eg. M8S_ENV_FOO=bar, will inject FOO=bar into the containers.
	extraEnvs := environ.Get()

	for name, service := range dc.Services {
		// Attach our addition environment variables to the service.
		service.Environment = append(service.Environment, extraEnvs...)

		// Build new images if the Docker Compose file is using the "build" option for a service.
		if service.Build != "" {
			fmt.Println("Detected Docker Compose file is using 'build' option. Packaging service:", name)

			tag := fmt.Sprintf("%s-%s", cmd.Name, name)

			err := buildAndPush(service.Build, cmd.DockerRepository, tag, dockercfg)
			if err != nil {
				return fmt.Errorf("failed to build image: %s", err)
			}

			// Pass this on so our API uses this image for the build.
			service.Image = fmt.Sprintf("%s:%s", cmd.DockerRepository, tag)
		}

		dc.Services[name] = service
	}

	ctx, cancel = context.WithTimeout(context.Background(), cmd.Timeout)
	defer cancel()

	// Start the build.
	stream, err := client.Create(ctx, &pb.CreateRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		Metadata: &pb.Metadata{
			Name:    strings.ToLower(cmd.Name),
			Domains: strings.Split(strings.ToLower(cmd.Domains), ","),
			BasicAuth: &pb.BasicAuth{
				User: cmd.BasicAuthUser,
				Pass: cmd.BasicAuthPass,
			},
		},
		GitCheckout: &pb.GitCheckout{
			Repository: cmd.GitRepository,
			Revision:   cmd.GitRevision,
		},
		Compose: dc.GRPC(),
	})
	if err != nil {
		return fmt.Errorf("the build has failed: %s", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read stream: %s", err)
		}

		fmt.Println(string(resp.Message))
	}

	for _, step := range steps {
		ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
		defer cancel()

		stream, err := client.Exec(ctx, &pb.ExecRequest{
			Credentials: &pb.Credentials{
				Token: cmd.Token,
			},
			Name:      strings.ToLower(cmd.Name),
			Container: cmd.ExecInside,
			Command:   step,
		})
		if err != nil {
			return fmt.Errorf("the exec command has failed: %s", err)
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read stream: %s", err)
			}

			fmt.Println(string(resp.Message))
		}
	}

	return nil
}

// Build declares the "build" sub command.
func Build(app *kingpin.Application) {
	c := new(cmdBuild)

	cmd := app.Command("build", "Build the environment").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default(defaultEndpoint).OverrideDefaultFromEnvar("M8S_API").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Default("").OverrideDefaultFromEnvar("M8S_TOKEN").StringVar(&c.Token)
	cmd.Flag("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
	cmd.Flag("domains", "Domains for this environment to run on").Required().StringVar(&c.Domains)
	cmd.Flag("basic-auth-user", "Basic auth user to assign to this environment").Default("").OverrideDefaultFromEnvar("M8S_BASIC_AUTH_USER").StringVar(&c.BasicAuthUser)
	cmd.Flag("basic-auth-pass", "Basic auth user to assign to this environment").Default("").OverrideDefaultFromEnvar("M8S_BASIC_AUTH_PASS").StringVar(&c.BasicAuthPass)
	cmd.Flag("git-repository", "Git repository to clone from").Default("").OverrideDefaultFromEnvar("M8S_GIT_REPO").StringVar(&c.GitRepository)
	cmd.Flag("git-revision", "Git revision to checkout during clone").Required().StringVar(&c.GitRevision)
	cmd.Flag("docker-compose", "Docker Compose file").Default("docker-compose.yml").OverrideDefaultFromEnvar("M8S_DOCKER_COMPOSE").StringVar(&c.DockerCompose)
	cmd.Flag("docker-repository", "Docker repository to push built images").Default("").OverrideDefaultFromEnvar("M8S_DOCKER_REPOSITORY").StringVar(&c.DockerRepository)
	cmd.Flag("exec-file", "Configuration file which contains execution steps").Default("m8s.yml").OverrideDefaultFromEnvar("M8S_EXEC_FILE").StringVar(&c.ExecFile)
	cmd.Flag("exec-step", "Step from the configuration file to use for execution").Default("build").OverrideDefaultFromEnvar("M8S_EXEC_STEP").StringVar(&c.ExecStep)
	cmd.Flag("exec-inside", "Docker repository to push built images").Default("php").OverrideDefaultFromEnvar("M8S_EXEC_INSIDE").StringVar(&c.ExecInside)
	cmd.Flag("timeout", "How long to wait for a step to finish").Default("30m").OverrideDefaultFromEnvar("M8S_TIMEOUT").DurationVar(&c.Timeout)
}

// A helper function to building and pushing an image to a Docker registry.
func buildAndPush(dir, repository, tag string, dockercfg *pb.DockerCfgResponse) error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	creds := docker.AuthConfigurations{
		Configs: map[string]docker.AuthConfiguration{
			dockercfg.Registry: {
				Username:      dockercfg.Username,
				Password:      dockercfg.Password,
				Email:         dockercfg.Email,
				ServerAddress: dockercfg.Registry,
			},
		},
	}

	err = client.BuildImage(docker.BuildImageOptions{
		Name:         fmt.Sprint("%s:%s", repository, tag),
		Dockerfile:   "Dockerfile",
		Pull:         true,
		OutputStream: os.Stdout,
		ContextDir:   dir,
		AuthConfigs:  creds,
	})
	if err != nil {
		return err
	}

	return client.PushImage(docker.PushImageOptions{
		Name: repository,
		Tag:  tag,
	}, creds.Configs[dockercfg.Registry])
}

// Helper function to load testing steps.
func loadSteps(f, step string) ([]string, error) {
	var steps []string

	s, err := ioutil.ReadFile(f)
	if err != nil {
		return steps, err
	}

	y, err := simpleyaml.NewYaml(s)
	if err != nil {
		return steps, err
	}

	raw, err := y.Get(step).Array()
	if err != nil {
		return steps, err
	}

	for _, val := range raw {
		steps = append(steps, to.String(val))
	}

	return steps, nil
}
