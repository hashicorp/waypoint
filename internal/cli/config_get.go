package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
)

type ConfigGetCommand struct {
	*baseCommand

	json bool
	raw  bool
}

func (c *ConfigGetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
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

	req := &pb.ConfigGetRequest{
		Scope:  &pb.ConfigGetRequest_Project{Project: c.project.Ref()},
		Prefix: prefix,
	}
	if c.flagApp != "" {
		req.Scope = &pb.ConfigGetRequest_Application{
			Application: &pb.Ref_Application{
				Project:     c.project.Ref().Project,
				Application: c.flagApp,
			},
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

	table := terminal.NewTable("Scope", "Name", "Value")
	for _, v := range resp.Variables {
		var app string
		if scope, ok := v.Scope.(*pb.ConfigVar_Application); ok {
			app = scope.Application.Application
		}

		value := ""
		switch v := v.Value.(type) {
		case *pb.ConfigVar_Static:
			value = v.Static

		case *pb.ConfigVar_Dynamic:
			value = fmt.Sprintf("<dynamic via %s>", v.Dynamic.From)
		}

		table.Rich([]string{
			app,
			v.Name,
			value,
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
Usage: waypoint config-get [prefix]

  Retrieve and print all config variables previously configured that have
  the given prefix. If no prefix is given, all variables are returned.

  By specifying the "-app" flag you can look at config variables for
  a specific application rather than the project.

` + c.Flags().Help())
}
