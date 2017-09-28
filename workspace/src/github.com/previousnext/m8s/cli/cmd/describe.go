package cmd

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdDescribe struct {
	API   string
	Token string
	Name  string
}

func (cmd *cmdDescribe) run(c *kingpin.ParseContext) error {
	client, err := buildClient(cmd.API)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	describe, err := client.Describe(context.Background(), &pb.DescribeRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		Name: cmd.Name,
	})
	if err != nil {
		return fmt.Errorf("failed list built environments: %s", err)
	}

	fmt.Println(describe)

	return err
}

// Describe declares the "describe" sub command.
func Describe(app *kingpin.Application) {
	c := new(cmdDescribe)

	cmd := app.Command("describe", "Describes an environment").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default(defaultEndpoint).OverrideDefaultFromEnvar("M8S_API").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Default("").OverrideDefaultFromEnvar("M8S_TOKEN").StringVar(&c.Token)
	cmd.Arg("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
}
