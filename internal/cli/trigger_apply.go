// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"fmt"
	"strings"

	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

type TriggerApplyCommand struct {
	*baseCommand

	// if true, this command will update the trigger. Set by invoking
	// `trigger update` as opposed to `trigger create`.
	Update bool

	flagTriggerName        string
	flagTriggerId          string
	flagTriggerDescription string
	flagTriggerTags        []string
	flagTriggerOperation   string
	flagTriggerNoAuth      bool

	// Operation options
	flagBuildSeq            int
	flagDeploySeq           int
	flagDisablePush         bool
	flagStatusReportDeploy  bool
	flagStatusReportRelease bool

	// Release options
	flagReleasePrune       bool
	flagReleasePruneRetain int
}

// Current supported trigger operation names.
var triggerOpValues = []string{"build", "push", "deploy", "destroy-workspace",
	"destroy-deployment", "release", "up", "init", "status-report-deploy", "status-report-release"}

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
		if !c.Update {
			c.ui.Output("Cannot specify id on create, must call 'waypoint trigger update'", terminal.WithErrorStyle())
			return 1
		}

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
	} else {
		if c.Update {
			c.ui.Output("Cannot update a trigger without specifying an id.\n\n%s",
				c.Help(), terminal.WithErrorStyle())
			return 1
		}
	}

	if diffTrigger != nil {
		// We're updating by id, so set some fields.

		if c.flagTriggerName == "" {
			c.flagTriggerName = diffTrigger.Name
		}
		if c.flagTriggerDescription == "" {
			c.flagTriggerDescription = diffTrigger.Description
		}
		if len(c.flagTriggerTags) == 0 {
			c.flagTriggerTags = diffTrigger.Tags
		}

		// Trigger target
		if diffTrigger.Workspace != nil && c.flagWorkspace == "" {
			c.flagWorkspace = diffTrigger.Workspace.Workspace
		}
		if diffTrigger.Project != nil && c.flagProject == "" {
			c.flagProject = diffTrigger.Project.Project
		}
		if diffTrigger.Application != nil && c.flagApp == "" {
			c.flagApp = diffTrigger.Application.Application
		}
	}

	// NOTE(briancain): there's probably a better way to set default workspace now
	if c.flagWorkspace == "" {
		c.flagWorkspace = "default"
	}

	createTrigger := &pb.Trigger{
		Name:          c.flagTriggerName,
		Description:   c.flagTriggerDescription,
		Tags:          c.flagTriggerTags,
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
			c.ui.Output("Must specify a build ID number for the \"push\" operation: %s",
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
		var artifact *pb.PushedArtifact
		if c.flagBuildSeq != 0 {
			artifact = &pb.PushedArtifact{
				Application: &pb.Ref_Application{
					Application: c.flagApp,
					Project:     c.flagProject,
				},
				Sequence: uint64(c.flagBuildSeq),
				Workspace: &pb.Ref_Workspace{
					Workspace: c.flagWorkspace,
				},
			}
		}

		// NOTE: nil artifact means "latest", the server will look up the latest in the DB and set it there
		createTrigger.Operation = &pb.Trigger_Deploy{
			Deploy: &pb.Job_DeployOp{
				Artifact: artifact,
			},
		}
	case c.flagTriggerOperation == "destroy-workspace":
		// NOTE(briancain): I don't think this operation actually works, takes no arguments...
		createTrigger.Operation = &pb.Trigger_Destroy{
			Destroy: &pb.Job_DestroyOp{
				Target: &pb.Job_DestroyOp_Workspace{
					Workspace: &empty.Empty{},
				},
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
	case c.flagTriggerOperation == "status-report-deploy":
		createTrigger.Operation = &pb.Trigger_StatusReport{
			StatusReport: &pb.Job_StatusReportOp{
				Target: &pb.Job_StatusReportOp_Deployment{
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
	case c.flagTriggerOperation == "status-report-release":
		createTrigger.Operation = &pb.Trigger_StatusReport{
			StatusReport: &pb.Job_StatusReportOp{
				Target: &pb.Job_StatusReportOp_Release{
					Release: &pb.Release{
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
	case c.flagTriggerOperation == "":
		if diffTrigger == nil {
			c.ui.Output("Empty operation type requested. Operation must be set with "+
				"'-op' and be one of the following values:\n%s\n\n%s",
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

	c.ui.Output("Trigger %q (%s) has been %s\n", resp.Trigger.Name, resp.Trigger.Id,
		action, terminal.WithSuccessStyle())

	triggerID := resp.Trigger.Id
	addr := strings.Split(c.clientContext.Server.Address, ":")[0]
	port := serverconfig.DefaultHTTPPort
	serverAddr := fmt.Sprintf("%s:%s", addr, port)
	serverTriggerURL := fmt.Sprintf("https://%s/v1/trigger/%s", serverAddr, triggerID)

	c.ui.Output(" Trigger ID: %s", triggerID, terminal.WithSuccessStyle())
	c.ui.Output("Trigger URL: %s", serverTriggerURL, terminal.WithSuccessStyle())

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
			Name:   "trigger-tag",
			Target: &c.flagTriggerTags,
			Usage:  "A collection of tags to apply to the trigger URL configuration. Can be specified multiple times.",
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
		fo.BoolVar(&flag.BoolVar{
			Name:    "disable-push",
			Target:  &c.flagDisablePush,
			Default: false,
			Usage:   "Disables pushing a build artifact to any configured registry for build operations.",
		})

		fo.IntVar(&flag.IntVar{
			Name:   "build-id",
			Target: &c.flagBuildSeq,
			Usage:  "The sequence number (short id) for the build to use for a deployment operation.",
		})

		fo.IntVar(&flag.IntVar{
			Name:   "deployment-id",
			Target: &c.flagDeploySeq,
			Usage:  "The sequence number (short id) for the deployment to use for a deployment operation.",
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
	if c.Update {
		return "Update a registered trigger URL."
	} else {
		return "Create and register a trigger URL."
	}
}

func (c *TriggerApplyCommand) Help() string {
	if c.Update {
		return formatHelp(`
Usage: waypoint trigger update [options]

  Update a trigger URL to Waypoint Server.

  If no sequence number is specified, the trigger will use the "latest" sequence
  for the given operation. I.e. if you create a deploy trigger with no specified
  build artifact sequence number, it will use whatever the latest artifact sequence is.

` + c.Flags().Help())
	} else {
		return formatHelp(`
Usage: waypoint trigger create [options]

  Create a trigger URL to Waypoint Server.

  If no sequence number is specified, the trigger will use the "latest" sequence
  for the given operation. I.e. if you create a deploy trigger with no specified
  build artifact sequence number, it will use whatever the latest artifact sequence is.

` + c.Flags().Help())
	}
}
