// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/waypoint/version"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ReleaseListCommand struct {
	*baseCommand

	flagWorkspaceAll bool
	flagVerbose      bool
	flagUrl          bool
	flagJson         bool
	flagId           idFormat
	filterFlags      filterFlags
}

func (c *ReleaseListCommand) Run(args []string) int {
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
		resp, err := client.UI_ListReleases(c.Ctx, &pb.UI_ListReleasesRequest{
			Application: app.Ref(),
			Workspace:   wsRef,
			Order:       c.filterFlags.orderOp(),
		})

		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.Unimplemented {
					var serverVersion string
					serverVersionResp := c.project.ServerVersion()
					if serverVersionResp != nil {
						serverVersion = serverVersionResp.Version
					}

					var clientVersion string
					clientVersionResp := version.GetVersion()
					if clientVersionResp != nil {
						clientVersion = clientVersionResp.Version
					}

					c.project.UI.Output(
						fmt.Sprintf("This CLI version %q is incompatible with the current server %q - missing UI_ListReleases method. Upgrade your server to v0.5.0 or higher or downgrade your CLI to v0.4 or older.", clientVersion, serverVersion),
						terminal.WithErrorStyle(),
					)
					return ErrSentinel
				}
			}
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		if c.flagJson {
			return c.displayJson(resp.Releases)
		}

		headers := []string{
			"", "ID", "Deployment ID", "Platform", "Details", "Started", "Completed", "Health",
		}

		if c.flagUrl {
			headers = append(headers, "URL")
		}

		tbl := terminal.NewTable(headers...)

		const bullet = "●"

		for _, releaseBundle := range resp.Releases {
			b := releaseBundle.Release

			// Determine our bullet
			status := ""
			statusColor := ""
			switch b.Status.State {
			case pb.Status_RUNNING:
				status = bullet
				statusColor = terminal.Yellow

			case pb.Status_SUCCESS:
				switch b.State {
				case pb.Operation_DESTROYED:
					status = bullet
				case pb.Operation_CREATED:
					status = "✔"
					statusColor = terminal.Green

					if resp.Releases[0] != nil && resp.Releases[0].Release.Id == b.Id {
						status = "🚀"
					}

				default:
					status = "?"
					statusColor = terminal.Yellow
				}
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

			// Add status report information if we have any
			statusReportComplete := "n/a"
			if releaseBundle.LatestStatusReport != nil {
				statusReport := releaseBundle.LatestStatusReport
				switch statusReport.Health.HealthStatus {
				case "READY":
					statusReportComplete = "✔"
				case "ALIVE":
					statusReportComplete = "✔"
				case "DOWN":
					statusReportComplete = "✖"
				case "PARTIAL":
					statusReportComplete = "●"
				case "UNKNOWN":
					statusReportComplete = "?"
				}

				if statusReport.GeneratedTime != nil {
					t := statusReport.GeneratedTime.AsTime()
					statusReportComplete = fmt.Sprintf("%s - %s", statusReportComplete, humanize.Time(t))
				}
			}

			var (
				extraDetails []string
				details      []string
			)

			if user, ok := b.Labels["common/user"]; ok {
				details = append(details, "user:"+user)
			} else if b.Preload.Build != nil {
				build := b.Preload.Build
				// labels have been set, safe to use them

				if user, ok := build.Labels["common/user"]; ok {
					details = append(details, "build-user:"+user)
				}
				if bp, ok := build.Labels["common/languages"]; ok {
					details = append(details, niceLanguages(bp))
				}

				if img, ok := build.Labels["common/image-id"]; ok {
					img, err = shortImg(img)
					if err != nil {
						app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
						return err
					}

					details = append(details, "image:"+img)
				}
			}

			j, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: releaseBundle.Release.JobId,
			})
			if err != nil {
				app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return err
			}
			if j.Pipeline != nil {
				pipeline := "name: " + j.Pipeline.PipelineName + ", run: " + strconv.FormatUint(j.Pipeline.RunSequence, 10) + ", step: " + j.Pipeline.Step
				details = append(details, pipeline)
			}

			if b.Preload.Artifact != nil {
				artDetails := fmt.Sprintf("artifact:%s", c.flagId.FormatId(b.Preload.Artifact.Sequence, b.Preload.Artifact.Id))
				if len(details) == 0 {
					details = append(details, artDetails)
				} else if c.flagVerbose && b.Preload.Build != nil {
					details = append(details,
						artDetails,
						fmt.Sprintf("build:%s", c.flagId.FormatId(b.Preload.Build.Sequence, b.Preload.Build.Id)))
				}
			}

			if c.flagVerbose {
				for k, val := range b.Labels {
					if strings.HasPrefix(k, "waypoint/") {
						continue
					}

					if len(val) > 30 {
						val = val[:30] + "..."
					}

					extraDetails = append(extraDetails, fmt.Sprintf("Release.%s:%s", k, val))
				}
				sort.Strings(extraDetails)
			}

			sort.Strings(details)
			var firstDetails string
			if len(details) > 0 {
				firstDetails = details[0]
			}

			var columns []string

			columns = []string{
				status,
				c.flagId.FormatId(b.Sequence, b.Id),
				c.flagId.FormatId(b.Preload.Deployment.Sequence, b.Id),
				b.Component.Name,
				firstDetails,
				startTime,
				completeTime,
				statusReportComplete,
			}

			if c.flagUrl {
				url := "n/a"
				if b.Url != "" {
					url = b.Url
				} else if releaseBundle.Release.Url != "" {
					url = releaseBundle.Release.Url
				}
				columns = append(columns, url)
			}

			// Omit Waypoint releases that didn't actually happen on the platform
			if !b.Unimplemented {
				tbl.Rich(
					columns,
					[]string{
						statusColor,
					},
				)

				if len(details) > 1 {
					for _, dr := range details[1:] {
						tbl.Rich([]string{"", "", "", dr}, nil)
					}
				}

				if len(extraDetails) > 0 {
					for _, dr := range extraDetails {
						tbl.Rich([]string{"", "", "", dr}, nil)
					}
				}
			}
		}

		c.ui.Table(tbl)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ReleaseListCommand) displayJson(releases []*pb.UI_ReleaseBundle) error {
	var output []map[string]interface{}

	for _, rel := range releases {
		if rel.Release.Unimplemented {
			continue
		}

		i := map[string]interface{}{}

		i["id"] = rel.Release.Sequence
		i["deploymentId"] = rel.Release.Preload.Deployment.Sequence
		i["application"] = rel.Release.Application
		i["workspace"] = rel.Release.Workspace.Workspace
		i["url"] = rel.Release.Url
		i["labels"] = rel.Release.Labels
		i["component"] = rel.Release.Component.Name
		i["status"] = c.statusJson(rel.Release.Status)
		i["latestStatusReport"] = rel.LatestStatusReport
		i["preloadDetails"] = rel.Release.Preload

		output = append(output, i)
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	c.ui.Output(string(data))

	return nil
}

func (c *ReleaseListCommand) statusJson(status *pb.Status) interface{} {
	if status == nil {
		return nil
	}
	i := map[string]interface{}{}

	i["state"] = status.State.String()
	i["complete_time"] = status.CompleteTime.AsTime().Format(time.RFC3339Nano)
	i["start_time"] = status.StartTime.AsTime().Format(time.RFC3339Nano)

	return i
}

func (c *ReleaseListCommand) Flags() *flag.Sets {
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
			Usage:   "Display more details about each release.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "url",
			Aliases: []string{"u"},
			Target:  &c.flagUrl,
			Usage:   "Display release URL.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage:  "Output the release information as JSON.",
		})

		initIdFormat(f, &c.flagId)
		initFilterFlags(set, &c.filterFlags, filterOptionAll)
	})
}

func (c *ReleaseListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ReleaseListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ReleaseListCommand) Synopsis() string {
	return "List releases."
}

func (c *ReleaseListCommand) Help() string {
	return formatHelp(`
Usage: waypoint release list [options]

  Lists the releases that were created if the platform includes a releaser.

` + c.Flags().Help())
}
