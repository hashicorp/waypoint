// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/posener/complete"
)

type ConfigGetCommand struct {
	*baseCommand

	json       bool
	raw        bool
	flagRunner bool
	flagLabels map[string]string
}

func (c *ConfigGetCommand) Run(args []string) int {
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
		if c.flagProject != "" {
			initOpts = append(initOpts,
				WithNoConfig(), // no waypoint.hcl
			)
		}
	}

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(initOpts...); err != nil {
		return 1
	}

	// Get our API client
	client := c.project.Client()

	var prefix string
	switch len(c.args) {
	case 0:
		// ok
	case 1:
		prefix = c.args[0]
	default:
		fmt.Fprintf(os.Stderr, "config-get requires 1 arguments: a variable name prefix")
		return 1
	}

	if c.refProject == nil {
		fmt.Fprintf(os.Stderr, "config-get requires project flag outside the project directory")
		return 1
	}

	req := &pb.ConfigGetRequest{
		Prefix: prefix,
		Scope:  &pb.ConfigGetRequest_Project{Project: c.refProject},
	}

	if c.flagApp != "" {
		req.Scope = &pb.ConfigGetRequest_Application{
			Application: &pb.Ref_Application{
				Project:     c.refProject.Project,
				Application: c.flagApp,
			},
		}

		// The workspace and label flags only apply if the application is set.
		req.Workspace = c.refWorkspace
		req.Labels = c.flagLabels
	}

	if c.flagRunner {
		req.Runner = &pb.Ref_RunnerId{
			// Specifying a non-existent ID will return the runner
			// vars set for all since none will match this ID (since
			// we use ULIDs).
			Id: "-",
		}
	}

	resp, err := client.GetConfig(c.Ctx, req)
	if err != nil {
		c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	if c.json {
		// Get our direct stdout handle cause we're going to be writing colors
		// and want color detection to work.
		out, _, err := c.project.UI.OutputWriters()
		if err != nil {
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		vars := map[string]string{}

		for _, cv := range resp.Variables {
			value := ""
			switch v := cv.Value.(type) {
			case *pb.ConfigVar_Static:
				value = v.Static

			case *pb.ConfigVar_Dynamic:
				value = fmt.Sprintf("<dynamic via %s>", v.Dynamic.From)
			}

			vars[cv.Name] = value
		}

		json.NewEncoder(out).Encode(vars)
		return 0
	}

	if c.raw {
		// Get our direct stdout handle cause we're going to be writing colors
		// and want color detection to work.
		out, _, err := c.project.UI.OutputWriters()
		if err != nil {
			c.project.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		if prefix == "" {
			for _, cv := range resp.Variables {
				fmt.Printf("%s=%s\n", cv.Name, cv.Value)
			}
			return 0
		}

		if len(resp.Variables) == 0 {
			fmt.Fprintf(os.Stderr, "named variable '%s' was not found in config", prefix)
			return 1
		}

		if resp.Variables[0].Name != prefix {
			fmt.Fprintf(os.Stderr, "name '%s' doesn't match prefix: %s", resp.Variables[0].Name, prefix)
			return 1
		}

		fmt.Fprintf(out, "%s=%s\n", resp.Variables[0].Name, resp.Variables[0].Value)
		return 0
	}

	table := terminal.NewTable("Scope", "Name", "Value", "Workspace", "Labels")
	for _, v := range resp.Variables {
		scope := "<unknown>"
		switch appscope := v.Target.AppScope.(type) {
		case *pb.ConfigVar_Target_Global:
			scope = "global"

		case *pb.ConfigVar_Target_Project:
			scope = "project"

		case *pb.ConfigVar_Target_Application:
			scope = "app: " + appscope.Application.Application
		}

		value := ""
		switch v := v.Value.(type) {
		case *pb.ConfigVar_Static:
			value = v.Static

		case *pb.ConfigVar_Dynamic:
			value = fmt.Sprintf("<dynamic via %s>", v.Dynamic.From)
		}

		ws := ""
		if v.Target.Workspace != nil {
			ws = v.Target.Workspace.Workspace
		}

		table.Rich([]string{
			scope,
			v.Name,
			value,
			ws,
			v.Target.LabelSelector,
		}, []string{
			"",
			terminal.Green,
			"",
		})
	}

	c.ui.Table(table)

	return 0
}

func (c *ConfigGetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.json,
			Usage:  "Output in JSON",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "raw",
			Target: &c.raw,
			Usage:  "Output in key=val",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "runner",
			Target: &c.flagRunner,
			Usage: "Show configuration that is set on runners. This will not " +
				"show any configuration that is set on any applications. " +
				"This only includes configuration set with the -runner flag.",
			Default: false,
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "label",
			Target: &c.flagLabels,
			Usage:  "Labels to set for this operation. Can be specified multiple times.",
		})
	})
}

func (c *ConfigGetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ConfigGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigGetCommand) Synopsis() string {
	return "Get config variables."
}

func (c *ConfigGetCommand) Help() string {
	return formatHelp(`
Usage: waypoint config get [prefix]

  Retrieve and print all config variables previously configured that have
  the given prefix. If no prefix is given, all variables are returned.

  The output of this command depends on whether an exact application is
  specified or not. If no app is specified with "-app", this will list
  all POSSIBLE config vars that could be set for applications within the project.
  It is the list of "possible" variables because some are dependent on the
  application name, workspace, or labels.

  By specifying the "-app" flag you can look at config variables that
  would be resolved for a specific application. This will never show duplicate
  variables and will show the full resolved list of variables that will be
  set for a specific application. This will default to a workspace of "default,"
  the current project, and empty labels.

  To further simulate what config variables an application would receive,
  you may specify the "-workspace" flag or the "-label" flag (repeatedly)
  to set the label context.

  The "-runner" flag can be specified to show the list of variables that
  would be set for a runner. All of the above filtering applies to runners,
  as well.

` + c.Flags().Help())
}
