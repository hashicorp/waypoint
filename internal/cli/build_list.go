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

type BuildListCommand struct {
	*baseCommand

	flagWorkspaceAll bool
}

func (c *BuildListCommand) Run(args []string) int {
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
		resp, err := client.ListBuilds(c.Ctx, &pb.ListBuildsRequest{
			Application: app.Ref(),
			Workspace:   wsRef,
		})
		if err != nil {
			c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}
		sort.Sort(serversort.BuildStartDesc(resp.Builds))

		// Get our direct stdout handle cause we're going to be writing colors
		// and want color detection to work.
		out, _, err := c.project.UI.OutputWriters()
		if err != nil {
			c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		const bullet = "●"

		table := tablewriter.NewWriter(out)
		table.SetHeader([]string{"", "ID", "Workspace", "Builder", "Started", "Completed"})
		table.SetBorder(false)
		for _, b := range resp.Builds {
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
				b.Workspace.Workspace,
				b.Component.Name,
				startTime,
				completeTime,
			}, []tablewriter.Colors{
				statusColor,
				tablewriter.Colors{},
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

func (c *BuildListCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "workspace-all",
			Target: &c.flagWorkspaceAll,
			Usage:  "List builds in all workspaces for this project and application.",
		})
	})
}

func (c *BuildListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *BuildListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *BuildListCommand) Synopsis() string {
	return "Build a new versioned artifact from source."
}

func (c *BuildListCommand) Help() string {
	helpText := `
Usage: waypoint artifact build [options]
Alias: waypoint build

  Build a new versioned artifact from source.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
