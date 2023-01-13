package cli

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/config/variables/formatter"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
		WithMultiAppTargets(),
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

		// Show input variable values used in build
		// We do this here so that if the list is long, it doesn't
		// push the deploy/release URLs off the top of the terminal.
		// BuildResult, DeployResult, and ReleaseResult all store
		// used VariableRefs. We use Release just because it's last.
		app.UI.Output("Variables used:", terminal.WithHeaderStyle())
		resp, err := c.project.Client().GetJob(ctx, &pb.GetJobRequest{
			JobId: result.Release.Release.JobId,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		tbl := fmtVariablesOutput(resp.VariableFinalValues)
		c.ui.Table(tbl)

		// Common reused values
		releaseUrl := result.Up.ReleaseUrl
		appUrl := result.Up.AppUrl
		deployUrl := result.Up.DeployUrl

		// inplace is true if this was an in-place deploy. We detect this
		// if we have a generation that uses a non-matching sequence number
		inplace := result.Deploy.Deployment.Generation != nil &&
			result.Deploy.Deployment.Generation.Id != "" &&
			result.Deploy.Deployment.Generation.InitialSequence != result.Deploy.Deployment.Sequence

		// Ensure deploy and release Urls have a scheme
		deployUrl, err = addUrlScheme(deployUrl, httpsScheme)
		if err != nil {
			return err
		}
		releaseUrl, err = addUrlScheme(releaseUrl, httpsScheme)
		if err != nil {
			return err
		}

		// Output
		app.UI.Output("")
		switch {
		case releaseUrl != "":
			if !inplace {
				app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			} else {
				app.UI.Output(strings.TrimSpace(deployInPlace)+"\n", terminal.WithSuccessStyle())
			}
			app.UI.Output("   Release URL: %s", releaseUrl, terminal.WithSuccessStyle())
			if deployUrl != "" {
				app.UI.Output("Deployment URL: %s", deployUrl, terminal.WithSuccessStyle())
			}

		case appUrl != "" && deployUrl != "":
			if !inplace {
				app.UI.Output(strings.TrimSpace(deployURLService)+"\n", terminal.WithSuccessStyle())
			} else {
				app.UI.Output(strings.TrimSpace(deployInPlace)+"\n", terminal.WithSuccessStyle())
			}
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
	return "Perform the build, deploy, and release steps"
}

func (c *UpCommand) Help() string {
	return formatHelp(`
Usage: waypoint up [options]

  Perform the build, deploy, and release steps.

` + c.Flags().Help())
}

// Helper functions for formatting variable final value output
func fmtVariablesOutput(values map[string]*pb.Variable_FinalValue) *terminal.Table {
	headers := []string{
		"Variable", "Value", "Type", "Source",
	}
	tbl := terminal.NewTable(headers...)
	output := formatter.ValuesForOutput(values)
	var columns []string
	for iv, v := range output {
		// We add a line break in the value here because the Table word wrap
		// alone can't accomodate the column headers to a long value
		if len(v.Value) > 45 {
			for i := 45; i < len(v.Value); i += 46 {
				v.Value = v.Value[:i] + "\n" + v.Value[i:]
			}
		}
		columns = []string{
			iv,
			v.Value,
			v.Type,
			v.Source,
		}
		tbl.Rich(
			columns,
			[]string{
				terminal.Green,
			},
		)
	}
	return tbl
}
