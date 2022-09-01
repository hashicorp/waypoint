package cli

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/hashicorp/waypoint/internal/version"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type DeploymentListCommand struct {
	*baseCommand

	flagWorkspaceAll bool
	flagVerbose      bool
	flagUrl          bool
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
		resp, err := client.UI_ListDeployments(c.Ctx, &pb.UI_ListDeploymentsRequest{
			Application:   app.Ref(),
			Workspace:     wsRef,
			PhysicalState: phyState,
			Status:        c.filterFlags.statusFilters(),
			Order:         c.filterFlags.orderOp(),
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
						fmt.Sprintf("This CLI version %q is incompatible with the current server %q - missing UI_ListDeployments method. Upgrade your server to v0.5.0 or higher or downgrade your CLI to v0.4 or older.", clientVersion, serverVersion),
						terminal.WithErrorStyle(),
					)
					return ErrSentinel
				}
			}
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if len(resp.Deployments) == 0 {
			c.project.UI.Output(
				"No deployments found for application %q",
				app.Ref().Application,
				terminal.WithWarningStyle(),
			)
			return nil
		}

		if c.flagJson {
			return c.displayJson(resp.Deployments)
		}

		headers := []string{
			"", "ID", "Platform", "Details", "Started", "Completed", "Health",
		}

		if c.flagUrl {
			headers = append(headers, "URL")
		}

		tbl := terminal.NewTable(headers...)

		const bullet = "●"

		for _, deployBundle := range resp.Deployments {
			b := deployBundle.Deployment
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
			if b.Status.StartTime != nil {
				startTime = humanize.Time(b.Status.StartTime.AsTime())
			}

			if b.Status.CompleteTime != nil {
				completeTime = humanize.Time(b.Status.CompleteTime.AsTime())
			}

			// Add status report information if we have any
			statusReportComplete := "n/a"
			if deployBundle.LatestStatusReport != nil {
				statusReport := deployBundle.LatestStatusReport
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
			} else if deployBundle.Build != nil {
				build := deployBundle.Build
				// labels have been set, safe to use them

				if user, ok := build.Labels["common/user"]; ok {
					details = append(details, "build-user:"+user)
				}
				if bp, ok := build.Labels["common/languages"]; ok {
					details = append(details, niceLanguages(bp))
				}

				if img, ok := build.Labels["common/image-id"]; ok {
					img = shortImg(img)

					details = append(details, "image:"+img)
				}
			}

			j, err := c.project.Client().GetJob(c.Ctx, &pb.GetJobRequest{
				JobId: deployBundle.Deployment.JobId,
			})
			if err != nil {
				app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return err
			}
			if j.Pipeline != nil {
				pipeline := "pipeline: " + j.Pipeline.PipelineName + "[run: " + strconv.FormatUint(j.Pipeline.RunSequence, 10) + "]" + "[step: " + j.Pipeline.Step + "]"
				details = append(details, pipeline)
			}

			if deployBundle.Artifact != nil {
				artDetails := fmt.Sprintf("artifact:%s", c.flagId.FormatId(deployBundle.Artifact.Sequence, deployBundle.Artifact.Id))
				if len(details) == 0 {
					details = append(details, artDetails)
				} else if c.flagVerbose && deployBundle.Build != nil {
					details = append(details,
						artDetails,
						fmt.Sprintf("build:%s", c.flagId.FormatId(deployBundle.Build.Sequence, deployBundle.Build.Id)))
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

					extraDetails = append(extraDetails, fmt.Sprintf("deployment.%s:%s", k, val))
				}

				if deployBundle.Artifact != nil {
					for k, val := range deployBundle.Artifact.Labels {
						if strings.HasPrefix(k, "waypoint/") {
							continue
						}

						if len(val) > 30 {
							val = val[:30] + "..."
						}

						extraDetails = append(extraDetails, fmt.Sprintf("artifact.%s:%s", k, val))
					}
				}

				if deployBundle.Build != nil {
					for k, val := range deployBundle.Build.Labels {
						if strings.HasPrefix(k, "waypoint/") {
							continue
						}

						if len(val) > 30 {
							val = val[:30] + "..."
						}

						extraDetails = append(extraDetails, fmt.Sprintf("build.%s:%s", k, val))
					}
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
				} else if deployBundle.DeployUrl != "" {
					url = deployBundle.DeployUrl
				}
				columns = append(columns, url)
			}

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

		c.ui.Table(tbl)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *DeploymentListCommand) displayJson(deployments []*pb.UI_DeploymentBundle) error {
	var output []map[string]interface{}

	for _, dep := range deployments {
		i := map[string]interface{}{}

		i["id"] = dep.Deployment.Id
		i["sequence"] = dep.Deployment.Sequence
		i["application"] = dep.Deployment.Application
		i["labels"] = dep.Deployment.Labels
		i["component"] = dep.Deployment.Component.Name
		i["physical_state"] = dep.Deployment.State.String()
		i["status"] = c.statusJson(dep.Deployment.Status)
		i["workspace"] = dep.Deployment.Workspace.Workspace
		i["artifact"] = c.artifactJson(dep.Artifact)
		i["build"] = c.buildJson(dep.Build)
		i["latestStatusReport"] = dep.LatestStatusReport

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
	if art == nil {
		return nil
	}

	i := map[string]interface{}{}

	i["id"] = art.Id
	i["sequence"] = art.Sequence
	i["labels"] = art.Labels
	i["status"] = c.statusJson(art.Status)

	return i
}

func (c *DeploymentListCommand) statusJson(status *pb.Status) interface{} {
	if status == nil {
		return nil
	}
	i := map[string]interface{}{}

	i["state"] = status.State.String()
	i["complete_time"] = status.CompleteTime.AsTime().Format(time.RFC3339Nano)
	i["start_time"] = status.StartTime.AsTime().Format(time.RFC3339Nano)

	return i
}

func (c *DeploymentListCommand) buildJson(b *pb.Build) interface{} {
	if b == nil {
		return nil
	}
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
			Name:    "url",
			Aliases: []string{"u"},
			Target:  &c.flagUrl,
			Usage:   "Display deployment URL.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage:  "Output the deployment information as JSON.",
		})

		initIdFormat(f, &c.flagId)
		initFilterFlags(set, &c.filterFlags, filterOptionAll)
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
Usage: waypoint deployment list [options] [project/app]

  Lists the deployments that were created.

` + c.Flags().Help())
}
