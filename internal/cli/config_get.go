package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
	"github.com/olekukonko/tablewriter"
	"github.com/posener/complete"
)

type ConfigGetCommand struct {
	*baseCommand

	json bool
	raw  bool
	app  string
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

	resp, err := client.GetConfig(c.Ctx, &pb.ConfigGetRequest{
		Prefix: prefix,
		App:    c.app,
	})

	if err != nil {
		c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	// Get our direct stdout handle cause we're going to be writing colors
	// and want color detection to work.
	out, _, err := c.project.UI.OutputWriters()
	if err != nil {
		c.project.UI.Output(err.Error(), terminal.WithErrorStyle())
		return 1
	}

	if c.json {
		vars := map[string]string{}

		for _, cv := range resp.Variables {
			vars[cv.Name] = cv.Value
		}

		json.NewEncoder(out).Encode(vars)
		return 0
	}

	if c.raw {
		if len(resp.Variables) == 0 {
			return 1
		}

		if resp.Variables[0].Name != prefix {
			return 1
		}

		fmt.Fprintln(out, resp.Variables[0].Value)
		return 0
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"Scope", "Name", "Value"})
	table.SetBorder(false)

	for _, v := range resp.Variables {
		table.Rich([]string{
			v.App,
			v.Name,
			v.Value,
		}, []tablewriter.Colors{
			{},
			{tablewriter.FgGreenColor},
			{},
		})
	}

	table.Render()

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
			Usage:  "Output the value for the named variable only (disables prefix matching)",
		})

		f.StringVar(&flag.StringVar{
			Name:   "app",
			Target: &c.app,
			Usage:  "Scope the variables to a specific app.",
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
	helpText := `
Usage: waypoint config-get [prefix]

  Retrieve and print all config variables previously configured that have the given prefix.
	If no prefix is given, all variables are returned.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
