package cli

import (
	"context"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/olekukonko/tablewriter"
	"github.com/posener/complete"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serversort "github.com/hashicorp/waypoint/internal/server/sort"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type DeploymentListCommand struct {
	*baseCommand

	flagWorkspaceAll bool
}

func (c *DeploymentListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	// Get our API client
	client := c.project.Client()

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		var wsRef *pb.Ref_Workspace
		if !c.flagWorkspaceAll {
			wsRef = c.project.WorkspaceRef()
		}

		// List builds
		resp, err := client.ListDeployments(c.Ctx, &pb.ListDeploymentsRequest{
			Application: app.Ref(),
			Workspace:   wsRef,
		})
		if err != nil {
			c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}
		sort.Sort(serversort.DeploymentCompleteDesc(resp.Deployments))

		// Get our direct stdout handle cause we're going to be writing colors
		// and want color detection to work.
		out, _, err := c.project.UI.OutputWriters()
		if err != nil {
			c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		const bullet = "●"

		table := tablewriter.NewWriter(out)
		table.SetHeader([]string{"", "ID", "Registry", "Started", "Completed"})
		table.SetBorder(false)
		for _, b := range resp.Deployments {
			// Determine our bullet
			status := ""
			statusColor := tablewriter.Colors{}
			switch b.Status.State {
			case pb.Status_RUNNING:
				status = bullet
				statusColor = tablewriter.Colors{tablewriter.FgYellowColor}

			case pb.Status_SUCCESS:
				status = "✔"
				statusColor = tablewriter.Colors{tablewriter.FgGreenColor}

			case pb.Status_ERROR:
				status = "✖"
				statusColor = tablewriter.Colors{tablewriter.FgRedColor}
			}

			// Parse our times
			var startTime, completeTime string
			if t, err := ptypes.Timestamp(b.Status.StartTime); err == nil {
				startTime = humanize.Time(t)
			}
			if t, err := ptypes.Timestamp(b.Status.CompleteTime); err == nil {
				completeTime = humanize.Time(t)
			}

			table.Rich([]string{
				status,
				b.Id,
				b.Component.Name,
				startTime,
				completeTime,
			}, []tablewriter.Colors{
				statusColor,
				tablewriter.Colors{},
				tablewriter.Colors{},
				tablewriter.Colors{},
				tablewriter.Colors{},
			})
		}
		table.Render()

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *DeploymentListCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "workspace-all",
			Target: &c.flagWorkspaceAll,
			Usage:  "List builds in all workspaces for this project and application.",
		})
	})
}

func (c *DeploymentListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *DeploymentListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *DeploymentListCommand) Synopsis() string {
	return "List deployments."
}

func (c *DeploymentListCommand) Help() string {
	helpText := `
Usage: waypoint deployment list [options]

  Lists the deployments that were created.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
