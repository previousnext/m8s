package cmd

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	pb "github.com/previousnext/m8s/pb"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdSlack struct {
	API          string
	Token        string
	SlackToken   string
	SlackChannel string
	SlackColor   string
	Name         string
}

func (cmd *cmdSlack) run(c *kingpin.ParseContext) error {
	client, err := buildClient(cmd.API)
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	describe, err := client.Describe(context.Background(), &pb.DescribeRequest{
		Credentials: &pb.Credentials{
			Token: cmd.Token,
		},
		Name: strings.ToLower(cmd.Name),
	})
	if err != nil {
		return fmt.Errorf("failed list built environments: %s", err)
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
						Title: "SSH",
						Value: "m8s ssh cwi-m8s <container>",
					},
				},
			},
		},
	}

	_, _, err = api.PostMessage(cmd.SlackChannel, fmt.Sprintf("Environment has been built: *%s*", describe.Name), msg)

	return nil
}

// Slack declares the "slack" sub command.
func Slack(app *kingpin.Application) {
	c := new(cmdSlack)

	cmd := app.Command("notify", "Slack notification command for environments").Action(c.run)
	cmd.Flag("api", "API endpoint which accepts our build requests").Default(defaultEndpoint).OverrideDefaultFromEnvar("M8S_API").StringVar(&c.API)
	cmd.Flag("token", "Token used for authenticating with the API service").Default("").OverrideDefaultFromEnvar("M8S_TOKEN").StringVar(&c.Token)
	cmd.Flag("slack-token", "Slack token for authentication").Default("").OverrideDefaultFromEnvar("M8S_SLACK_TOKEN").StringVar(&c.SlackToken)
	cmd.Flag("slack-channel", "Slack channel for posting updates").Default("").OverrideDefaultFromEnvar("M8S_SLACK_CHANNEL").StringVar(&c.SlackChannel)
	cmd.Flag("slack-color", "Color to use for Slack notifications").Default("#32cd32").OverrideDefaultFromEnvar("M8S_SLACK_COLOR").StringVar(&c.SlackColor)
	cmd.Arg("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
}
