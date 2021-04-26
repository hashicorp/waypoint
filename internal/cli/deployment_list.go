package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serversort "github.com/hashicorp/waypoint/internal/server/sort"
)

type DeploymentListCommand struct {
	*baseCommand

	flagWorkspaceAll bool
	flagVerbose      bool
	flagJson         bool
	flagId           idFormat
	filterFlags      filterFlags
}

func shortImg(img string) string {
	if strings.HasPrefix(img, "sha256:") {
		return img[7:14]
	}

	return img[:7]
}

// Add either language: or languages: based on how many values are specified
func niceLanguages(langs string) string {
	parts := strings.Split(langs, ",")

	if len(parts) == 1 {
		return "language:" + strings.TrimSpace(parts[0])
	}

	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}

	return "languages:" + strings.Join(parts, ", ")
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

		phyState, err := c.filterFlags.physState()
		if err != nil {
			return err
		}

		// Get the latest release
		release, err := client.GetLatestRelease(ctx, &pb.GetLatestReleaseRequest{
			Application: app.Ref(),
			Workspace:   c.project.WorkspaceRef(),
		})
		if status.Code(err) == codes.NotFound {
			err = nil
			release = nil
		}
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// List builds
		resp, err := client.ListDeployments(c.Ctx, &pb.ListDeploymentsRequest{
			Application:   app.Ref(),
			Workspace:     wsRef,
			PhysicalState: phyState,
			Status:        c.filterFlags.statusFilters(),
			Order:         c.filterFlags.orderOp(),
			LoadDetails:   pb.Deployment_BUILD,
		})
		if err != nil {
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		sort.Sort(serversort.DeploymentCompleteDesc(resp.Deployments))

		// get status reports
		statusReportsResp, err := client.ListStatusReports(ctx, &pb.ListStatusReportsRequest{
			Application: app.Ref(),
			Workspace:   wsRef,
		})

		if status.Code(err) == codes.NotFound || status.Code(err) == codes.Unimplemented {
			err = nil
			statusReportsResp = nil
		}
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		if c.flagJson {
			return c.displayJson(resp.Deployments)
		}

		tbl := terminal.NewTable("", "ID", "Platform", "Details", "Started", "Completed", "URL", "Health")

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
				switch b.State {
				case pb.Operation_DESTROYED:
					status = bullet
				case pb.Operation_CREATED:
					status = "✔"
					statusColor = terminal.Green

					if release != nil && release.DeploymentId == b.Id {
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
			if t, err := ptypes.Timestamp(b.Status.StartTime); err == nil {
				startTime = humanize.Time(t)
			}
			if t, err := ptypes.Timestamp(b.Status.CompleteTime); err == nil {
				completeTime = humanize.Time(t)
			}

			// Add status report information if we have any
			statusReportComplete := "n/a"
			for _, statusReport := range statusReportsResp.StatusReports {
				if deploymentTargetId, ok := statusReport.TargetId.(*pb.StatusReport_DeploymentId); ok {
					if deploymentTargetId.DeploymentId == b.Id {
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

						if t, err := ptypes.Timestamp(statusReport.GeneratedTime); err == nil {
							statusReportComplete = fmt.Sprintf("%s - %s", statusReportComplete, humanize.Time(t))
						}
					}
				}
			}

			var (
				extraDetails []string
				details      []string
			)

			if user, ok := b.Labels["common/user"]; ok {
				details = append(details, "user:"+user)
			} else if user, ok := b.Preload.Build.Labels["common/user"]; ok {
				details = append(details, "build-user:"+user)
			}

			if bp, ok := b.Preload.Build.Labels["common/languages"]; ok {
				details = append(details, niceLanguages(bp))
			}

			if img, ok := b.Preload.Build.Labels["common/image-id"]; ok {
				img = shortImg(img)

				details = append(details, "image:"+img)
			}

			artdetails := fmt.Sprintf("artifact:%s", c.flagId.FormatId(b.Preload.Artifact.Sequence, b.Preload.Artifact.Id))
			if len(details) == 0 {
				details = append(details, artdetails)
			} else if c.flagVerbose {
				details = append(details,
					artdetails,
					fmt.Sprintf("build:%s", c.flagId.FormatId(b.Preload.Build.Sequence, b.Preload.Build.Id)))
			}

			if c.flagVerbose {
				for k, val := range b.Labels {
					if strings.HasPrefix(k, "waypoint/") {
						continue
					}

					if len(val) > 30 {
						val = val[:30] + "..."
					}

					extraDetails = append(extraDetails, fmt.Sprintf("deployment.%s:%s", k, val))
				}

				for k, val := range b.Preload.Artifact.Labels {
					if strings.HasPrefix(k, "waypoint/") {
						continue
					}

					if len(val) > 30 {
						val = val[:30] + "..."
					}

					extraDetails = append(extraDetails, fmt.Sprintf("artifact.%s:%s", k, val))
				}

				for k, val := range b.Preload.Build.Labels {
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

			tbl.Rich(
				[]string{
					status,
					c.flagId.FormatId(b.Sequence, b.Id),
					b.Component.Name,
					details[0],
					startTime,
					completeTime,
					b.Preload.DeployUrl,
					statusReportComplete,
				},
				[]string{
					statusColor,
				},
			)

			if len(details[1:]) > 0 {
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

		c.ui.Table(tbl)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *DeploymentListCommand) displayJson(deployments []*pb.Deployment) error {
	var output []map[string]interface{}

	for _, dep := range deployments {
		i := map[string]interface{}{}

		i["id"] = dep.Id
		i["sequence"] = dep.Sequence
		i["application"] = dep.Application
		i["labels"] = dep.Labels
		i["component"] = dep.Component.Name
		i["physical_state"] = dep.State.String()
		i["status"] = c.statusJson(dep.Status)
		i["workspace"] = dep.Workspace.Workspace
		i["artifact"] = c.artifactJson(dep.Preload.Artifact)
		i["build"] = c.buildJson(dep.Preload.Build)

		output = append(output, i)
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	c.ui.Output(string(data))

	return nil
}

func (c *DeploymentListCommand) artifactJson(art *pb.PushedArtifact) interface{} {
	i := map[string]interface{}{}

	i["id"] = art.Id
	i["sequence"] = art.Sequence
	i["labels"] = art.Labels
	i["status"] = c.statusJson(art.Status)

	return i
}

func (c *DeploymentListCommand) statusJson(status *pb.Status) interface{} {
	i := map[string]interface{}{}

	i["state"] = status.State.String()
	i["complete_time"] = status.CompleteTime.AsTime().Format(time.RFC3339Nano)
	i["start_time"] = status.StartTime.AsTime().Format(time.RFC3339Nano)

	return i
}

func (c *DeploymentListCommand) buildJson(b *pb.Build) interface{} {
	i := map[string]interface{}{}

	i["id"] = b.Id
	i["sequence"] = b.Sequence
	i["labels"] = b.Labels
	i["status"] = c.statusJson(b.Status)

	return i
}

func (c *DeploymentListCommand) Flags() *flag.Sets {
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
		initFilterFlags(set, &c.filterFlags, fillterOptionAll)
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
	return formatHelp(`
Usage: waypoint deployment list [options]

  Lists the deployments that were created.

` + c.Flags().Help())
}
