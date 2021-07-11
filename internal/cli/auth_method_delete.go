package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
)

type AuthMethodDeleteCommand struct {
	*baseCommand
}

func (c *AuthMethodDeleteCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	if len(c.args) == 0 {
		c.ui.Output("auth method name required for deletion", terminal.WithErrorStyle())
		return 1
	}
	name := c.args[0]

	_, err := c.project.Client().DeleteAuthMethod(c.Ctx, &pb.DeleteAuthMethodRequest{
		AuthMethod: &pb.Ref_AuthMethod{Name: name},
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	c.ui.Output("Auth method deleted.", terminal.WithSuccessStyle())
	return 0
}

func (c *AuthMethodDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, nil)
}

func (c *AuthMethodDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AuthMethodDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AuthMethodDeleteCommand) Synopsis() string {
	return "Delete a previously configured auth method."
}

func (c *AuthMethodDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint auth-method delete NAME

  Delete a previously configured auth method.

  This will not delete any users, although users may no longer be able to
  log in. Already authenticated users will remain logged in even if they
  authenticated using this auth method.

`)
}
