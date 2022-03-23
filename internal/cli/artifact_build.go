package cli

import (
	"context"
	"sort"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ArtifactBuildCommand struct {
	*baseCommand

	flagPush bool
}

func (c *ArtifactBuildCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		app.UI.Output("Building %s...", app.Ref().Application, terminal.WithHeaderStyle())
		buildResult, err := app.Build(ctx, &pb.Job_BuildOp{
			DisablePush: !c.flagPush,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if buildResult.Push != nil {
			app.UI.Output("\nCreated artifact v%d", buildResult.Push.Sequence)
		}

		// Show input variable values used in build
		app.UI.Output("Variables used:", terminal.WithHeaderStyle())
		headers := []string{
			"Variable", "Value", "Type", "Source",
		}

		tbl := terminal.NewTable(headers...)
		// sort alphabetically for joy
		inputVars := make([]string, 0, len(buildResult.Build.VariableRefs))
		for iv := range buildResult.Build.VariableRefs {
			inputVars = append(inputVars, iv)
		}
		sort.Strings(inputVars)
		for _, iv := range inputVars {
			// We add a line break in the value here because the Table word wrap
			// alone can't accomodate the column headers to a long value
			val := buildResult.Build.VariableRefs[iv].Value
			if len(val) > 45 {
				for i := range val {
					// line break every 45 characters
					if i%46 == 0 && i != 0 {
						val = val[:i] + "\n" + val[i:]
					}
				}
			}
			columns := []string{
				iv,
				val,
				buildResult.Build.VariableRefs[iv].Type,
				buildResult.Build.VariableRefs[iv].Source,
			}
			tbl.Rich(
				columns,
				[]string{
					terminal.Green,
				},
			)
		}
		c.ui.Table(tbl)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ArtifactBuildCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "push",
			Target:  &c.flagPush,
			Default: true,
			Usage:   "Push the artifact to the configured registry.",
		})
	})
}

func (c *ArtifactBuildCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ArtifactBuildCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ArtifactBuildCommand) Synopsis() string {
	return "Build a new versioned artifact from source"
}

func (c *ArtifactBuildCommand) Help() string {
	return formatHelp(`
Usage: waypoint artifact build [options]
Alias: waypoint build [options]

  Build a new versioned artifact from source.

` + c.Flags().Help())
}
