package main

import (
	"fmt"

	cloudwatchlogs "github.com/previousnext/go-cloudwatchlogs"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cliRegion = kingpin.Flag("region", "Region which logs reside").Default("ap-southeast-2").String()
	cliGroup  = kingpin.Flag("group", "CloudWatch Logs group").Required().String()
	cliStream = kingpin.Flag("stream", "CloudWatch Logs stream").String()
	cliStart  = kingpin.Flag("start", "Time ago to search from").Default("10m").String()
	cliEnd    = kingpin.Flag("end", "Time ago to end search").Default("0").String()
)

func main() {
	kingpin.Parse()
	logs, err := cloudwatchlogs.GetStreams(*cliRegion, *cliGroup, *cliStream, *cliStart, *cliEnd)
	if err != nil {
		panic(err)
	}

	for _, l := range logs {
		fmt.Println(l.Message)
	}
}
