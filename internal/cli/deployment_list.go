package cli

import (
	"context"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/posener/complete"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
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
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		sort.Sort(serversort.DeploymentCompleteDesc(resp.Deployments))

		tbl := terminal.NewTable("", "ID", "Registry", "Started", "Completed")

		const bullet = "●"

		for _, b := range resp.Deployments {
			// Determine our bullet
			status := ""
			statusColor := ""
			switch b.Status.State {
			case pb.Status_RUNNING:
				status = bullet
				statusColor = terminal.Yellow

			case pb.Status_SUCCESS:
				status = "✔"
				statusColor = terminal.Green

			case pb.Status_ERROR:
				status = "✖"
				statusColor = terminal.Red
			}

			// Parse our times
			var startTime, completeTime string
			if t, err := ptypes.Timestamp(b.Status.StartTime); err == nil {
				startTime = humanize.Time(t)
			}
			if t, err := ptypes.Timestamp(b.Status.CompleteTime); err == nil {
				completeTime = humanize.Time(t)
			}

			tbl.Rich([]string{
				status,
				b.Id,
				b.Component.Name,
				startTime,
				completeTime,
			}, []string{
				statusColor,
			})
		}

		c.ui.Table(tbl)

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
