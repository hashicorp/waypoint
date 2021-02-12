package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
)

type ConfigSetCommand struct {
	*baseCommand
}

func (c *ConfigSetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	// If there are no command arguments, check if the command has
	// been invoked with a pipe like `cat .env | waypoint config set`.
	if len(c.args) == 0 {
		info, err := os.Stdin.Stat()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to get console mode for stdin")
			return 1
		}

		// If there's no pipe, there are no arguments. Fail.
		if info.Mode()&os.ModeNamedPipe == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "config set requires at least one key=value entry")
			return 1
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			c.args = append(c.args, scanner.Text())
		}
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
			Name: arg[:idx],
			Value: &pb.ConfigVar_Static{
				Static: arg[idx+1:],
			},
		}

		if c.flagApp == "" {
			configVar.Scope = &pb.ConfigVar_Project{
				Project: c.project.Ref(),
			}
		} else {
			configVar.Scope = &pb.ConfigVar_Application{
				Application: &pb.Ref_Application{
					Project:     c.project.Ref().Project,
					Application: c.flagApp,
				},
			}
		}

		req.Variables = append(req.Variables, configVar)
	}

	_, err := client.SetConfig(c.Ctx, &req)
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *ConfigSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
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
	return formatHelp(`
Usage: waypoint config set <name>=<value>

  Set a config variable that will be available to deployments as an
  environment variable.

  This will scope the variable to the entire project by default.
  Specify the "-app" flag to set a config variable for a specific app.

`)
}
