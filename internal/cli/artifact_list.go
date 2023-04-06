// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ArtifactListCommand struct {
	*baseCommand

	flagWorkspaceAll bool
	flagVerbose      bool
	flagJson         bool
	flagId           idFormat
	filterFlags      filterFlags
}

func (c *ArtifactListCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	// Get our API client
	client := c.project.Client()

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		if !c.flagJson {
			// UI -- this should happen at the top so that the app name shows clearly
			// for any errors we may encounter prior to the actual table output
			// but we also don't want to corrupt the json
			app.UI.Output("%s", app.Ref().Application, terminal.WithHeaderStyle())
		}

		var wsRef *pb.Ref_Workspace
		if !c.flagWorkspaceAll {
			wsRef = c.project.WorkspaceRef()
		}

		// List builds
		resp, err := client.ListPushedArtifacts(c.Ctx, &pb.ListPushedArtifactsRequest{
			Application:  app.Ref(),
			Workspace:    wsRef,
			Order:        c.filterFlags.orderOp(),
			IncludeBuild: true,
		})
		if err != nil {
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if len(resp.Artifacts) == 0 {
			c.project.UI.Output(
				"No artifacts found for application %q",
				app.Ref().Application,
				terminal.WithWarningStyle(),
			)
			return nil
		}

		if c.flagJson {
			return c.displayJson(resp.Artifacts)
		}

		const bullet = "●"

		table := terminal.NewTable("", "ID", "Registry", "Details", "Started", "Completed")
		for _, b := range resp.Artifacts {
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
			if b.Status.StartTime != nil {
				startTime = humanize.Time(b.Status.StartTime.AsTime())
			}
			if b.Status.CompleteTime != nil {
				completeTime = humanize.Time(b.Status.CompleteTime.AsTime())
			}

			var (
				extraDetails []string
				details      []string
			)

			if user, ok := b.Labels["common/user"]; ok {
				details = append(details, "user:"+user)
			}

			details = append(details, fmt.Sprintf("build:%s", c.flagId.FormatId(b.Build.Sequence, b.Build.Id)))

			if c.flagVerbose {
				for k, val := range b.Labels {
					if strings.HasPrefix(k, "waypoint/") {
						continue
					}

					if len(val) > 30 {
						val = val[:30] + "..."
					}

					extraDetails = append(extraDetails, fmt.Sprintf("artifact.%s:%s", k, val))
				}

				for k, val := range b.Build.Labels {
					if strings.HasPrefix(k, "waypoint/") {
						continue
					}

					if len(val) > 30 {
						val = val[:30] + "..."
					}

					extraDetails = append(extraDetails, fmt.Sprintf("build.%s:%s", k, val))
				}
				sort.Strings(extraDetails)
			}

			sort.Strings(details)

			table.Rich(
				[]string{
					status,
					c.flagId.FormatId(b.Sequence, b.Id),
					b.Component.Name,
					details[0],
					startTime,
					completeTime,
				}, []string{
					statusColor,
				},
			)

			if len(details[1:]) > 0 {
				for _, dr := range details[1:] {
					table.Rich([]string{"", "", "", dr}, nil)
				}
			}

			if len(extraDetails) > 0 {
				for _, dr := range extraDetails {
					table.Rich([]string{"", "", "", dr}, nil)
				}
			}
		}

		c.ui.Table(table)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ArtifactListCommand) displayJson(artifacts []*pb.PushedArtifact) error {
	var output []map[string]interface{}

	for _, art := range artifacts {
		i := map[string]interface{}{}

		i["id"] = art.Id
		i["sequence"] = art.Sequence
		i["application"] = art.Application
		i["labels"] = art.Labels
		i["component"] = art.Component.Name
		i["status"] = c.statusJson(art.Status)
		i["workspace"] = art.Workspace.Workspace
		i["build"] = c.buildJson(art.Build)

		output = append(output, i)
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	c.ui.Output(string(data))

	return nil
}

func (c *ArtifactListCommand) statusJson(status *pb.Status) interface{} {
	i := map[string]interface{}{}

	i["state"] = status.State.String()
	i["complete_time"] = status.CompleteTime.AsTime().Format(time.RFC3339Nano)
	i["start_time"] = status.StartTime.AsTime().Format(time.RFC3339Nano)

	return i
}

func (c *ArtifactListCommand) buildJson(b *pb.Build) interface{} {
	i := map[string]interface{}{}

	i["id"] = b.Id
	i["sequence"] = b.Sequence
	i["labels"] = b.Labels
	i["status"] = c.statusJson(b.Status)

	return i
}

func (c *ArtifactListCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "workspace-all",
			Target: &c.flagWorkspaceAll,
			Usage:  "List builds in all workspaces for this project and application.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "verbose",
			Aliases: []string{"V"},
			Target:  &c.flagVerbose,
			Usage:   "Display more details about each deployment.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage:  "Output the deployment information as JSON.",
		})

		initIdFormat(f, &c.flagId)
		initFilterFlags(set, &c.filterFlags, filterOptionOrder)
	})
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
	return formatHelp(`
Usage: waypoint artifact list [options]

  Lists the artifacts that are pushed to a registry. This does not
  list the artifacts that are just part of local builds.

` + c.Flags().Help())
}
