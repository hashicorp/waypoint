// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lambda

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/aws/utils"
)

// Logs fetches logs from cloudwatch
func (p *Platform) Logs(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	es *component.LogViewer,
	app *component.Source,
	dep *Deployment,
) error {
	defer log.Debug("finished with cloudwatchlogs")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region: p.config.Region,
		Logger: log,
	})
	if err != nil {
		return err
	}

	logs := cloudwatchlogs.New(sess)

	group := fmt.Sprintf("/aws/lambda/%s", app.App)

	var lastLSToken *string

	for {
		streams, err := logs.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(group),
			Descending:   aws.Bool(false),
			OrderBy:      aws.String("LastEventTime"),
			NextToken:    lastLSToken,
		})

		if err != nil {
			return errors.Wrapf(err, "failed to describe log stream for group %q in region %q", group, p.config.Region)
		}

		if len(streams.LogStreams) == 0 {
			return nil
		}

		lastLSToken = streams.NextToken

		limit := int64(es.Limit)
		if limit == 0 {
			limit = -1
		}

		log.Debug("fetching log events", "streams", len(streams.LogStreams))

		// 2021/02/16/[25]
		filterRe, err := regexp.Compile(`\d{1,5}/\d{1,2}/\d{1,2}/\[` + dep.Version + `\]`)
		if err != nil {
			return err
		}

		for _, stream := range streams.LogStreams {
			if !filterRe.MatchString(*stream.LogStreamName) {
				continue
			}

			log.Debug("fetching stream", "stream", *stream.LogStreamName)

			gei := &cloudwatchlogs.GetLogEventsInput{
				StartFromHead: aws.Bool(true),
				LogGroupName:  aws.String(group),
				LogStreamName: stream.LogStreamName,
			}

			if !es.StartingAt.IsZero() {
				gei.StartTime = aws.Int64(int64(aws.TimeUnixMilli(es.StartingAt)))
			}

			for {
				if limit >= 0 {
					gei.Limit = &limit
				}

				output, err := logs.GetLogEvents(gei)
				if err != nil {
					return err
				}

				// this stream has no more logs, switch to the next one
				if len(output.Events) == 0 {
					break
				}

				log.Debug("chunk of cloudwatch logs",
					"size", len(output.Events),
					"stream", stream.LogStreamName,
					"token", *output.NextForwardToken,
					"start-time", *output.Events[0].IngestionTime,
				)

				gei.NextToken = output.NextForwardToken

				for _, ev := range output.Events {
					cle := component.LogEvent{
						Partition: *stream.LogStreamName,
						Timestamp: aws.MillisecondsTimeValue(ev.Timestamp),
						Message:   strings.TrimRight(*ev.Message, "\n\t"),
					}

					select {
					case <-ctx.Done():
						return ctx.Err()
					case es.Output <- cle:
						// ok
					}
				}

				log.Debug("processed cloudwatch log chunk")

				if limit >= 0 {
					limit -= int64(len(output.Events))
					if limit <= 0 {
						return nil
					}
				}
			}
		}

		if lastLSToken == nil {
			break
		}
	}

	return nil
}
