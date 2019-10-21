package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/previousnext/m8s/cmd/config"

	"github.com/pkg/errors"
	"github.com/previousnext/compose"
	"github.com/previousnext/m8s/cmd/environ"
	"github.com/previousnext/m8s/cmd/metadata"
	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// ServiceSkip will skip a service if the annotation is set.
	ServiceSkip = "m8s.io/skip"
	// ServiceType is used for indentifying the type of service for extra handling.
	ServiceType = "m8s.io/type"
)

type cmdBuild struct {
	API           string
	Insecure      bool
	Token         string
	Name          string
	Domains       string
	BasicAuthUser string
	BasicAuthPass string
	Retention     time.Duration
	GitRepository string
	GitRevision   string
	DockerCompose []string
	ExecFile      string
	ExecStep      string
	ExecInside    string
	Timeout       time.Duration
}

func (cmd *cmdBuild) run(c *kingpin.ParseContext) error {
	// Load the Docker Compose file, we are going to use alot of its
	// configuration for this build.
	dc, err := compose.Load(cmd.DockerCompose)
	if err != nil {
		return errors.Wrap(err, "failed to load Docker Compose file")
	}

	// Load the steps required to run the build, these are bespoke steps used
	// for bootstrapping and testing the application.

	cfg, err := config.Load(cmd.ExecFile)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	client, err := buildClient(cmd.API, cmd.Insecure)
	if err != nil {
		return errors.Wrap(err, "failed to connect")
	}

	// These are additional environment variables that have been provided outside of this build, with the intent
	// for them to be injected into our running containers.
	//   eg. M8S_ENV_FOO=bar, will inject FOO=bar into the containers.
	extraEnvs := environ.Get()

	for name, service := range dc.Services {
		// Attach our addition environment variables to the service.
		service.Environment = append(service.Environment, extraEnvs...)
		dc.Services[name] = service
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
	defer cancel()

	annotations, err := metadata.Annotations(os.Environ())
	if err != nil {
		return errors.Wrap(err, "failed to get annotations")
	}

	// Start the build.
	createRequest := pb.CreateRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		Metadata: &pb.Metadata{
			Name:        strings.ToLower(cmd.Name),
			Annotations: annotations,
			Domains:     strings.Split(strings.ToLower(cmd.Domains), ","),
			BasicAuth: &pb.BasicAuth{
				User: cmd.BasicAuthUser,
				Pass: cmd.BasicAuthPass,
			},
			Retention: cmd.Retention.String(),
		},
		GitCheckout: &pb.GitCheckout{
			Repository: cmd.GitRepository,
			Revision:   cmd.GitRevision,
		},
		Compose: composeToGRPC(dc),
	}
	for _, init := range cfg.Init {
		createRequest.Steps = append(createRequest.Steps, &pb.Init{
			Name:  init.Name,
			Image: init.Image,
			Reservations: &pb.Resource{
				CPU:    init.Resources.Reservations.CPUs,
				Memory: init.Resources.Reservations.Memory,
			},
			Limits: &pb.Resource{
				CPU:    init.Resources.Limits.CPUs,
				Memory: init.Resources.Limits.Memory,
			},
			Steps:   init.Steps,
			Volumes: init.Volumes,
		})
	}
	stream, err := client.Create(ctx, &createRequest)
	if err != nil {
		return errors.Wrap(err, "the build has failed")
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read stream")
		}

		fmt.Print(resp.Message)
	}

	for _, step := range cfg.Build {
		ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
		defer cancel()

		stream, err := client.Step(ctx, &pb.StepRequest{
			Credentials: &pb.Credentials{
				Token: cmd.Token,
			},
			Name:      strings.ToLower(cmd.Name),
			Container: cmd.ExecInside,
			Command:   step,
		})
		if err != nil {
			return errors.Wrap(err, "the exec command has failed")
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return errors.Wrap(err, "failed to read stream")
			}

			fmt.Print(resp.Message)
		}
	}

	return nil
}

// Build declares the "build" sub command.
func Build(app *kingpin.Application) {
	c := new(cmdBuild)

	cmd := app.Command("build", "Build the environment").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default(defaultEndpoint).OverrideDefaultFromEnvar("M8S_API").StringVar(&c.API)
	cmd.Flag("insecure", "Use insecure connections for API interactions").BoolVar(&c.Insecure)
	cmd.Flag("token", "Token used for authenticating with the API service").Default("").OverrideDefaultFromEnvar("M8S_TOKEN").StringVar(&c.Token)
	cmd.Flag("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
	cmd.Flag("domains", "Domains for this environment to run on").Required().StringVar(&c.Domains)
	cmd.Flag("basic-auth-user", "Basic auth user to assign to this environment").Default("").OverrideDefaultFromEnvar("M8S_BASIC_AUTH_USER").StringVar(&c.BasicAuthUser)
	cmd.Flag("basic-auth-pass", "Basic auth user to assign to this environment").Default("").OverrideDefaultFromEnvar("M8S_BASIC_AUTH_PASS").StringVar(&c.BasicAuthPass)
	cmd.Flag("retention", "How long to keep an environment").Default("120h").OverrideDefaultFromEnvar("M8S_RETENTION").DurationVar(&c.Retention)
	cmd.Flag("git-repository", "Git repository to clone from").Default("").OverrideDefaultFromEnvar("M8S_GIT_REPO").StringVar(&c.GitRepository)
	cmd.Flag("git-revision", "Git revision to checkout during clone").Required().StringVar(&c.GitRevision)
	cmd.Flag("docker-compose", "Docker Compose file(s)").Default("docker-compose.yml").OverrideDefaultFromEnvar("M8S_DOCKER_COMPOSE").StringsVar(&c.DockerCompose)
	cmd.Flag("exec-file", "Configuration file which contains execution steps").Default("m8s.yml").OverrideDefaultFromEnvar("M8S_EXEC_FILE").StringVar(&c.ExecFile)
	cmd.Flag("exec-step", "Step from the configuration file to use for execution").Default("build").OverrideDefaultFromEnvar("M8S_EXEC_STEP").StringVar(&c.ExecStep)
	cmd.Flag("exec-inside", "Docker repository to push built images").Default("app").OverrideDefaultFromEnvar("M8S_EXEC_INSIDE").StringVar(&c.ExecInside)
	cmd.Flag("timeout", "How long to wait for a step to finish").Default("30m").OverrideDefaultFromEnvar("M8S_TIMEOUT").DurationVar(&c.Timeout)
}

// Helper function used for marshalling a Docker Compose file into a M8s object.
func composeToGRPC(dc compose.DockerCompose) *pb.Compose {
	resp := new(pb.Compose)

	for name, service := range dc.Services {
		if _, ok := service.Labels[ServiceSkip]; ok {
			fmt.Println("Skipping service:", name)
			continue
		}

		newService := &pb.ComposeService{
			Name:         name,
			Image:        service.Image,
			Entrypoint:   service.Entrypoint,
			Volumes:      service.Volumes,
			Ports:        service.Ports,
			Environment:  service.Environment,
			Tmpfs:        service.Tmpfs,
			Capabilities: service.CapAdd,
			Extrahosts:	  service.ExtraHosts,
			Limits: &pb.Resource{
				CPU:    service.Deploy.Resources.Limits.CPUs,
				Memory: service.Deploy.Resources.Limits.Memory,
			},
			Reservations: &pb.Resource{
				CPU:    service.Deploy.Resources.Reservations.CPUs,
				Memory: service.Deploy.Resources.Reservations.Memory,
			},
		}

		if val, ok := service.Labels[ServiceType]; ok {
			newService.Type = val
		}

		resp.Services = append(resp.Services, newService)
	}

	return resp
}
