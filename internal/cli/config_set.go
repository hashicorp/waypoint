package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
	"github.com/posener/complete"
)

type ConfigSetCommand struct {
	*baseCommand

	app string
}

func (c *ConfigSetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		fmt.Fprintf(os.Stderr, "config-set requires at least one key=value entry")
		return 1
	}

	// Get our API client
	client := c.project.Client()

	var req pb.ConfigSetRequest

	for _, arg := range c.args {
		idx := strings.IndexByte(arg, '=')
		if idx == -1 || idx == 0 {
			fmt.Fprintf(os.Stderr, "variables must be in the form key=value")
			return 1
		}

		configVar := &pb.ConfigVar{
			Name:  arg[:idx],
			Value: arg[idx+1:],
		}

		if c.app == "" {
			configVar.Scope = &pb.ConfigVar_Project{
				Project: c.project.Ref(),
			}
		} else {
			configVar.Scope = &pb.ConfigVar_Application{
				Application: &pb.Ref_Application{
					Project:     c.project.Ref().Project,
					Application: c.app,
				},
			}
		}

		req.Variables = append(req.Variables, configVar)
	}

	_, err := client.SetConfig(c.Ctx, &req)
	if err != nil {
		c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *ConfigSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:   "app",
			Target: &c.app,
			Usage:  "Scope the variables to a specific app.",
		})
	})
}

func (c *ConfigSetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigSetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigSetCommand) Synopsis() string {
	return "Set a config variable."
}

func (c *ConfigSetCommand) Help() string {
	helpText := `
Usage: waypoint config-set <name> <value>

  Set a config variable that will be available to deployments as an environment variable.

`

	return strings.TrimSpace(helpText)
}
