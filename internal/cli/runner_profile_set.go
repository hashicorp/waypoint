package cli

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/installutil"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerProfileSetCommand struct {
	*baseCommand
	//TODO(XX): after `-env-vars` as a slice is deprecated, rename flagEnvVar to flagEnvVars
	flagName               string
	flagOCIUrl             *string
	flagEnvVar             map[string]string
	flagEnvVars            []string
	flagPluginType         *string
	flagPluginConfig       string
	flagDefault            *bool
	flagTargetRunnerAny    *bool
	flagTargetRunnerId     *string
	flagTargetRunnerLabels map[string]string
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
	if c.flagTargetRunnerId != nil {
		od.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: *c.flagTargetRunnerId,
				},
			},
		}
		if c.flagTargetRunnerAny != nil {
			c.ui.Output("Both -target-runner-id and -target-runner-any detected, only one can be set at a time. ID takes priority.",
				terminal.WithWarningStyle())
		}
		if c.flagTargetRunnerLabels != nil {
			c.ui.Output("Both -target-runner-id and -target-runner-label detected, only one can be set at a time. ID takes priority.",
				terminal.WithWarningStyle())
		}
	} else if c.flagTargetRunnerLabels != nil {
		od.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Labels{
				Labels: &pb.Ref_RunnerLabels{
					Labels: c.flagTargetRunnerLabels,
				},
			},
		}
		if c.flagTargetRunnerAny != nil {
			c.ui.Output("Both -target-runner-label and -target-runner-any detected, only one can be set at a time. Labels take priority.",
				terminal.WithWarningStyle())
		}
	} else if od.TargetRunner == nil || (c.flagTargetRunnerAny != nil && *c.flagTargetRunnerAny) {
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

	if od.OciUrl == "" {
		od.OciUrl = installutil.DefaultODRImage
	}
	if c.flagOCIUrl != nil {
		od.OciUrl = *c.flagOCIUrl
	}

	if c.flagDefault != nil {
		od.Default = *c.flagDefault
	}

	if c.flagEnvVars != nil {
		//TODO(XX): Deprecate -env-vars and this logic
		c.ui.Output(
			"Flag '-env-vars' is deprecated, please use flag '-env-var=k=v'",
			terminal.WithWarningStyle(),
		)
		for _, kv := range c.flagEnvVars {
			idx := strings.IndexByte(kv, '=')
			if idx != -1 {
				od.EnvironmentVariables[kv[:idx]] = kv[idx+1:]
			}
		}
	}

	if c.flagEnvVar != nil {
		for k, v := range c.flagEnvVar {
			if v == "" {
				delete(od.EnvironmentVariables, k)
			} else {
				od.EnvironmentVariables[k] = v
			}
		}
	}

	if c.flagPluginType != nil || od.PluginType == "" {
		if *c.flagPluginType == "" {
			c.ui.Output(
				"Flag '-plugin-type' must be set to a valid plugin type like 'docker' or 'kubernetes'.\n\n%s",
				c.Help(),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		od.PluginType = *c.flagPluginType
	}

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

		f.StringPtrVar(&flag.StringPtrVar{
			Name:    "oci-url",
			Target:  &c.flagOCIUrl,
			Default: installutil.DefaultODRImage,
			Usage:   "The url for the OCI image to launch for the on-demand runner.",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "env-var",
			Target: &c.flagEnvVar,
			Usage: "Environment variable to expose to the on-demand runner set in 'k=v' format. Typically used to " +
				"introduce configuration for the plugins that the runner will execute. Can be specified multiple times.",
		})

		//TODO(XX): deprecate and remove this
		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "env-vars",
			Target: &c.flagEnvVars,
			Usage:  "DEPRECATED. Please see `-env-var`.",
		})

		f.StringPtrVar(&flag.StringPtrVar{
			Name:   "plugin-type",
			Target: &c.flagPluginType,
			Usage:  "The type of the plugin to launch for the on-demand runner, such as aws-ecs, kubernetes, etc.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "plugin-config",
			Target:  &c.flagPluginConfig,
			Default: "",
			Usage: "Path to an hcl file that contains the configuration for the plugin. " +
				"This is only necessary when the plugin's defaults need to be adjusted for " +
				"the environment the plugin will launch the on-demand runner in.",
		})

		f.BoolPtrVar(&flag.BoolPtrVar{
			Name:   "default",
			Target: &c.flagDefault,
			Usage: "Indicates that this remote runner profile should be the default for any project that doesn't " +
				"otherwise specify its own remote runner.",
		})

		f.StringPtrVar(&flag.StringPtrVar{
			Name:   "target-runner-id",
			Target: &c.flagTargetRunnerId,
			Usage:  "ID of the runner to target for this remote runner profile.",
		})

		f.StringMapVar(&flag.StringMapVar{
			Name:   "target-runner-label",
			Target: &c.flagTargetRunnerLabels,
			Usage: "Labels on the runner to target for this remote runner profile. " +
				"e.g. `-target-runner-label=k=v`. Can be specified multiple times.",
		})

		f.BoolPtrVar(&flag.BoolPtrVar{
			Name:   "target-runner-any",
			Target: &c.flagTargetRunnerAny,
			Usage:  "Set profile to target any available runner.",
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
