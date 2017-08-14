package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/gosexy/to"
	"github.com/previousnext/pr/cli/compose"
	pb "github.com/previousnext/pr/pb"
	"github.com/smallfish/simpleyaml"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdBuild struct {
	API              string
	Token            string
	Name             string
	Domains          string
	GitRepository    string
	GitRevision      string
	DockerCompose    string
	DockerRepository string
	ExecFile         string
	ExecStep         string
	ExecInside       string
	Keep             string
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

	// Query the API for the Docker configuration.
	dockercfg, err := client.DockerCfgGet(context.Background(), &pb.DockerCfgGetRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to request Docker configuration to pushing built images: %s", err)
	}

	// Build new images if the Docker Compose file is using the "build" option for a service.
	for name, service := range dc.Services {
		if service.Build != "" {
			fmt.Println("Detected Docker Compose file is using 'build' option. Packaging service:", name)

			tag := fmt.Sprintf("%s-%s", cmd.Name, name)

			err := buildAndPush(service.Build, cmd.DockerRepository, tag, dockercfg.DockerCfg)
			if err != nil {
				return fmt.Errorf("failed to build image: %s", err)
			}

			// Pass this on so our API uses this image for the build.
			service.Image = fmt.Sprintf("%s:%s", cmd.DockerRepository, tag)

			dc.Services[name] = service
		}
	}

	// Start the build.
	stream, err := client.Build(context.Background(), &pb.BuildRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		Metadata: &pb.Metadata{
			Name:    cmd.Name,
			Domains: strings.Split(cmd.Domains, ","),
		},
		GitCheckout: &pb.GitCheckout{
			Repository: cmd.GitRepository,
			Revision:   cmd.GitRevision,
		},
		Compose: dc.GRPC(),
		Exec: &pb.Exec{
			Container: cmd.ExecInside,
			Steps:     steps,
		},
		Keep: cmd.Keep,
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

	return nil
}

// Build declares the "build" sub command.
func Build(app *kingpin.Application) {
	c := new(cmdBuild)

	cmd := app.Command("build", "Build the environment").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default("pr.ci.pnx.com.au:433").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Required().StringVar(&c.Token)
	cmd.Flag("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
	cmd.Flag("domains", "Domains for this environment to run on").Required().StringVar(&c.Domains)
	cmd.Flag("git-repository", "Git repository to clone from").Required().StringVar(&c.GitRepository)
	cmd.Flag("git-revision", "Git revision to checkout during clone").Required().StringVar(&c.GitRevision)
	cmd.Flag("docker-compose", "Docker Compose file").Default("docker-compose.yml").StringVar(&c.DockerCompose)
	cmd.Flag("docker-repository", "Docker repository to push built images").Required().StringVar(&c.DockerRepository)
	cmd.Flag("exec-file", "Configuration file which contains execution steps").Required().StringVar(&c.ExecFile)
	cmd.Flag("exec-step", "Step from the configuration file to use for execution").Required().StringVar(&c.ExecStep)
	cmd.Flag("exec-inside", "Docker repository to push built images").Required().StringVar(&c.ExecInside)
	cmd.Flag("Keep", "How many days before an environment can be deleted").Default("5d").StringVar(&c.Keep)
}

// A helper function to building and pushing an image to a Docker registry.
func buildAndPush(dir, repository, tag string, dockercfg *pb.DockerCfg) error {
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
