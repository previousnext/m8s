package cmd

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdNotify struct {
	API          string
	Token        string
	SlackToken   string
	SlackChannel string
	SlackColor   string
	Container    string
	Name         string
}

func (cmd *cmdNotify) run(c *kingpin.ParseContext) error {
	client, err := buildClient(cmd.API)
	if err != nil {
		return errors.Wrap(err, "failed to build client")
	}

	describe, err := client.Describe(context.Background(), &pb.DescribeRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		Name: strings.ToLower(cmd.Name),
	})
	if err != nil {
		return errors.Wrap(err, "failed to describe environment")
	}

	api := slack.New(cmd.SlackToken)

	msg := slack.PostMessageParameters{
		Username:  "M8s",
		IconEmoji: ":m8s:",
		Attachments: []slack.Attachment{
			{
				Color: cmd.SlackColor,
				Fields: []slack.AttachmentField{
					{
						Title: "Domains",
						Value: strings.Join(describe.Domains, "\n"),
					},
					{
						Title: "Containers",
						Value: strings.Join(describe.Containers, ", "),
					},
					{
						Title: "Command Line Access",
						Value: fmt.Sprintf("ssh %s~%s~%s~$(whoami)@%s", describe.Namespace, describe.Name, cmd.Container, describe.SSH),
					},
				},
			},
		},
	}

	_, _, err = api.PostMessage(cmd.SlackChannel, fmt.Sprintf("Temporary environment is ready: *%s*", describe.Name), msg)

	return nil
}

// Notify declares the "slack" sub command.
func Notify(app *kingpin.Application) {
	c := new(cmdNotify)

	cmd := app.Command("notify", "Slack notification command for environments").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default(defaultEndpoint).OverrideDefaultFromEnvar("M8S_API").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Default("").OverrideDefaultFromEnvar("M8S_TOKEN").StringVar(&c.Token)
	cmd.Flag("slack-token", "Slack token for authentication").Default("").OverrideDefaultFromEnvar("M8S_SLACK_TOKEN").StringVar(&c.SlackToken)
	cmd.Flag("slack-channel", "Slack channel for posting updates").Default("").OverrideDefaultFromEnvar("M8S_SLACK_CHANNEL").StringVar(&c.SlackChannel)
	cmd.Flag("slack-color", "Color to use for Slack notifications").Default("#32cd32").OverrideDefaultFromEnvar("M8S_SLACK_COLOR").StringVar(&c.SlackColor)
	cmd.Flag("container", "Default container name to use").Default("php").OverrideDefaultFromEnvar("M8S_EXEC_INSIDE").StringVar(&c.Container)
	cmd.Arg("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
}
