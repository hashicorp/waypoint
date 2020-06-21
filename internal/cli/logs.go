package cli

import (
	"context"
	"strings"

	"github.com/fatih/color"
	"github.com/posener/complete"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
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

	client := c.project.Client()
	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		// Get the latest deployment
		resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Order: &pb.OperationOrder{
				Limit: 1,
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if len(resp.Deployments) == 0 {
			app.UI.Output("No successful deployments found.", terminal.WithErrorStyle())
			return ErrSentinel
		}

		lv, err := app.Logs(ctx, resp.Deployments[0])
		if err != nil {
			app.UI.Output("Error reading logs: %s", err, terminal.WithErrorStyle())
			return ErrSentinel
		}

		var pv component.PartitionViewer
		for {
			batch, err := lv.NextLogBatch(ctx)
			if err != nil {
				app.UI.Output("Error reading logs: %s", err, terminal.WithErrorStyle())
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
				short := pv.Short(event.Partition)

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
	return ""
}

func (c *LogsCommand) Help() string {
	return ""
}
