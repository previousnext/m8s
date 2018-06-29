package types

// ListOutput is used for returning a list of environment.
type ListOutput struct {
	Environments []Environment `environments`
}

// Environment is used for describing an environment.
type Environment struct {
	Name       string   `json:"name"`
	Domain     string   `json:"domain"`
	Containers []string `json:"containers"`
}

// LogsOutput is used for returning logs data.
type LogsOutput struct {
	Logs string `json:"logs"`
}
