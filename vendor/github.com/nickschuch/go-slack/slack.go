package slack

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

func Send(user, emoji, text, url string) error {
	request := gorequest.New()
	_, body, errs := request.Post(url).
		Set("Notes", "gorequst is coming!").
		Send(Message{
			Username: user,
			Emoji:    emoji,
			Text:     text,
		}).
		End()

	if len(errs) > 0 {
		return fmt.Errorf(body)
	}

	return nil
}
