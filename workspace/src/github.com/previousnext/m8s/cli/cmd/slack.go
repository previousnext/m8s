package cmd

import (
	//"fmt"
	//"github.com/nlopes/slack"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdSlack struct {
	Token   string
	Channel string
	Color   string
	Name    string
}

func (cmd *cmdSlack) run(c *kingpin.ParseContext) error {
	/*client, err := buildClient(cmd.API)
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

	api := slack.New(cmd.Token)

	// Get the details from the api.
	domains := ""
	ssh := "ssh foo~foo~foo@foo.com"

	params := slack.PostMessageParameters{
		Username:  "M8s",
		IconEmoji: ":m8s:",
		Attachments: []slack.Attachment{
			{
				Color: cmd.Color,
				Fields: []slack.AttachmentField{
					{
						Title: "Name",
						Value: cmd.Name,
						Short: true,
					},
					{
						Title: "Domains",
						Value: domains,
						Short: true,
					},
					{
						Title: "SSH",
						Value: ssh,
					},
				},
			},
		},
	}

	_, _, err = api.PostMessage(cmd.Channel, "Environment has been built", params)*/

	return nil
}

// Slack declares the "slack" sub command.
func Slack(app *kingpin.Application) {
	c := new(cmdSlack)

	cmd := app.Command("notify", "Slack notification command for environments").Action(c.run)
	cmd.Flag("slack-token", "Slack token for authentication").Default("").OverrideDefaultFromEnvar("M8S_SLACK_TOKEN").StringVar(&c.Token)
	cmd.Flag("slack-channel", "Slack channel for posting updates").Default("").OverrideDefaultFromEnvar("M8S_SLACK_CHANNEL").StringVar(&c.Channel)
	cmd.Flag("slack-color", "Color to use for Slack notifications").Default("").OverrideDefaultFromEnvar("M8S_SLACK_COLOR").StringVar(&c.Color)
	cmd.Arg("name", "Unique identifier for the environment").Required().StringVar(&c.Name)
}
