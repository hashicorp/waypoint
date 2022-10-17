package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/version"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type StatusCommand struct {
	*baseCommand

	flagContextName      string
	flagVerbose          bool
	flagJson             bool
	flagAllProjects      bool
	flagRefreshAppStatus bool

	serverCtx *clicontext.Config
}

func (c *StatusCommand) Run(args []string) int {
	flagSet := c.Flags()
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
	); err != nil {
		return 1
	}

	var ctxName string
	defaultName, err := c.contextStorage.Default()
	if err != nil {
		c.ui.Output(
			"Error getting default context: %s",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	ctxName = defaultName

	if ctxName != "" {
		ctxConfig, err := c.contextStorage.Load(ctxName)
		if err != nil {
			c.ui.Output("Error loading context %q: %s", ctxName, err.Error(), terminal.WithErrorStyle())
			return 1
		}
		c.serverCtx = ctxConfig
	} else {
		c.ui.Output(wpNoServerContext, terminal.WithWarningStyle())
	}

	// Optionally refresh status
	if c.flagRefreshAppStatus {
		if err := c.RefreshApplicationStatus(); err != nil {
			c.ui.Output("CLI failed to refresh project statuses: "+clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		c.ui.Output("")
	}

	// Generate a status view
	if c.refProject == nil || c.flagAllProjects {
		// Show high-level status of all projects
		err = c.FormatProjectStatus()
		if err != nil {
			c.ui.Output("CLI failed to build project statuses: "+clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	} else if c.refProject != nil && c.flagApp == "" {
		// Show status of apps inside project
		projectTarget := c.refProject.Project
		err = c.FormatProjectAppStatus(projectTarget)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				var serverAddress string
				if c.serverCtx != nil {
					serverAddress = c.serverCtx.Server.Address
				}

				c.ui.Output(wpProjectNotFound, projectTarget, serverAddress, terminal.WithErrorStyle())
			} else {
				c.ui.Output("CLI failed to format project app statuses:"+clierrors.Humanize(err), terminal.WithErrorStyle())
			}
			return 1
		}
	} else if c.refProject != nil && c.flagApp != "" {
		// Advanced view of a single app status
		projectTarget := c.refProject.Project
		appTarget := c.flagApp
		err = c.FormatAppStatus(projectTarget, appTarget)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				var serverAddress string
				if c.serverCtx != nil {
					serverAddress = c.serverCtx.Server.Address
				}

				c.ui.Output(wpAppNotFound, appTarget, projectTarget, serverAddress, terminal.WithErrorStyle())
			} else {
				c.ui.Output("CLI failed to format app status:"+clierrors.Humanize(err), terminal.WithErrorStyle())
			}
			return 1
		}
	}

	return 0
}

// RefreshApplicationStatus takes a project and application target and generates
// a list of applications to refresh the status on. If all projects are requested
// to be refreshed, the CLI will do its best to honor the request. However, if
// a project is local and the CLI was not invoked inside that project dir, the
// CLI won't be able to refresh that project's application statuses.
func (c *StatusCommand) RefreshApplicationStatus() error {
	// Get our API client
	client := c.project.Client()

	// Get the entire list of apps
	// Determine project locality
	// Do the Work (local or remote)
	if c.flagAllProjects {
		c.ui.Output("This command does not support refreshing statuses for all "+
			"defined projects in Waypoint. Use the project argument to narrow down "+
			"which projects you hope to refresh a status on.", terminal.WithWarningStyle())
		return nil
	}

	if c.refProject == nil {
		return errors.New("No project specified - use the -project flag")
	}

	if len(c.refApps) == 0 {
		// Technically the user could have no apps in the
		return errors.New("No apps specified - please use the -app flag, or run from within " +
			"a project directory that contains apps.")
	}

	projectResp, err := client.GetProject(c.Ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: c.refProject.Project,
		},
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			var serverAddress string
			if c.serverCtx != nil {
				serverAddress = c.serverCtx.Server.Address
			}

			c.ui.Output(wpProjectNotFound, c.refProject.Project, serverAddress, terminal.WithErrorStyle())
		}
		return err
	}
	project := projectResp.Project

	workspace, err := c.getWorkspaceFromProject(projectResp)
	if err != nil {
		return err
	}

	// Only refresh specified applications
	var appsToRefresh []*pb.Application
	for _, refApp := range c.refApps {
		for _, app := range project.Applications {
			if app.Name == refApp.Application {
				appsToRefresh = append(appsToRefresh, app)
			}
		}
	}

	// Useful for printing
	var appNames []string
	for _, refApp := range c.refApps {
		appNames = append(appNames, refApp.Application)
	}

	if len(appsToRefresh) == 0 {
		// Corner case - will happen if they typo an app given to the -app flag.
		c.ui.Output("Specified app(s) %q not found in project %q",
			strings.Join(appNames, ", "),
			project.Name,
			clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}

	if err = c.doRefresh(project, appsToRefresh, workspace); err != nil {
		c.ui.Output("Failed to refresh app(s) %q status in project %q: %s",
			strings.Join(appNames, ", "),
			project.Name,
			clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}

	return nil
}

// doRefresh takes a list of applications and their parent project and will
// attempt to run the statusfunc configured for each requested application.
func (c *StatusCommand) doRefresh(
	project *pb.Project,
	appList []*pb.Application,
	workspace string,
) error {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("")
	c.ui.Output("")
	for _, app := range appList {
		s.Update("Refreshing status for app %q in project %q. Refreshing could take a while...",
			app.Name, project.Name)

		err := c.refreshAppStatus(project, app, workspace)
		if err != nil {
			s.Update("Failed to refresh app status\n")
			s.Status(terminal.StatusError)
			s.Done()

			c.ui.Output(
				"Error attempting to refresh application %q: %s",
				app.Name,
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return err
		}
	}

	s.Update("Finished refreshing app statuses in project %q", project.Name)
	s.Done()
	c.ui.Output("")

	return nil
}

// refreshAppStatus takes a project and a single application and uses
// the local runner to execute a StatusReport refresh function by getting
// the latest Deployment or Release and using the application client to
// refresh the apps Status Report.
func (c *StatusCommand) refreshAppStatus(
	project *pb.Project,
	app *pb.Application,
	workspace string,
) error {
	// We must setup app ref so DoApp has the context for which app to execute
	// its operations on.
	c.refApps = []*pb.Ref_Application{{
		Application: app.Name,
		Project:     project.Name,
	}}

	err := c.DoApp(c.Ctx, func(ctx context.Context, appClient *clientpkg.App) error {
		// Get our API client
		client := c.project.Client()

		// Get Latest Deployment and Release

		// Deployments
		deploymentsResp, err := ListDeployments(c.Ctx, client, app.Name, project.Name, workspace)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.Unimplemented {
					return c.errAPIUnimplemented(err)
				}
			}
			return err
		}

		var deployment *pb.Deployment
		if deploymentsResp.Deployments != nil && len(deploymentsResp.Deployments) > 0 {
			deployment = deploymentsResp.Deployments[0].Deployment
		} else {
			c.ui.Output("No deployments in project %q for app %q to refresh a status on!",
				project.Name,
				app.Name,
				terminal.WithWarningStyle())
			// We probably don't have a release either, so return early
			return nil
		}

		_, err = appClient.StatusReport(c.Ctx, &pb.Job_StatusReportOp{
			Target: &pb.Job_StatusReportOp_Deployment{
				Deployment: deployment,
			},
		})
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return err
		}

		// Releases
		releaseResp, err := ListReleases(c.Ctx, client, app.Name, project.Name, workspace)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.Unimplemented {
					return c.errAPIUnimplemented(err)
				}
			}
			return err
		}

		var release *pb.Release
		if releaseResp.Releases != nil && len(releaseResp.Releases) > 0 {
			release = releaseResp.Releases[0].Release
		}

		if release != nil && !release.Unimplemented {
			_, err = appClient.StatusReport(c.Ctx, &pb.Job_StatusReportOp{
				Target: &pb.Job_StatusReportOp_Release{
					Release: release,
				},
			})
			if err != nil {
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return err
			}
		}

		return nil
	})
	if err != nil {
		if err != ErrSentinel {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		}

		return err
	}

	return nil
}

//
// Helper functions
//

func ListDeployments(c context.Context, client pb.WaypointClient, app string, project string, workspace string) (*pb.UI_ListDeploymentsResponse, error) {
	resp, err := client.UI_ListDeployments(c, &pb.UI_ListDeploymentsRequest{
		Application: &pb.Ref_Application{
			Application: app,
			Project:     project,
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: workspace,
		},
		PhysicalState: pb.Operation_CREATED,
		Order: &pb.OperationOrder{
			Order: pb.OperationOrder_COMPLETE_TIME,
			Limit: 1,
		},
	})
	return resp, err
}

func ListReleases(c context.Context, client pb.WaypointClient, app string, project string, workspace string) (*pb.UI_ListReleasesResponse, error) {
	resp, err := client.UI_ListReleases(c, &pb.UI_ListReleasesRequest{
		Application: &pb.Ref_Application{
			Application: app,
			Project:     project,
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: workspace,
		},
		Order: &pb.OperationOrder{
			Order: pb.OperationOrder_COMPLETE_TIME,
			Limit: 1,
		},
	})
	return resp, err
}

//
// UI format functions
//

// FormatProjectAppStatus formats all applications inside a project
func (c *StatusCommand) FormatProjectAppStatus(projectTarget string) error {
	if !c.flagJson && c.serverCtx != nil {
		c.ui.Output(wpStatusProjectMsg, projectTarget, c.serverCtx.Server.Address)
	}

	// Get our API client
	client := c.project.Client()

	resp, err := client.GetProject(c.Ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: projectTarget,
		},
	})
	if err != nil {
		return err
	}
	project := resp.Project

	workspace, err := c.getWorkspaceFromProject(resp)
	if err != nil {
		return err
	}

	// Summary
	//   App list

	appHeaders := []string{
		"App", "Workspace", "Deployment Status", "Deployment Checked", "Release Status", "Release Checked",
	}

	appTbl := terminal.NewTable(appHeaders...)

	appFailures := false
	for _, app := range project.Applications {
		// Get the latest deployment
		deploymentsResp, err := ListDeployments(c.Ctx, client, app.Name, project.Name, workspace)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.Unimplemented {
					return c.errAPIUnimplemented(err)
				}
			}
			return err
		}

		var appDeployStatus *pb.StatusReport
		if deploymentsResp.Deployments != nil && len(deploymentsResp.Deployments) > 0 {
			appDeployStatus = deploymentsResp.Deployments[0].LatestStatusReport
		}

		statusReportComplete, statusReportCheckTime, err := c.FormatStatusReportComplete(appDeployStatus)
		if err != nil {
			return err
		}

		if appDeployStatus != nil && appDeployStatus.Health != nil {
			if appDeployStatus.Health.HealthStatus == "ERROR" ||
				appDeployStatus.Health.HealthStatus == "DOWN" {
				appFailures = true
			}
		}

		// Get the latest release, if there was one
		releasesResp, err := ListReleases(c.Ctx, client, app.Name, project.Name, workspace)
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.Unimplemented {
					return c.errAPIUnimplemented(err)
				}
			}
			return err
		}
		var appReleaseStatus *pb.StatusReport
		if releasesResp.Releases != nil && len(releasesResp.Releases) > 0 {
			appReleaseStatus = releasesResp.Releases[0].LatestStatusReport
		}

		statusReportCompleteRelease, statusReportCheckTimeRelease, err := c.FormatStatusReportComplete(appReleaseStatus)
		if err != nil {
			return err
		}

		if appReleaseStatus != nil && appReleaseStatus.Health != nil {
			if appReleaseStatus.Health.HealthStatus == "ERROR" ||
				appReleaseStatus.Health.HealthStatus == "DOWN" {
				appFailures = true
			}
		}

		statusColor := ""
		columns := []string{
			app.Name,
			workspace,
			statusReportComplete,
			statusReportCheckTime,
			statusReportCompleteRelease,
			statusReportCheckTimeRelease,
		}

		// Add column data to table
		appTbl.Rich(
			columns,
			[]string{
				statusColor,
			},
		)
	}

	if c.flagJson {
		c.outputJsonProjectAppStatus(appTbl, project)
	} else {
		c.ui.Output("")
		c.ui.Table(appTbl, terminal.WithStyle("Simple"))
		c.ui.Output("")
		c.ui.Output(wpStatusProjectSuccessMsg)
	}

	if appFailures {
		c.ui.Output("")

		c.ui.Output(wpStatusHealthTriageMsg, projectTarget, terminal.WithWarningStyle())
	}

	return nil
}

func (c *StatusCommand) FormatAppStatus(projectTarget string, appTarget string) error {
	if !c.flagJson && c.serverCtx != nil {
		c.ui.Output(wpStatusAppProjectMsg, appTarget, projectTarget, c.serverCtx.Server.Address)
	}

	// Get our API client
	client := c.project.Client()

	projResp, err := client.GetProject(c.Ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: projectTarget,
		},
	})
	if err != nil {
		return err
	}
	project := projResp.Project

	workspace, err := c.getWorkspaceFromProject(projResp)
	if err != nil {
		return err
	}

	// App Summary
	//  Summary of single app
	var app *pb.Application
	for _, a := range project.Applications {
		if a.Name == appTarget {
			app = a
			break
		}
	}
	if app == nil {
		return fmt.Errorf("Did not find application %q in project %q", appTarget, projectTarget)
	}

	// Deployment Summary
	//   Deployment List

	// Get the latest deployment
	respDeployList, err := ListDeployments(c.Ctx, client, app.Name, project.Name, workspace)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.Unimplemented {
				return c.errAPIUnimplemented(err)
			}
		}
		return err
	}

	// deployment and releases use the same headers, with the exception that
	// deployment lists an additional item for instances count
	releaseHeaders := []string{
		"App Name", "Version", "Workspace", "Platform", "Artifact", "Lifecycle State",
	}

	// Add "Instance Count" for deployment summary headers
	deployHeaders := append(releaseHeaders, "Instances Count")

	deployTbl := terminal.NewTable(deployHeaders...)

	resourcesHeaders := []string{
		"Type", "Name", "Platform", "Health", "Time Created",
	}

	resourcesTbl := terminal.NewTable(resourcesHeaders...)

	deployStatusReportComplete := "N/A"
	var deployStatusReportCheckTime string
	appFailures := false
	if len(respDeployList.Deployments) > 0 {
		deployBundle := respDeployList.Deployments[0]
		deploy := deployBundle.Deployment
		appDeployStatus := deployBundle.LatestStatusReport

		var instancesCount uint32
		if appDeployStatus != nil {
			instancesCount = appDeployStatus.InstancesCount
		}
		statusColor := ""

		var details string
		if deployBundle.Build != nil {
			if deployBundle.Artifact != nil {
				artDetails := fmt.Sprintf("id:%d", deployBundle.Artifact.Sequence)
				details = artDetails
			}
			if img, ok := deployBundle.Build.Labels["common/image-id"]; ok {
				img = shortImg(img)

				details = details + " image:" + img
			}
		}

		columns := []string{
			deploy.Application.Application,
			fmt.Sprintf("v%d", deploy.Sequence),
			deploy.Workspace.Workspace,
			deploy.Component.Name,
			details,
			deploy.Status.State.String(),
			fmt.Sprintf("%d", instancesCount),
		}

		// Add column data to table
		deployTbl.Rich(
			columns,
			[]string{
				statusColor,
			},
		)

		deployStatusReportComplete, deployStatusReportCheckTime, err = c.FormatStatusReportComplete(appDeployStatus)
		if err != nil {
			return err
		}

		// Deployment Resources Summary
		//   Resources List
		if appDeployStatus != nil {
			if appDeployStatus.Health.HealthStatus == "ERROR" ||
				appDeployStatus.Health.HealthStatus == "DOWN" {
				appFailures = true
			}

			for _, dr := range appDeployStatus.Resources {
				var createdTime string
				if dr.CreatedTime != nil {
					createdTime = humanize.Time(dr.CreatedTime.AsTime())
				}

				columns := []string{
					dr.Type,
					dr.Name,
					dr.Platform,
					dr.Health.String(),
					createdTime,
				}

				// Add column data to table
				resourcesTbl.Rich(
					columns,
					[]string{
						statusColor,
					},
				)
			}
		}

	} // else show no table

	// Release Summary
	//   Release List

	releasesResp, err := ListReleases(c.Ctx, client, app.Name, project.Name, workspace)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.Unimplemented {
				return c.errAPIUnimplemented(err)
			}
		}
		return err
	}

	// Same headers as deploy
	releaseTbl := terminal.NewTable(releaseHeaders...)
	releaseResourcesTbl := terminal.NewTable(resourcesHeaders...)

	releaseUnimplemented := true
	releaseStatusReportComplete := "N/A"
	var releaseStatusReportCheckTime string
	if releasesResp.Releases != nil {
		release := releasesResp.Releases[0].Release
		releaseUnimplemented = release.Unimplemented

		if !release.Unimplemented {
			appReleaseStatus := releasesResp.Releases[0].LatestStatusReport

			statusColor := ""

			var details string
			if release.Preload.Artifact != nil {
				artDetails := fmt.Sprintf("id:%d", release.Preload.Artifact.Sequence)
				details = artDetails
			}
			if img, ok := release.Preload.Build.Labels["common/image-id"]; ok {
				img = shortImg(img)

				details = details + " image:" + img
			}

			columns := []string{
				release.Application.Application,
				fmt.Sprintf("v%d", release.Sequence),
				release.Workspace.Workspace,
				release.Component.Name,
				details,
				release.Status.State.String(),
			}

			// Add column data to table
			releaseTbl.Rich(
				columns,
				[]string{
					statusColor,
				},
			)

			releaseStatusReportComplete, releaseStatusReportCheckTime, err = c.FormatStatusReportComplete(appReleaseStatus)
			if err != nil {
				return err
			}

			// Release Resources Summary
			//   Resources List
			if appReleaseStatus != nil {
				if appReleaseStatus.Health.HealthStatus == "ERROR" ||
					appReleaseStatus.Health.HealthStatus == "DOWN" {
					appFailures = true
				}

				for _, rr := range appReleaseStatus.Resources {
					var createdTime string
					if rr.CreatedTime != nil {
						createdTime = humanize.Time(rr.CreatedTime.AsTime())
					}

					columns := []string{
						rr.Type,
						rr.Name,
						rr.Platform,
						rr.Health.String(),
						createdTime,
					}

					// Add column data to table
					releaseResourcesTbl.Rich(
						columns,
						[]string{
							statusColor,
						},
					)
				}
			}

		}
	} // else show no table

	appHeaders := []string{
		"App", "Workspace", "Deployment Status", "Deployment Checked", "Release Status", "Release Checked",
	}

	appTbl := terminal.NewTable(appHeaders...)

	statusColor := ""
	columns := []string{
		app.Name,
		workspace,
		deployStatusReportComplete,
		deployStatusReportCheckTime,
		releaseStatusReportComplete,
		releaseStatusReportCheckTime,
	}

	// Add column data to table
	appTbl.Rich(
		columns,
		[]string{
			statusColor,
		},
	)

	// TODO(briancain): we don't yet store a list of recent events per app
	// but it would go here if we did.
	// Recent Events
	//   Events List

	if c.flagJson {
		c.outputJsonAppStatus(appTbl, deployTbl, resourcesTbl, releaseTbl, releaseResourcesTbl, project)
	} else {
		c.ui.Output("Application Summary", terminal.WithHeaderStyle())
		c.ui.Table(appTbl, terminal.WithStyle("Simple"))
		c.ui.Output("Deployment Summary", terminal.WithHeaderStyle())
		c.ui.Table(deployTbl, terminal.WithStyle("Simple"))
		c.ui.Output("Deployment Resources Summary", terminal.WithHeaderStyle())
		c.ui.Table(resourcesTbl, terminal.WithStyle("Simple"))

		if !releaseUnimplemented {
			c.ui.Output("Release Summary", terminal.WithHeaderStyle())
			c.ui.Table(releaseTbl, terminal.WithStyle("Simple"))
			c.ui.Output("Release Resources Summary", terminal.WithHeaderStyle())
			c.ui.Table(releaseResourcesTbl, terminal.WithStyle("Simple"))
			c.ui.Output("")
		}

		c.ui.Output(wpStatusAppSuccessMsg)
	}

	if appFailures {
		c.ui.Output(wpStatusHealthTriageMsg, projectTarget, terminal.WithWarningStyle())
	}

	return nil
}

// FormatProjectStatus formats all known projects into a table
func (c *StatusCommand) FormatProjectStatus() error {
	if !c.flagJson && c.serverCtx != nil {
		c.ui.Output(wpStatusMsg, c.serverCtx.Server.Address)
	}

	// Get our API client
	client := c.project.Client()

	projectResp, err := client.ListProjects(c.Ctx, &pb.ListProjectsRequest{PaginationOptions: &pb.PaginationRequest{}})
	if err != nil {
		c.ui.Output("Failed to retrieve all projects:"+clierrors.Humanize(err), terminal.WithErrorStyle())
		return err
	}
	projNameList := projectResp.Projects

	headers := []string{
		"Project", "Workspace", "Deployment Statuses", "Release Statuses",
	}

	tbl := terminal.NewTable(headers...)

	for _, projectRef := range projNameList {
		resp, err := client.GetProject(c.Ctx, &pb.GetProjectRequest{
			Project: projectRef,
		})
		if err != nil {
			return err
		}

		workspace, err := c.getWorkspaceFromProject(resp)
		if err != nil {
			return err
		}

		// Get App Statuses
		var appDeployStatusReports []*pb.StatusReport
		var appReleaseStatusReports []*pb.StatusReport
		for _, app := range resp.Project.Applications {
			// Latest Deployment for app
			respDeployList, err := ListDeployments(c.Ctx, client, app.Name, resp.Project.Name, workspace)
			if err != nil {
				if s, ok := status.FromError(err); ok {
					if s.Code() == codes.Unimplemented {
						return c.errAPIUnimplemented(err)
					}
				}
				return err
			}

			var appStatusReportDeploy *pb.StatusReport
			if respDeployList.Deployments != nil && len(respDeployList.Deployments) > 0 {
				appStatusReportDeploy = respDeployList.Deployments[0].LatestStatusReport

				if appStatusReportDeploy != nil {
					appDeployStatusReports = append(appDeployStatusReports, appStatusReportDeploy)
				}
			}

			// Latest Release for app
			respReleaseList, err := ListReleases(c.Ctx, client, app.Name, resp.Project.Name, workspace)
			if err != nil {
				if s, ok := status.FromError(err); ok {
					if s.Code() == codes.Unimplemented {
						return c.errAPIUnimplemented(err)
					}
				}
				return err
			}

			var appStatusReportRelease *pb.StatusReport
			if respReleaseList.Releases != nil && len(respReleaseList.Releases) > 0 {
				appStatusReportRelease = respReleaseList.Releases[0].LatestStatusReport

				if appStatusReportRelease != nil {
					appReleaseStatusReports = append(appReleaseStatusReports, appStatusReportRelease)
				}
			}
		}

		deployStatusReportComplete := c.buildAppStatus(appDeployStatusReports)
		releaseStatusReportComplete := c.buildAppStatus(appReleaseStatusReports)

		statusColor := ""
		columns := []string{
			resp.Project.Name,
			workspace,
			deployStatusReportComplete,
			releaseStatusReportComplete,
		}

		// Add column data to table
		tbl.Rich(
			columns,
			[]string{
				statusColor,
			},
		)
	}

	// Render the table
	if c.flagJson {
		c.outputJsonProjectStatus(tbl)
	} else {
		c.ui.Output("")
		c.ui.Table(tbl, terminal.WithStyle("Simple"))
		c.ui.Output("")
		c.ui.Output(wpStatusSuccessMsg)
	}

	return nil
}

func (c *StatusCommand) outputJsonProjectStatus(t *terminal.Table) error {
	output := make(map[string]interface{})

	// Add server context
	serverContext := map[string]interface{}{}

	var serverAddress, serverPlatform string
	if c.serverCtx != nil {
		serverAddress = c.serverCtx.Server.Address
		serverPlatform = c.serverCtx.Server.Platform
	}

	serverContext["Address"] = serverAddress
	serverContext["ServerPlatform"] = serverPlatform

	output["ServerContext"] = serverContext

	projects := c.formatJsonMap(t)
	output["Projects"] = projects

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	c.ui.Output(string(data))

	return nil
}

func (c *StatusCommand) outputJsonProjectAppStatus(
	t *terminal.Table,
	project *pb.Project,
) error {
	output := make(map[string]interface{})

	// Add server context
	serverContext := map[string]interface{}{}

	var serverAddress, serverPlatform string
	if c.serverCtx != nil {
		serverAddress = c.serverCtx.Server.Address
		serverPlatform = c.serverCtx.Server.Platform
	}

	serverContext["Address"] = serverAddress
	serverContext["ServerPlatform"] = serverPlatform

	output["ServerContext"] = serverContext

	// Add project info
	projectInfo := map[string]interface{}{}
	projectInfo["Name"] = project.Name

	output["Project"] = projectInfo

	app := c.formatJsonMap(t)
	output["Applications"] = app

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	c.ui.Output(string(data))

	return nil
}

func (c *StatusCommand) outputJsonAppStatus(
	appTbl *terminal.Table,
	deployTbl *terminal.Table,
	resourcesTbl *terminal.Table,
	releaseTbl *terminal.Table,
	releaseResourcesTbl *terminal.Table,
	project *pb.Project,
) error {
	output := make(map[string]interface{})

	// Add server context
	serverContext := map[string]interface{}{}

	var serverAddress, serverPlatform string
	if c.serverCtx != nil {
		serverAddress = c.serverCtx.Server.Address
		serverPlatform = c.serverCtx.Server.Platform
	}

	serverContext["Address"] = serverAddress
	serverContext["ServerPlatform"] = serverPlatform

	output["ServerContext"] = serverContext

	// Add project info
	projectInfo := map[string]interface{}{}
	projectInfo["Name"] = project.Name

	output["Project"] = projectInfo

	app := c.formatJsonMap(appTbl)
	output["Applications"] = app

	deploySummary := c.formatJsonMap(deployTbl)
	output["DeploymentSummary"] = deploySummary

	deployResourcesSummary := c.formatJsonMap(resourcesTbl)
	output["DeploymentResourcesSummary"] = deployResourcesSummary

	releasesSummary := c.formatJsonMap(releaseTbl)
	output["ReleasesSummary"] = releasesSummary

	releaseResourcesSummary := c.formatJsonMap(releaseResourcesTbl)
	output["ReleasesResourcesSummary"] = releaseResourcesSummary

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	c.ui.Output(string(data))

	return nil
}

//
// CLI Status Helpers
//

func (c *StatusCommand) FormatStatusReportComplete(
	statusReport *pb.StatusReport,
) (string, string, error) {
	statusReportComplete := "N/A"

	if statusReport == nil {
		return statusReportComplete, "", nil
	}

	switch statusReport.Health.HealthStatus {
	case "READY":
		statusReportComplete = "✔ READY"
	case "ALIVE":
		statusReportComplete = "✔ ALIVE"
	case "DOWN":
		statusReportComplete = "✖ DOWN"
	case "PARTIAL":
		statusReportComplete = "● PARTIAL"
	case "UNKNOWN":
		statusReportComplete = "? UNKNOWN"
	}

	reportTime := "(unknown)"
	if statusReport.GeneratedTime.IsValid() {
		reportTime = humanize.Time(statusReport.GeneratedTime.AsTime())
	}

	return statusReportComplete, reportTime, nil
}

func (c *StatusCommand) getWorkspaceFromProject(pr *pb.GetProjectResponse) (string, error) {
	var workspace string

	if len(pr.Workspaces) != 0 {
		wp, err := c.workspace()
		if err != nil {
			return "", err
		}
		if wp != "" {
			for _, ws := range pr.Workspaces {
				if ws.Workspace.Workspace == wp {
					workspace = ws.Workspace.Workspace
					break
				}
			}

			if workspace == "" {
				return "", fmt.Errorf("Failed to find project in requested workspace %q", wp)
			}
		} else {
			// No workspace flag specified, try the "first" one
			workspace = pr.Workspaces[0].Workspace.Workspace
		}
	}

	return workspace, nil
}

// buildAppStatus takes a list of Status Reports and builds a string
// that details each apps status in a human readable format.
func (c *StatusCommand) buildAppStatus(reports []*pb.StatusReport) string {
	var ready, alive, down, unknown int

	for _, sr := range reports {
		switch sr.Health.HealthStatus {
		case "DOWN":
			down++
		case "UNKNOWN":
			unknown++
		case "READY":
			ready++
		case "ALIVE":
			alive++
		}
	}

	var result string
	if ready > 0 {
		result = result + fmt.Sprintf("%v READY ", ready)
	}
	if alive > 0 {
		result = result + fmt.Sprintf("%v ALIVE ", alive)
	}
	if down > 0 {
		result = result + fmt.Sprintf("%v DOWN ", down)
	}
	if alive > 0 {
		result = result + fmt.Sprintf("%v UNKNOWN ", unknown)
	}

	if result == "" {
		result = "N/A"
	}

	return result
}

// Takes a terminal Table and formats it into a map of key values to be used
// for formatting a JSON output response
func (c *StatusCommand) formatJsonMap(t *terminal.Table) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, row := range t.Rows {
		c := map[string]interface{}{}

		for j, r := range row {
			// Remove any whitespacess in key
			header := strings.ReplaceAll(t.Headers[j], " ", "")
			c[header] = r.Value
		}
		result = append(result, c)
	}

	return result
}

func (c *StatusCommand) errAPIUnimplemented(err error) error {
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

	c.project.UI.Output("This CLI version %q is incompatible with the current "+
		"server %q - missing API method. Upgrade your server to v0.5.0 or higher: %s",
		clierrors.Humanize(err),
		clientVersion, serverVersion,
		terminal.WithErrorStyle(),
	)
	return err
}

func (c *StatusCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:    "verbose",
			Aliases: []string{"V"},
			Target:  &c.flagVerbose,
			Usage:   "Display more details.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage:  "Output the status information as JSON.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "all-projects",
			Target: &c.flagAllProjects,
			Usage:  "Output status about every project in a workspace.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "refresh",
			Target: &c.flagRefreshAppStatus,
			Usage:  "Refresh application status for the requested app or apps in a project.",
		})
	})
}

func (c *StatusCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *StatusCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *StatusCommand) Synopsis() string {
	return "List and refresh application statuses."
}

func (c *StatusCommand) Help() string {
	return formatHelp(`
Usage: waypoint status [options] [project]

  View the current status of projects, applications, and their resources
  managed by Waypoint.

  When the '-refresh' flag is included, this command will attempt to regenerate
  every requested application's status report on-demand for both local and remote
  data sourced projects.

` + c.Flags().Help())
}

var (
	// Success or info messages

	wpStatusSuccessMsg = strings.TrimSpace(`
The projects listed above represent their current state known
in the Waypoint server. For more information about a project’s applications and
their current state, run ‘waypoint status -project=PROJECT-NAME’.
`)

	wpStatusProjectSuccessMsg = strings.TrimSpace(`
The project and its apps listed above represents its current state known
in the Waypoint server. For more information about a project’s applications and
their current state, run ‘waypoint status -app=APP-NAME -project=PROJECT-NAME’.
`)

	wpStatusAppSuccessMsg = strings.TrimSpace(`
The application and its resources listed above represents its current state known
in the Waypoint server.
`)

	wpStatusMsg = "Current project statuses in server context %q"

	wpStatusProjectMsg = "Current status for project %q in server context %q."

	wpStatusAppProjectMsg = strings.TrimSpace(`
Current status for application % q in project %q in server context %q.
`)

	// Failure messages

	wpNoServerContext = strings.TrimSpace(`
No default server context set for the Waypoint CLI. To set a default, use
'waypoint context use <context-name>'. To see a full list of known contexts,
run 'waypoint context list'. If Waypoint is running in local mode, this is expected.
`)

	wpStatusHealthTriageMsg = strings.TrimSpace(`
To see more information about the failing application, please check out the application logs:

waypoint logs -app=APP-NAME

The projects listed above represent their current state known
in Waypoint server. For more information about an application defined in the
project %[1]q can be viewed by running the command:

waypoint status -app=APP-NAME -project=%[1]s
`)

	wpProjectNotFound = strings.TrimSpace(`
No project named %q was found for the server context %q. To see a list of
currently configured projects, run “waypoint project list”.

If you want more information for a specific application, use the '-app' flag
with “waypoint status -app=APP-NAME -project=PROJECT-NAME”.
`)

	wpAppNotFound = strings.TrimSpace(`
No application named %q was found in project %q for the server context %q. To see a
list of currently configured projects, run “waypoint project list”.
`)
)
