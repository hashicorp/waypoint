package cli

import (
	"context"
	"sort"
	"strings"

	"github.com/posener/complete"

	"github.com/hashicorp/go-argmapper"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type AppDocsCommand struct {
	*baseCommand

	flagBuiltin bool
}

func (c *AppDocsCommand) builtinDocs(args []string) int {
	factories := []struct {
		f *factory.Factory
		t string
	}{
		{plugin.Builders, "builders"},
		{plugin.Registries, "registry"},
		{plugin.Platforms, "platforms"},
		{plugin.Releasers, "releasemanagers"},
	}

	for _, f := range factories {
		types := f.f.Registered()
		sort.Strings(types)

		for _, t := range types {
			fn := f.f.Func(t)
			res := fn.Call(argmapper.Typed(c.Log))
			if res.Err() != nil {
				panic(res.Err())
			}

			raw := res.Out(0)

			// If we have a plugin.Instance then we can extract other information
			// from this plugin. We accept pure factories too that don't return
			// this so we type-check here.
			if pinst, ok := raw.(*plugin.Instance); ok {
				raw = pinst.Component
				defer pinst.Close()
			}

			docs, err := component.Documentation(raw)
			if err != nil {
				panic(err)
			}

			c.ui.Output("%s (%s)", t, f.t, terminal.WithHeaderStyle())

			dets := docs.Details()

			if dets.Description != "" {
				c.ui.Output("\n%s\n", dets.Description)
			}

			for _, f := range docs.Fields() {
				c.ui.NamedValues([]terminal.NamedValue{
					{
						Name:  "Field",
						Value: f.Field,
					},
					{
						Name:  "Type",
						Value: f.Type,
					},
					{
						Name:  "Optional",
						Value: f.Optional,
					},
					{
						Name:  "Synopsis",
						Value: f.Synopsis,
					},
					{
						Name:  "Default",
						Value: f.Default,
					},
					{
						Name:  "Help",
						Value: f.Help,
					},
				})
			}
		}
	}

	return 0
}

func (c *AppDocsCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
	); err != nil {
		return 1
	}

	if c.flagBuiltin {
		return c.builtinDocs(args)
	}

	c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		docs, err := app.Docs(ctx, &pb.Job_DocsOp{})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}

		var (
			builder  *pb.Job_DocsResult_Result
			registry *pb.Job_DocsResult_Result
			platform *pb.Job_DocsResult_Result
			release  *pb.Job_DocsResult_Result
		)

		for _, r := range docs.Results {
			switch r.Component.Type {
			case pb.Component_BUILDER:
				builder = r
			case pb.Component_PLATFORM:
				platform = r
			case pb.Component_REGISTRY:
				registry = r
			case pb.Component_RELEASEMANAGER:
				release = r
			}
		}

		for _, r := range []*pb.Job_DocsResult_Result{builder, registry, platform, release} {
			if r == nil {
				continue
			}

			c.ui.Output("%s (%s)", r.Component.Name, r.Component.Type, terminal.WithHeaderStyle())

			var keys []string
			for k := range r.Docs.Fields {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				f := r.Docs.Fields[k]

				c.ui.NamedValues([]terminal.NamedValue{
					{
						Name:  "Field",
						Value: f.Name,
					},
					{
						Name:  "Type",
						Value: f.Type,
					},
					{
						Name:  "Optional",
						Value: f.Optional,
					},
					{
						Name:  "Synopsis",
						Value: f.Synopsis,
					},
					{
						Name:  "Default",
						Value: f.Default,
					},
					{
						Name:  "Help",
						Value: f.Summary,
					},
				})
			}
		}

		return nil
	})

	return 0
}

func (c *AppDocsCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:   "builtin",
			Target: &c.flagBuiltin,
			Usage:  "Show documentation on all builtin plugins",
		})
	})
}

func (c *AppDocsCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AppDocsCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AppDocsCommand) Synopsis() string {
	return "Build a new versioned artifact from source."
}

func (c *AppDocsCommand) Help() string {
	helpText := `
Usage: waypoint artifact build [options]
Alias: waypoint build

  Build a new versioned artifact from source.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
