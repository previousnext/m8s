package cloudwatchlogs

import "time"

type Log struct {
	Stream    string    `json:"stream"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type Logs []*Log
