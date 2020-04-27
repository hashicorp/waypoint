package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/core"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	servercomponent "github.com/hashicorp/waypoint/internal/server/component"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type ReleaseCreateCommand struct {
	*baseCommand
}

func (c *ReleaseCreateCommand) Run(args []string) int {
	defer c.Close()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
	}

	// Get the args
	args = c.Flags().Args()
	if len(args) == 0 {
		args = []string{"100"}
	}
	if len(args) > 1 {
		c.project.UI.Output(c.Flags().Help(), terminal.WithErrorStyle())
		return 1
	}

	// Parse the single argument. We allow a trailing % since its expected
	// users might try this and we don't want to error for what is a reasonable
	// (though alternate) input style.
	value := args[0]
	if len(value) > 0 && value[len(value)-1] == '%' {
		value = value[:len(value)-1]
	}

	// Parse the percentage number
	number, err := strconv.Atoi(value)
	if err != nil {
		c.project.UI.Output("Failed to parse the percentage value %s: %s", value, err,
			terminal.WithErrorStyle())
		return 1
	}
	if number < 0 || number > 100 {
		c.project.UI.Output("Percentage value %q must be greater than 0 and less than or equal to 100", value,
			terminal.WithErrorStyle())
		return 1
	}

	client := c.project.Client()
	err = c.DoApp(c.Ctx, func(ctx context.Context, app *core.App) error {
		// Get the latest deployment
		resp, err := client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Limit:     2,
			Order:     pb.ListDeploymentsRequest_COMPLETE_TIME,
			OrderDesc: true,
		})
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}
		if len(resp.Deployments) == 0 {
			app.UI.Output("No successful deployments found.", terminal.WithErrorStyle())
			return ErrSentinel
		}

		type target struct {
			Deployment *pb.Deployment
			Target     component.ReleaseTarget
		}

		// Build our targets
		var targets []target
		targets = append(targets, target{
			resp.Deployments[0],
			component.ReleaseTarget{
				DeploymentId: resp.Deployments[0].Id,
				Deployment:   servercomponent.Deployment(resp.Deployments[0]),
				Percent:      uint(number),
			},
		})
		if number < 100 {
			if len(resp.Deployments) < 2 {
				app.UI.Output("Traffic splitting requires multiple successful deploys.", terminal.WithErrorStyle())
				return ErrSentinel
			}

			targets = append(targets, target{
				resp.Deployments[1],
				component.ReleaseTarget{
					Deployment: servercomponent.Deployment(resp.Deployments[1]),
					Percent:    uint(100 - number),
				},
			})
		}

		// UI
		app.UI.Output("Releasing...", terminal.WithHeaderStyle())
		for _, t := range targets {
			completeTime, _ := ptypes.Timestamp(t.Deployment.Status.CompleteTime)

			app.UI.Output("%d%%: ID %s (%s)",
				t.Target.Percent,
				t.Deployment.Id,
				humanize.Time(completeTime),
				terminal.WithStatusStyle(),
			)
		}

		// Release
		targetArgs := make([]component.ReleaseTarget, len(targets))
		for i, target := range targets {
			targetArgs[i] = target.Target
		}
		release, err := app.Release(ctx, targetArgs)
		if err != nil {
			app.UI.Output(err.Error(), terminal.WithErrorStyle())
			return ErrSentinel
		}

		app.UI.Output("\nURL: %s", release.URL(), terminal.WithSuccessStyle())
		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *ReleaseCreateCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *ReleaseCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ReleaseCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ReleaseCreateCommand) Synopsis() string {
	return "Release a deployment."
}

func (c *ReleaseCreateCommand) Help() string {
	helpText := `
Usage: devflow release [options] [percentage...]

  Open a deployment to traffic. This will default to shifting traffic
  100% to the latest deployment. You can specify multiple percentages to
  split traffic between releases.

Examples:

  "devflow release" - will send 100% of traffic to the latest deployment.

  "devflow release 90" - will send 90% of traffic to the latest deployment
  and 10% of traffic to the prior deployment.

  "devflow release '+10'" - will send 10% more traffic to the latest deployment.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
