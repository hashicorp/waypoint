package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type FmtCommand struct {
	*baseCommand

	flagWrite bool
}

func (c *FmtCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithClient(false),
	); err != nil {
		return 1
	}

	// If we have too many args, error immediately.
	if len(c.args) > 1 {
		c.ui.Output("At most one argument is expected.\n\n"+c.Help(), terminal.WithErrorStyle())
		return 1
	}

	// If we have no args, default to the filename
	if len(c.args) == 0 {
		c.args = []string{configpkg.Filename}
	}

	// Read the input
	src, err := c.readInput()
	if err != nil {
		c.ui.Output(
			"Error reading input to format: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// Format it
	name := "<stdin>"
	stdin := true
	if c.args[0] != "-" {
		name = filepath.Base(c.args[0])
		stdin = false
	}
	out, err := configpkg.Format(src, name)
	if err != nil {
		c.ui.Output(
			"Error formatting: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// If we're writing then write it to the file. stdin never writes to a file
	if c.flagWrite && !stdin {
		if err := ioutil.WriteFile(c.args[0], out, 0644); err != nil {
			c.ui.Output(
				"Error writing formatted output: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}
	} else {
		// We must use fmt here and not c.ui since c.ui may wordwrap and trim.
		fmt.Print(string(out))
	}

	return 0
}

func (c *FmtCommand) readInput() ([]byte, error) {
	// If we have non-stdin input then read it
	if c.args[0] != "-" {
		return ioutil.ReadFile(c.args[0])
	}

	// Otherwise it is stdin
	return ioutil.ReadAll(os.Stdin)
}

func (c *FmtCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:    "write",
			Target:  &c.flagWrite,
			Default: true,
			Usage: "Overwrite the input file. If this is false, the formatted " +
				"output will be written to STDOUT. This has no effect when formatting " +
				"from STDIN.",
		})
	})
}

func (c *FmtCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *FmtCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *FmtCommand) Synopsis() string {
	return "Rewrite waypoint.hcl configuration to a canonical format"
}

func (c *FmtCommand) Help() string {
	return formatHelp(`
Usage: waypoint fmt [FILE]

  Rewrite a waypoint.hcl file to a canonical format.

  This only works for HCL-formatted Waypoint configuration files. JSON-formatted
  files do not work and will result in an error.

  If FILE is not specified, then the current directory will be searched
  for a "waypoint.hcl" file. If FILE is "-" then the content will be read
  from stdin.

  This command does not validate the waypoint.hcl configuration. This will
  work for older and newer configuration formats.

` + c.Flags().Help())
}
