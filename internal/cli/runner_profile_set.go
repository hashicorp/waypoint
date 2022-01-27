package cli

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RunnerProfileSetCommand struct {
	*baseCommand

	flagName           string
	flagOCIUrl         string
	flagEnvVars        []string
	flagPluginType     string
	flagPluginConfig   string
	flagDefault        bool
	flagTargetRunnerId string
}

func (c *RunnerProfileSetCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flagSet.Args()
	ctx := c.Ctx

	sg := c.ui.StepGroup()
	defer sg.Wait()

	var s terminal.Step

	defer func() {
		if s != nil {
			s.Abort()
		}
	}()

	var (
		od      *pb.OnDemandRunnerConfig
		updated bool
	)

	if c.flagName != "" {
		s = sg.Add("Checking for an existing runner profile: %s", c.flagName)
		// Check for an existing project of the same name.
		resp, err := c.project.Client().GetOnDemandRunnerConfig(ctx, &pb.GetOnDemandRunnerConfigRequest{
			Config: &pb.Ref_OnDemandRunnerConfig{
				Name: c.flagName,
			},
		})
		if status.Code(err) == codes.NotFound {
			// If the error is a not found error, act as though there is no error
			// and the project is nil so that we can handle that later.
			resp = nil
			err = nil
		}
		if err != nil {
			c.ui.Output(
				"Error checking for project: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}

		if resp != nil {
			od = resp.Config
			s.Update("Updating runner profile %q (%q)...", od.Name, od.Id)
			updated = true
		} else {
			s.Update("No existing runner profile found for id %q...command will create a new profile", c.flagName)
			od = &pb.OnDemandRunnerConfig{
				Name: c.flagName,
			}
		}
	} else {
		s = sg.Add("Creating new runner profile")
		od = &pb.OnDemandRunnerConfig{
			Name: c.flagName,
		}
	}

	// Set target runner for profile
	if c.flagTargetRunnerId != "" {
		od.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: c.flagTargetRunnerId,
				},
			},
		}
	} else {
		od.TargetRunner = &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{}}
	}

	// If we were specified a file then we're going to load that up.
	if c.flagPluginConfig != "" {
		path, err := filepath.Abs(c.flagPluginConfig)
		if err != nil {
			c.ui.Output(
				"Error loading HCL file specified with the -plugin-config flag: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			c.ui.Output(
				"Error reading HCL plugin config: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		od.PluginConfig = data

		switch filepath.Ext(path) {
		case ".hcl":
			od.ConfigFormat = pb.Hcl_HCL
			_, diag := hclsyntax.ParseConfig(data, "<waypoint-hcl>", hcl.Pos{})
			if diag.HasErrors() {
				c.ui.Output(
					"Syntax errors in file specified with -plugin-config: %s",
					clierrors.Humanize(diag),
					terminal.WithErrorStyle(),
				)

				return 1
			}

		case ".json":
			od.ConfigFormat = pb.Hcl_JSON
			_, diag := hcljson.Parse(data, "<waypoint-hcl>")
			if diag.HasErrors() {
				c.ui.Output(
					"Syntax errors in file specified with -plugin-config: %s",
					clierrors.Humanize(diag),
					terminal.WithErrorStyle(),
				)

				return 1
			}

		default:
			c.ui.Output(
				"File specified via -plugin-config must end in '.hcl' or '.json'",
				terminal.WithErrorStyle(),
			)

			return 1
		}
	}

	od.OciUrl = c.flagOCIUrl
	od.EnvironmentVariables = map[string]string{}
	od.Default = c.flagDefault

	for _, kv := range c.flagEnvVars {
		idx := strings.IndexByte(kv, '=')
		if idx != -1 {
			od.EnvironmentVariables[kv[:idx]] = kv[idx+1:]
		}
	}

	if c.flagPluginType == "" {
		c.ui.Output(
			"Flag '-plugin-type' must be set to a valid plugin type like 'docker' or 'kubernetes'.\n\n%s",
			c.Help(),
			terminal.WithErrorStyle(),
		)

		return 1
	}

	od.PluginType = c.flagPluginType

	// Upsert
	_, err := c.project.Client().UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: od,
	})
	if err != nil {
		c.ui.Output(
			"Error upserting runner profile: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	if updated {
		s.Update("Runner profile updated")
	} else {
		s.Update("Runner profile created")
	}
	s.Done()

	return 0
}

func (c *RunnerProfileSetCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "name",
			Target:  &c.flagName,
			Default: "",
			Usage:   "The name of an existing runner profile to update.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "oci-url",
			Target:  &c.flagOCIUrl,
			Default: "hashicorp/waypoint-odr:latest",
			Usage:   "The url for the OCI image to launch for the on-demand runner.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "env-vars",
			Target: &c.flagEnvVars,
			Usage: "Environment variable to expose to the on-demand runner. Typically used to " +
				"introduce configuration for the plugins that the runner will execute. Can be specified multiple times.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "plugin-type",
			Target:  &c.flagPluginType,
			Default: "",
			Usage:   "The type of the plugin to launch for the on-demand runner, such as aws-ecs, kubernetes, etc.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "plugin-config",
			Target:  &c.flagPluginConfig,
			Default: "",
			Usage: "Path to an hcl file that contains the configuration for the plugin. " +
				"This is only necessary when the plugin's defaults need to be adjusted for " +
				"the environment the plugin will launch the on-demand runner in.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "default",
			Target:  &c.flagDefault,
			Default: false,
			Usage: "Indicates that this remote runner profile should be the default for any project that doesn't " +
				"otherwise specify its own remote runner.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "target-runner-id",
			Target:  &c.flagTargetRunnerId,
			Default: "",
			Usage:   "ID of the remote runner to target for the profile.",
		})
	})
}

func (c *RunnerProfileSetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerProfileSetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerProfileSetCommand) Synopsis() string {
	return "Create or update a runner profile."
}

func (c *RunnerProfileSetCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner profile set [OPTIONS]

  Create or update a runner profile.

  This will register a new runner profile with the given options. If
  a runner profile with the same id already exists, this will update the
  existing runner profile using the fields that are set.

  Waypoint will use a runner profile to spawn containers for
  various kinds of work as needed on the platform requested during any given
  lifecycle operation.

` + c.Flags().Help())
}
