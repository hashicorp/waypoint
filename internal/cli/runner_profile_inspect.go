// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cli

import (
	"encoding/json"
	"fmt"

	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type RunnerProfileInspectCommand struct {
	*baseCommand

	flagJson bool
}

func (c *RunnerProfileInspectCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	var (
		name   string
		odrCfg *pb.OnDemandRunnerConfig
	)
	if len(c.args) == 0 {
		c.ui.Output("Runner profile name not specified, inspecting default profile.", terminal.WithWarningStyle())
		resp, err := c.project.Client().GetDefaultOnDemandRunnerConfig(c.Ctx, &empty.Empty{})
		if err != nil && status.Code(err) != codes.NotFound {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		if resp.Config == nil {
			c.ui.Output("No default profile found, please specify a profile name to inspect instead.", terminal.WithWarningStyle())
			return 1
		}

		odrCfg = resp.Config
	} else {
		name = c.args[0]

		resp, err := c.project.Client().GetOnDemandRunnerConfig(c.Ctx, &pb.GetOnDemandRunnerConfigRequest{
			Config: &pb.Ref_OnDemandRunnerConfig{
				Name: name,
			},
		})
		if err != nil && status.Code(err) != codes.NotFound {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}

		// Try again with arg as the ID
		if resp == nil {
			resp, err = c.project.Client().GetOnDemandRunnerConfig(c.Ctx, &pb.GetOnDemandRunnerConfigRequest{
				Config: &pb.Ref_OnDemandRunnerConfig{
					Id: name,
				},
			})
			if err != nil {
				if status.Code(err) != codes.NotFound {
					c.ui.Output("runner profile not found", terminal.WithErrorStyle())
					return 1
				}
				c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
				return 1
			}
		}

		odrCfg = resp.Config
	}

	if c.flagJson {
		data, err := protojson.MarshalOptions{
			Indent: "\t",
		}.Marshal(odrCfg)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
		fmt.Println(string(data))
		return 0
	}

	var targetRunner string
	if odrCfg.TargetRunner != nil {
		switch t := odrCfg.TargetRunner.Target.(type) {
		case *pb.Ref_Runner_Any:
			targetRunner = "*"
		case *pb.Ref_Runner_Id:
			targetRunner = t.Id.Id
		case *pb.Ref_Runner_Labels:
			s, _ := json.Marshal(t.Labels.Labels)
			targetRunner = "labels: " + string(s)
		}
	}
	c.ui.Output("Runner profile:", terminal.WithHeaderStyle())
	c.ui.NamedValues([]terminal.NamedValue{
		{
			Name: "Name", Value: odrCfg.Name,
		},
		{
			Name: "ID", Value: odrCfg.Id,
		},
		{
			Name: "Default", Value: odrCfg.Default,
		},
		{
			Name: "OCI URL", Value: odrCfg.OciUrl,
		},
		{
			Name: "Plugin Type", Value: odrCfg.PluginType,
		},
		{
			Name: "Target Runner", Value: targetRunner,
		},
		{
			Name: "Environment Variables", Value: odrCfg.EnvironmentVariables,
		},
	}, terminal.WithInfoStyle())

	if len(odrCfg.PluginConfig) > 0 {
		c.ui.Output("Additional Plugin Configuration:", terminal.WithHeaderStyle())

		// We have to do the %s here in case the plugin config contains
		// formatting chars we don't want to error.
		c.ui.Output("\n%s", string(odrCfg.PluginConfig))
	}

	return 0
}

func (c *RunnerProfileInspectCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage: "Output runner profile information as JSON. This includes " +
				"more fields since this is the complete API structure.",
		})
	})
}

func (c *RunnerProfileInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerProfileInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerProfileInspectCommand) Synopsis() string {
	return "Show detailed information about a runner profile."
}

func (c *RunnerProfileInspectCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner profile inspect <name>

  Show detailed information about a runner profile.

` + c.Flags().Help())
}
