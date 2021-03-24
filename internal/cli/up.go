package cli

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type UpCommand struct {
	*baseCommand

	flagPrune       bool
	flagPruneRetain int
}

func (c *UpCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		result, err := app.Up(ctx, &pb.Job_UpOp{
			Release: &pb.Job_ReleaseOp{
				Prune:               c.flagPrune,
				PruneRetain:         int32(c.flagPruneRetain),
				PruneRetainOverride: c.flagPruneRetain >= 0,
			},
		})
		if c.legacyRequired(err) {
			// An older Waypoint server version that doesn't support the
			// "up" operation, so fall back.
			c.Log.Warn("server doesn't support 'up' operation, falling back")
			return c.legacyUp(ctx, app)
		}
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		// Common reused values
		releaseUrl := result.Up.ReleaseUrl
		appUrl := result.Up.AppUrl
		deployUrl := result.Up.DeployUrl

		// Output
		app.UI.Output("")
		switch {
		case releaseUrl != "":
			app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			app.UI.Output("   Release URL: %s", releaseUrl, terminal.WithSuccessStyle())
			if deployUrl != "" {
				app.UI.Output("Deployment URL: %s", deployUrl, terminal.WithSuccessStyle())
			}

		case appUrl != "" && deployUrl != "":
			app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			app.UI.Output("           URL: %s", appUrl, terminal.WithSuccessStyle())
			app.UI.Output("Deployment URL: %s", deployUrl, terminal.WithSuccessStyle())

		default:
			app.UI.Output(strings.TrimSpace(deployNoURL)+"\n", terminal.WithSuccessStyle())
		}

		return nil
	})

	if err != nil {
		if err != ErrSentinel {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		}

		return 1
	}

	return 0
}

// legacyRequired returns true if we need to execute the legacy "up" logic.
func (c *UpCommand) legacyRequired(err error) bool {
	switch status.Code(err) {
	case codes.Unimplemented:
		// Waypoint 0.3+ returns this.
		return true

	case codes.FailedPrecondition:
		// Waypoint 0.2 the only way we can detect is this error message
		// since operation will be seen as "nil" to an older server since
		// it can't understand our newer op type.
		return err != nil && strings.Contains(strings.ToLower(err.Error()), "operation")

	default:
		return false
	}
}

// legacyUp implements the exact "up" logic from WP 0.2.x. In WP 0.3.0 we
// introduced a remote "up" job type that we use instead. If the user uses
// a new client but an old server that doesn't support the up job type, this
// will be executed instead.
func (c *UpCommand) legacyUp(
	ctx context.Context,
	app *clientpkg.App,
) error {
	client := c.project.Client()

	// Build it
	app.UI.Output("Building...", terminal.WithHeaderStyle())

	_, err := app.Build(ctx, &pb.Job_BuildOp{})
	if err != nil {
		app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return ErrSentinel
	}

	// Get the most recent pushed artifact
	push, err := client.GetLatestPushedArtifact(ctx, &pb.GetLatestPushedArtifactRequest{
		Application: app.Ref(),
		Workspace:   c.project.WorkspaceRef(),
	})
	if err != nil {
		app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return ErrSentinel
	}

	// Push it
	app.UI.Output("Deploying...", terminal.WithHeaderStyle())

	result, err := app.Deploy(ctx, &pb.Job_DeployOp{
		Artifact: push,
	})
	if err != nil {
		app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return ErrSentinel
	}
	deployUrl := result.Deployment.Preload.DeployUrl

	// Try to get the hostname
	var hostname *pb.Hostname
	hostnamesResp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{
		Target: &pb.Hostname_Target{
			Target: &pb.Hostname_Target_Application{
				Application: &pb.Hostname_TargetApp{
					Application: result.Deployment.Application,
					Workspace:   result.Deployment.Workspace,
				},
			},
		},
	})
	if err == nil && len(hostnamesResp.Hostnames) > 0 {
		hostname = hostnamesResp.Hostnames[0]
	}

	// We're releasing, do that too.
	app.UI.Output("Releasing...", terminal.WithHeaderStyle())
	releaseResult, err := app.Release(ctx, &pb.Job_ReleaseOp{
		Deployment: result.Deployment,
		Prune:      true,
	})
	if err != nil {
		app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return ErrSentinel
	}

	releaseUrl := releaseResult.Release.Url

	// Output
	app.UI.Output("")
	switch {
	case releaseUrl != "":
		app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
		app.UI.Output("   Release URL: %s", releaseUrl, terminal.WithSuccessStyle())
		if deployUrl != "" {
			app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())
		}

	case hostname != nil && deployUrl != "":
		app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
		app.UI.Output("           URL: https://%s", hostname.Fqdn, terminal.WithSuccessStyle())
		app.UI.Output("Deployment URL: https://%s", deployUrl, terminal.WithSuccessStyle())

	default:
		app.UI.Output(strings.TrimSpace(deployNoURL)+"\n", terminal.WithSuccessStyle())
	}

	return nil
}

func (c *UpCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:    "prune",
			Target:  &c.flagPrune,
			Usage:   "Prune old unreleased deployments.",
			Default: true,
		})

		f.IntVar(&flag.IntVar{
			Name:   "prune-retain",
			Target: &c.flagPruneRetain,
			Usage: "The number of unreleased deployments to keep. " +
				"If this isn't set or is set to any negative number, " +
				"then this will default to 1 on the server. If you want to prune " +
				"all unreleased deployments, set this to 0.",
			Default: -1,
		})
	})
}

func (c *UpCommand) Synopsis() string {
	return "Perform the build, deploy, and release steps for the app"
}

func (c *UpCommand) Help() string {
	return formatHelp(`
Usage: waypoint up [options]

  Perform the build, deploy, and release steps for the app.

` + c.Flags().Help())
}
