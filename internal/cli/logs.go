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
)

type LogsCommand struct {
	*baseCommand
}

var headerColor = color.New(color.FgCyan)

func (c *LogsCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		lv, err := app.Logs(ctx)
		if err != nil {
			if !clierrors.IsCanceled(err) {
				app.UI.Output("Error reading logs: %s", err, terminal.WithErrorStyle())
			}
			return ErrSentinel
		}

		for {
			batch, err := lv.NextLogBatch(ctx)
			if err != nil {
				if !clierrors.IsCanceled(err) {
					app.UI.Output("Error reading logs: %s", err, terminal.WithErrorStyle())
				}
				return ErrSentinel
			}

			if len(batch) == 0 {
				break
			}

			for _, event := range batch {
				event.Message = strings.TrimSuffix(event.Message, "\n")

				// We use this format rather than regular RFC3339Nano because we use .0
				// instead of .9, which preserves the spacing so the output is always
				// lined up
				ts := event.Timestamp.Format("2006-01-02T15:04:05.000Z07:00")
				short := event.Partition
				if len(short) > 6 {
					short = short[len(short)-6:]
				}

				header := headerColor.Sprintf("%s %s: ", ts, short)
				if strings.IndexByte(event.Message, '\n') != -1 {
					parts := strings.Split(event.Message, "\n")

					for _, part := range parts {
						m := header + part
						c.ui.Output(m)
					}
				} else {
					m := header + event.Message
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
	return c.flagSet(0, nil)
}

func (c *LogsCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LogsCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *LogsCommand) Synopsis() string {
	return "Show log output from the current application deployment"
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
