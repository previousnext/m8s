package cmd

import (
	"fmt"
	"io"
	"strings"
	"time"

	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdExec struct {
	API     string
	Token   string
	Name    string
	Inside  string
	Command string
	Timeout time.Duration
}

func (cmd *cmdExec) run(c *kingpin.ParseContext) error {
	client, err := buildClient(cmd.API)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
	defer cancel()

	stream, err := client.Exec(ctx, &pb.ExecRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		Name:      strings.ToLower(cmd.Name),
		Container: cmd.Inside,
		Command:   cmd.Command,
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

	return nil
}

// Exec declares the "exec" sub command.
func Exec(app *kingpin.Application) {
	c := new(cmdExec)

	cmd := app.Command("exec", "Exec command in the environment").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default(defaultEndpoint).OverrideDefaultFromEnvar("M8S_API").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Default("").OverrideDefaultFromEnvar("M8S_TOKEN").StringVar(&c.Token)
	cmd.Flag("timeout", "How long to wait for a step to finish").Default("30m").OverrideDefaultFromEnvar("M8S_TIMEOUT").DurationVar(&c.Timeout)
	cmd.Arg("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
	cmd.Arg("inside", "Unique identifier for the environment").Required().StringVar(&c.Inside)
	cmd.Arg("command", "Unique identifier for the environment").Required().StringVar(&c.Command)
}
