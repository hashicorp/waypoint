package cli

import (
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olekukonko/tablewriter"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serversort "github.com/hashicorp/waypoint/internal/server/sort"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ArtifactListCommand struct {
	*baseCommand
}

func (c *ArtifactListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	// Get our API client
	client := c.project.Client()

	// List builds
	resp, err := client.ListPushedArtifacts(c.Ctx, &empty.Empty{})
	if err != nil {
		c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}
	sort.Sort(serversort.ArtifactStartDesc(resp.Artifacts))

	// Get our direct stdout handle cause we're going to be writing colors
	// and want color detection to work.
	out, _, err := c.project.UI.OutputWriters()
	if err != nil {
		c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	const bullet = "●"

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"", "ID", "Registry", "Started", "Completed"})
	table.SetBorder(false)
	for _, b := range resp.Artifacts {
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

	return 0
}

func (c *ArtifactListCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ArtifactListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ArtifactListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ArtifactListCommand) Synopsis() string {
	return "List pushed artifacts."
}

func (c *ArtifactListCommand) Help() string {
	helpText := `
Usage: waypoint artifact list [options]

  Lists the artifacts that are pushed to a registry. This does not
  list the artifacts that are just part of local builds.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
