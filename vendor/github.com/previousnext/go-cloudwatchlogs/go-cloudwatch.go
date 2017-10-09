package cloudwatchlogs

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func Events(region, group, stream, start, end string) (Logs, error) {
	var logs Logs

	svc := cloudwatchlogs.New(session.New(), &aws.Config{Region: aws.String(region)})

	startDuration, err := time.ParseDuration(start)
	if err != nil {
		return logs, err
	}

	endDuration, err := time.ParseDuration(end)
	if err != nil {
		return logs, err
	}

	var (
		from = aws.TimeUnixMilli(time.Now().Add(-startDuration).UTC())
		to   = aws.TimeUnixMilli(time.Now().Add(-endDuration).UTC())
	)

	streams, err := streams(svc, group, stream)
	if err != nil {
		return logs, err
	}

	for _, s := range streams {
		// Ensure that we are not querying for streams which have finished prior to
		if s.LastEventTimestamp != nil && *s.LastEventTimestamp < from {
			continue
		}

		newLogs, err := events(svc, group, *s.LogStreamName, from, to)
		if err != nil {
			return logs, err
		}

		if len(newLogs) > 0 {
			logs = MergeLogs(logs, newLogs)
		}
	}

	return logs, nil
}

// Helper function to get Log Streams (with Token support).
func streams(svc *cloudwatchlogs.CloudWatchLogs, group, prefix string) ([]*cloudwatchlogs.LogStream, error) {
	var streams []*cloudwatchlogs.LogStream

	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(group),
		LogStreamNamePrefix: aws.String(prefix),
		Descending:          aws.Bool(true),
	}

	for {
		resp, err := svc.DescribeLogStreams(params)
		if err != nil {
			return streams, err
		}

		for _, s := range resp.LogStreams {
			streams = append(streams, s)
		}

		if resp.NextToken == nil {
			return streams, nil
		}

		params.NextToken = resp.NextToken
	}

	return streams, nil
}

// Helper function to get Events (with Token support).
func events(svc *cloudwatchlogs.CloudWatchLogs, group, stream string, from, to int64) (Logs, error) {
	var logs Logs

	params := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
		StartTime:     aws.Int64(from),
		EndTime:       aws.Int64(to),
	}

	for {
		resp, err := svc.GetLogEvents(params)
		if err != nil {
			return logs, err
		}

		if len(resp.Events) < 1 {
			return logs, nil
		}

		for _, e := range resp.Events {
			logs = append(logs, &Log{
				Stream:    stream,
				Timestamp: time.Unix(*e.Timestamp/1000, 0),
				Message:   *e.Message,
			})
		}

		if resp.NextForwardToken == nil {
			return logs, nil
		}

		params.NextToken = resp.NextForwardToken
	}

	return logs, nil
}
