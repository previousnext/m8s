package cmd

import (
	"fmt"
	"io/ioutil"

	pb "github.com/previousnext/pr/pb"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdSSHSet struct {
	API        string
	Token      string
	KnownHosts string
	PrivateKey string
}

func (cmd *cmdSSHSet) run(c *kingpin.ParseContext) error {
	client, err := buildClient(cmd.API)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	// Load the "known hosts" file.
	known, err := ioutil.ReadFile(cmd.KnownHosts)
	if err != nil {
		return fmt.Errorf("failed to load known hosts file: %s", err)
	}

	// Load the "private key" file.
	private, err := ioutil.ReadFile(cmd.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to load private key file: %s", err)
	}

	_, err = client.SSHSet(context.Background(), &pb.SSHSetRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		SSH: &pb.SSH{
			KnownHosts: known,
			PrivateKey: private,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// SSHSet declares the "ssh-set" sub command.
func SSHSet(app *kingpin.Application) {
	c := new(cmdSSHSet)

	cmd := app.Command("ssh-set", "Updates the SSH secret configuration").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default("pr.ci.pnx.com.au:433").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Required().StringVar(&c.Token)
	cmd.Flag("known-hosts", "Path to the known_hosts file").Required().StringVar(&c.KnownHosts)
	cmd.Flag("private-key", "Path to the private key file").Required().StringVar(&c.PrivateKey)
}
