package cli

import (
	"fmt"
	"strings"

	"github.com/posener/complete"

	"github.com/mitchellh/devflow/internal/pkg/flag"
	"github.com/mitchellh/devflow/internal/serverinstall"
	"github.com/mitchellh/devflow/sdk/terminal"
)

type InstallCommand struct {
	*baseCommand

	config serverinstall.Config
}

func (c *InstallCommand) Run(args []string) int {
	defer c.Close()

	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	// Decode our configuration
	output, err := serverinstall.Render(&c.config)
	if err != nil {
		c.ui.Output(
			"Error generating configuration: %s", err.Error(),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	out, _, err := c.ui.OutputWriters()
	if err != nil {
		panic(err)
	}

	fmt.Fprint(out, output)
	if output[:len(output)-1] != "\n" {
		fmt.Fprint(out, "\n")
	}

	return 0
}

func (c *InstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.StringVar(&flag.StringVar{
			Name:    "namespace",
			Target:  &c.config.Namespace,
			Usage:   "Kubernetes namespace install into.",
			Default: "default",
		})

		f.StringVar(&flag.StringVar{
			Name:    "service",
			Target:  &c.config.ServiceName,
			Usage:   "Name of the Kubernetes service for the server.",
			Default: "devflow",
		})

		f.StringVar(&flag.StringVar{
			Name:    "server-image",
			Target:  &c.config.ServerImage,
			Usage:   "Docker image for the server image.",
			Default: "gcr.io/mitchellh-test/devflow:latest",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "annotate-service",
			Target: &c.config.ServiceAnnotations,
			Usage:  "Annotations for the Service generated.",
		})
	})
}

func (c *InstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *InstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *InstallCommand) Synopsis() string {
	return "Output Kubernetes configurations to run a self-hosted server."
}

func (c *InstallCommand) Help() string {
	helpText := `
Usage: devflow install [options]

  Outputs the Kubernetes configurations required to run a self-hosted
  Devflow server. You can deploy to Kubernetes by piping this to kubectl.

  Example: devflow install | kubectl apply -f -

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
