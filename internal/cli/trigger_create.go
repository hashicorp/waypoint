package cli

import (
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type TriggerCreateCommand struct {
	*baseCommand

	flagTriggerName        string
	flagTriggerDescription string
	flagTriggerLabels      []string
	flagTriggerOperation   string
	flagTriggerNoAuth      bool

	// Operation options
	flagArtifactSeq string
	flagBuildSeq    string
	flagDeploySeq   string
	flagReleaseSeq  string
	flagDisablePush bool

	// Release options
	flagReleasePrune       bool
	flagReleasePruneRetain int
}

var triggerOpValues = []string{"build", "push", "deploy", "destroy", "release", "up", "init"}

func (c *TriggerCreateCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}

	return 0
}

func (c *TriggerCreateCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "name",
			Target:  &c.flagTriggerName,
			Default: "",
			Usage:   "The name the trigger configuration should be defined as.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "description",
			Target:  &c.flagTriggerDescription,
			Default: "",
			Usage:   "A human readable description about the trigger URL configuration.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "trigger-label",
			Target: &c.flagTriggerLabels,
			Usage:  "A collection of labels to apply to the trigger URL configuration. Can be specified multiple times.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "no-auth",
			Target:  &c.flagTriggerNoAuth,
			Default: false,
			Usage:   "If set, the trigger URL configuration will not require authentication to initiate a request.",
		})

		f.EnumSingleVar(&flag.EnumSingleVar{
			Name:   "op",
			Target: &c.flagTriggerOperation,
			Values: triggerOpValues,
			Usage:  "The operation the trigger should execute when requested.",
		})

		// Operation specific flags
		fo := set.NewSet("Operation Options")
		fo.StringVar(&flag.StringVar{
			Name:   "artifact-sequence",
			Target: &c.flagArtifactSeq,
			Usage:  "The sequence number for the artifact to use in an operation.",
		})

		fo.BoolVar(&flag.BoolVar{
			Name:    "disable-push",
			Target:  &c.flagDisablePush,
			Default: false,
			Usage:   "Disables pushing a build artifact to any configured registry for build operations.",
		})

		fo.StringVar(&flag.StringVar{
			Name:   "build-sequence",
			Target: &c.flagBuildSeq,
			Usage:  "The sequence number for the build to use in an operation.",
		})

		fo.StringVar(&flag.StringVar{
			Name:   "deployment-sequence",
			Target: &c.flagDeploySeq,
			Usage:  "The sequence number for the deployment to use in an operation.",
		})

		fo.StringVar(&flag.StringVar{
			Name:   "release-sequence",
			Target: &c.flagReleaseSeq,
			Usage:  "The sequence number for the release to use in an operation.",
		})

		// Release operation specific flags
		fro := set.NewSet("Release Operation Options")
		fro.BoolVar(&flag.BoolVar{
			Name:    "prune",
			Target:  &c.flagReleasePrune,
			Default: false,
			Usage:   "If true, will prune deployments that aren't released.",
		})

		fro.IntVar(&flag.IntVar{
			Name:   "prune-retain",
			Target: &c.flagReleasePruneRetain,
			Usage:  "This sets the number of unreleased deployments to retain when pruning.",
		})
	})
}

func (c *TriggerCreateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerCreateCommand) Synopsis() string {
	return "Generate a trigger URL and register it to Waypoint server"
}

func (c *TriggerCreateCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger create [options]

  Create and register a trigger URL to Waypoint Server.

	If no sequence number is specified, the trigger will use the "latest" sequence
	for the given operation. I.e. if you create a deploy trigger with no specified
	build sequence number, it will use whatever the latest build sequence is.

` + c.Flags().Help())
}
