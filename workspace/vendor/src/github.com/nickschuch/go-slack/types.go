package slack

type Message struct {
	Username string `json:"username"`
	Emoji    string `json:"icon_emoji"`
	Channel  string `json:"channel"`
	Text     string `json:"text"`
}
