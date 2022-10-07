package cli

import (
	"context"
	"strings"

	"github.com/fatih/color"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type LogsCommand struct {
	*baseCommand

	flagDeploySeq string
}

var logColors = map[pb.LogBatch_Entry_Source]*color.Color{
	pb.LogBatch_Entry_APP:        color.New(color.FgGreen),
	pb.LogBatch_Entry_ENTRYPOINT: color.New(color.FgCyan),
}

func (c *LogsCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleAppTarget(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		stream, err := app.Logs(ctx, c.flagDeploySeq)
		if err != nil {
			if !clierrors.IsCanceled(err) {
				app.UI.Output("Error reading logs: %s", err, terminal.WithErrorStyle())
			}
			return ErrSentinel
		}

		for {
			batch, err := stream.Recv()
			if err != nil {
				if !clierrors.IsCanceled(err) {
					app.UI.Output("Error reading logs: %s", err, terminal.WithErrorStyle())
				}

				return ErrSentinel
			}

			if len(batch.Lines) == 0 {
				break
			}

			for _, event := range batch.Lines {
				event.Line = strings.TrimSuffix(event.Line, "\n")

				// We use this format rather than regular RFC3339Nano because we use .0
				// instead of .9, which preserves the spacing so the output is always
				// lined up
				tsRaw := event.Timestamp.AsTime()
				ts := tsRaw.Format("2006-01-02T15:04:05.000Z07:00")
				short := batch.InstanceId
				if len(short) > 6 {
					short = short[len(short)-6:]
				}

				color, ok := logColors[event.Source]
				if !ok {
					color = logColors[pb.LogBatch_Entry_APP]
				}

				header := color.Sprintf("%s %s: ", ts, short)
				if strings.IndexByte(event.Line, '\n') != -1 {
					parts := strings.Split(event.Line, "\n")

					for _, part := range parts {
						m := header + part
						c.ui.Output(m)
					}
				} else {
					m := header + event.Line
					c.ui.Output(m)
				}
			}
		}

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *LogsCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "deployment-seq",
			Usage:  "Get logs for a specific deployment of the app using the deployment sequence number.",
			Target: &c.flagDeploySeq,
		})
	})
}

func (c *LogsCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LogsCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LogsCommand) Synopsis() string {
	return "Show log output from an application's current deployment"
}

func (c *LogsCommand) Help() string {
	return formatHelp(`
Usage: waypoint logs [options]

  Show log output from all current deployments.

  The logs will include output from deployments that aren't released.
  As new deployments are made, their logs will appear automatically.

  The six character text after the date on a log line is the last six
  characters of the instance ID. This can be used to trace any logs back
  to a specific deployment or filter out certain log messages.

` + c.Flags().Help())
}
