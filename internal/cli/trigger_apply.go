package cli

import (
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type TriggerApplyCommand struct {
	*baseCommand

	flagTriggerName        string
	flagTriggerId          string
	flagTriggerDescription string
	flagTriggerLabels      []string
	flagTriggerOperation   string
	flagTriggerNoAuth      bool

	// Operation options
	flagArtifactSeq int
	flagBuildSeq    int
	flagDeploySeq   int
	flagReleaseSeq  int
	flagDisablePush bool

	// Release options
	flagReleasePrune       bool
	flagReleasePruneRetain int
}

// Current supported trigger operation names.
var triggerOpValues = []string{"build", "push", "deploy", "destroy-workspace",
	"destroy-deployment", "release", "up", "init"}

func (c *TriggerApplyCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	ctx := c.Ctx

	var diffTrigger *pb.Trigger
	if c.flagTriggerId != "" {
		// Look for an existing trigger if id specified
		respTrigger, err := c.project.Client().GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: c.flagTriggerId,
			},
		})
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		diffTrigger = respTrigger.Trigger
	}

	if diffTrigger != nil {
		// We're updating by id, so set some fields.

		if c.flagTriggerName == "" {
			c.flagTriggerName = diffTrigger.Name
		}
		if c.flagTriggerDescription == "" {
			c.flagTriggerDescription = diffTrigger.Description
		}
		if len(c.flagTriggerLabels) == 0 {
			c.flagTriggerLabels = diffTrigger.Labels
		}

		// Trigger target
		if diffTrigger.Workspace != nil {
			c.flagWorkspace = diffTrigger.Workspace.Workspace
		}
		if diffTrigger.Project != nil {
			c.flagProject = diffTrigger.Project.Project
		}
		if diffTrigger.Application != nil {
			c.flagApp = diffTrigger.Application.Application
		}
	}

	createTrigger := &pb.Trigger{
		Name:          c.flagTriggerName,
		Description:   c.flagTriggerDescription,
		Labels:        c.flagTriggerLabels,
		Authenticated: !c.flagTriggerNoAuth,
		Workspace: &pb.Ref_Workspace{
			Workspace: c.flagWorkspace,
		},
		Project: &pb.Ref_Project{
			Project: c.flagProject,
		},
		Application: &pb.Ref_Application{
			Application: c.flagApp,
			Project:     c.flagProject,
		},
	}

	// Set the operation
	switch {
	case c.flagTriggerOperation == "build":
		createTrigger.Operation = &pb.Trigger_Build{
			Build: &pb.Job_BuildOp{
				DisablePush: c.flagDisablePush,
			},
		}
	case c.flagTriggerOperation == "push":
		if c.flagBuildSeq == 0 {
			c.ui.Output("Must specify a build sequence number for the \"push\" operation: %s",
				c.Flags().Help(), terminal.WithErrorStyle())
			return 1
		}

		createTrigger.Operation = &pb.Trigger_Push{
			Push: &pb.Job_PushOp{
				Build: &pb.Build{
					Sequence: uint64(c.flagBuildSeq),
					Workspace: &pb.Ref_Workspace{
						Workspace: c.flagWorkspace,
					},
					Application: &pb.Ref_Application{
						Application: c.flagApp,
						Project:     c.flagProject,
					},
				},
			},
		}
	case c.flagTriggerOperation == "deploy":
		// TODO/FUTURE NOTE: If no sequence number is specififed (i.e. seq 0), the backend should default to using the "latest" artifact instead

		createTrigger.Operation = &pb.Trigger_Deploy{
			Deploy: &pb.Job_DeployOp{
				Artifact: &pb.PushedArtifact{
					Sequence: uint64(c.flagBuildSeq),
					Workspace: &pb.Ref_Workspace{
						Workspace: c.flagWorkspace,
					},
					Application: &pb.Ref_Application{
						Application: c.flagApp,
						Project:     c.flagProject,
					},
				},
			},
		}
	case c.flagTriggerOperation == "destroy-workspace":
		// NOTE(briancain): I don't think this operation actually works, takes no arguments...
		createTrigger.Operation = &pb.Trigger_Destroy{
			Destroy: &pb.Job_DestroyOp{
				Target: &pb.Job_DestroyOp_Workspace{},
			},
		}
	case c.flagTriggerOperation == "destroy-deployment":
		createTrigger.Operation = &pb.Trigger_Destroy{
			Destroy: &pb.Job_DestroyOp{
				Target: &pb.Job_DestroyOp_Deployment{
					Deployment: &pb.Deployment{
						Sequence: uint64(c.flagDeploySeq),
						Workspace: &pb.Ref_Workspace{
							Workspace: c.flagWorkspace,
						},
						Application: &pb.Ref_Application{
							Application: c.flagApp,
							Project:     c.flagProject,
						},
					},
				},
			},
		}
	case c.flagTriggerOperation == "release":
		// if no deployment seq is specified, the backend should default to latest
		rt := &pb.Trigger_Release{
			Release: &pb.Job_ReleaseOp{
				Deployment: &pb.Deployment{
					Sequence: uint64(c.flagDeploySeq),
					Workspace: &pb.Ref_Workspace{
						Workspace: c.flagWorkspace,
					},
					Application: &pb.Ref_Application{
						Application: c.flagApp,
						Project:     c.flagProject,
					},
				},
				Prune: c.flagReleasePrune,
			},
		}

		if c.flagReleasePruneRetain > 0 {
			rt.Release.PruneRetain = int32(c.flagReleasePruneRetain)
			rt.Release.PruneRetainOverride = true
		}

		createTrigger.Operation = rt
	case c.flagTriggerOperation == "up":
		releaseOp := &pb.Job_ReleaseOp{
			Prune: c.flagReleasePrune,
		}

		if c.flagReleasePruneRetain > 0 {
			releaseOp.PruneRetain = int32(c.flagReleasePruneRetain)
			releaseOp.PruneRetainOverride = true
		}

		createTrigger.Operation = &pb.Trigger_Up{
			Up: &pb.Job_UpOp{
				Release: releaseOp,
			},
		}
	case c.flagTriggerOperation == "init":
		createTrigger.Operation = &pb.Trigger_Init{
			Init: &pb.Job_InitOp{},
		}
	case c.flagTriggerOperation == "":
		if diffTrigger == nil {
			c.ui.Output("Empty operation type requested. Must be one of the following values:\n%s\n\n%s",
				strings.Join(triggerOpValues[:], ", "), c.Help(), terminal.WithErrorStyle())
			return 1
		} else {
			createTrigger.Operation = diffTrigger.Operation
		}
	default:
		// This shouldn't happened because the flag package should technically be handling the parsing
		// and fail if any value was not recognized in the defined Enum
		c.ui.Output("Unrecognized operation type %q: ", c.flagTriggerOperation, terminal.WithErrorStyle())
		return 1
	}

	action := "created"
	if diffTrigger != nil {
		action = "updated"
		createTrigger.Id = diffTrigger.Id
	}

	resp, err := c.project.Client().UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: createTrigger,
	})
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// TODO(briancain): update output to show trigger URL with wp server attached once http service is implemented
	c.ui.Output("Trigger %q (%s) has been %s", resp.Trigger.Name, resp.Trigger.Id,
		action, terminal.WithSuccessStyle())

	return 0
}

func (c *TriggerApplyCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "name",
			Target:  &c.flagTriggerName,
			Default: "",
			Usage:   "The name the trigger configuration should be defined as.",
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Target: &c.flagTriggerId,
			Usage: "If specified, will look up an existing trigger by this id and " +
				"attempt to update the configuration.",
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
		fo.IntVar(&flag.IntVar{
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

		fo.IntVar(&flag.IntVar{
			Name:   "build-sequence",
			Target: &c.flagBuildSeq,
			Usage:  "The sequence number for the build to use in an operation.",
		})

		fo.IntVar(&flag.IntVar{
			Name:   "deployment-sequence",
			Target: &c.flagDeploySeq,
			Usage:  "The sequence number for the deployment to use in an operation.",
		})

		fo.IntVar(&flag.IntVar{
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

func (c *TriggerApplyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *TriggerApplyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TriggerApplyCommand) Synopsis() string {
	return "Generate and Update a trigger URL and register it to Waypoint server"
}

func (c *TriggerApplyCommand) Help() string {
	return formatHelp(`
Usage: waypoint trigger apply [options]

  Create or update a trigger URL to Waypoint Server.

  If no sequence number is specified, the trigger will use the "latest" sequence
  for the given operation. I.e. if you create a deploy trigger with no specified
  build artifact sequence number, it will use whatever the latest artifact sequence is.

` + c.Flags().Help())
}
