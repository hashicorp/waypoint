// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type ConfigDeleteCommand struct {
	*baseCommand

	flagGlobal     bool
	flagScope      string
	flagLabelScope string
}

func (c *ConfigDeleteCommand) Run(args []string) int {
	initOpts := []Option{
		WithArgs(args),
		WithFlags(c.Flags()),

		// Don't allow a local in-mem server because configuration
		// makes no sense with the local server.
		WithNoLocalServer(),
	}

	// We parse our flags twice in this command because we need to
	// determine if we're loading a config or not.
	//
	// NOTE we specifically ignore errors here because if we have errors
	// they'll happen again on Init and Init will output to the CLI.
	if err := c.Flags().Parse(args); err == nil {
		// If we're global scoped OR we have a project explicitly set
		// then we do not need a config. If we're not global scoped and
		// we do not have a project explicitly set, we need a config because
		// we need a way to load that project name.
		if c.flagScope == "global" || c.flagProject != "" {
			initOpts = append(initOpts,
				WithNoConfig(), // no waypoint.hcl
			)
		}
	}

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(initOpts...); err != nil {
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
			_, _ = fmt.Fprintf(os.Stderr, "config delete requires at least one key entry")
			return 1
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			c.args = append(c.args, scanner.Text())
		}
	}

	// Pre-calculate our project ref since we reuse this.
	projectRef := &pb.Ref_Project{Project: c.flagProject}
	if c.flagProject == "" && c.project != nil {
		projectRef = c.project.Ref()
	}

	// Get our API client
	client := c.project.Client()

	var req pb.ConfigDeleteRequest

	for _, arg := range c.args {
		// Build our initial config var
		configVar := &pb.ConfigVar{
			Target: &pb.ConfigVar_Target{},
			Name:   arg,
			Value:  &pb.ConfigVar_Unset{},
		}

		// Depending on the scoping set our target
		switch c.flagScope {
		case "global":
			configVar.Target.AppScope = &pb.ConfigVar_Target_Global{
				Global: &pb.Ref_Global{},
			}

		case "project":
			configVar.Target.AppScope = &pb.ConfigVar_Target_Project{
				Project: projectRef,
			}

		case "app":
			if c.flagApp == "" {
				fmt.Fprintf(os.Stderr, "-scope requires -app set if scope is 'app'")
				return 1
			}
			configVar.Target.AppScope = &pb.ConfigVar_Target_Application{
				Application: &pb.Ref_Application{
					Project:     projectRef.Project,
					Application: c.flagApp,
				},
			}

		default:
			err := fmt.Errorf("-scope needs to be one of 'global', 'project', or 'app'")
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		// If we have a workspace flag set, set that.
		if v := c.flagWorkspace; v != "" {
			configVar.Target.Workspace = &pb.Ref_Workspace{
				Workspace: v,
			}
		}

		// If we have a label flag set, set that.
		if v := c.flagLabelScope; v != "" {
			configVar.Target.LabelSelector = v
		}

		req.Variables = append(req.Variables, configVar)
	}

	_, err := client.DeleteConfig(c.Ctx, &req)
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	return 0
}

func (c *ConfigDeleteCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:   "scope",
			Target: &c.flagScope,
			Usage: "The scope for this configuration to delete. The configuration will only " +
				"delete within this scope. This can be one of 'global', 'project', or " +
				"'app'.",
			Default: "project",
		})

		f.StringVar(&flag.StringVar{
			Name:   "label-scope",
			Target: &c.flagLabelScope,
			Usage: "If set, configuration will only be deleted if the deployment " +
				"or operation has a matching label set.",
			Default: "",
		})
	})
}

func (c *ConfigDeleteCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigDeleteCommand) Synopsis() string {
	return "Delete a config variable."
}

func (c *ConfigDeleteCommand) Help() string {
	return formatHelp(`
Usage: waypoint config delete <name> <name2> ...

  Delete a config variable from the system. This cannot be undone.

` + c.Flags().Help())
}
