package cmd

import (
	"context"
	"fmt"

	pb "github.com/previousnext/pr/pb"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdDockerCfgSet struct {
	API      string
	Token    string
	Registry string
	Username string
	Password string
	Email    string
	Auth     string
}

func (cmd *cmdDockerCfgSet) run(c *kingpin.ParseContext) error {
	client, err := buildClient(cmd.API)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	_, err = client.DockerCfgSet(context.Background(), &pb.DockerCfgSetRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		DockerCfg: &pb.DockerCfg{
			Registry: cmd.Registry,
			Username: cmd.Username,
			Password: cmd.Password,
			Email:    cmd.Email,
			Auth:     cmd.Auth,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// DockerCfgSet declares the "dockercfg-set" sub command.
func DockerCfgSet(app *kingpin.Application) {
	c := new(cmdDockerCfgSet)

	cmd := app.Command("dockercfg-set", "Updates the Docker secret configuration").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default("pr.ci.pnx.com.au:433").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Required().StringVar(&c.Token)
	cmd.Flag("registry", "The registry to store Docker images").Default("https://index.docker.io/v1/").StringVar(&c.Registry)
	cmd.Flag("username", "Username credential to use for the registry").Required().StringVar(&c.Username)
	cmd.Flag("password", "Password credential to use for the registry").Required().StringVar(&c.Password)
	cmd.Flag("email", "Email credential to use for the registry").Required().StringVar(&c.Email)
	cmd.Flag("auth", "Auth token credential to use for the registry").Required().StringVar(&c.Auth)
}
