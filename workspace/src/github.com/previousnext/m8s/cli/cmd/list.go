package cmd

import (
	"fmt"

	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdList struct {
	API   string
	Token string
}

func (cmd *cmdList) run(c *kingpin.ParseContext) error {
	conn, err := grpc.Dial(cmd.API, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return fmt.Errorf("failed to connect to CI API: %s", err)
	}

	client := pb.NewM8SClient(conn)

	envs, err := client.List(context.Background(), &pb.ListRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
	})
	if err != nil {
		return fmt.Errorf("failed list built environments: %s", err)
	}

	for _, env := range envs.Environments {
		fmt.Println(env)
	}

	return nil
}

// List declares the "list" sub command.
func List(app *kingpin.Application) {
	c := new(cmdList)

	cmd := app.Command("list", "List all the built environments").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default("M8S.ci.pnx.com.au:433").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Required().StringVar(&c.Token)
}
