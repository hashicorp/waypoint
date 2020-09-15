package cli

import (
	"context"
	"sort"
	"strings"

	"github.com/posener/complete"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type AppDocsCommand struct {
	*baseCommand

	flagPush bool
}

func (c *AppDocsCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithSingleApp(),
	); err != nil {
		return 1
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
			Name:    "push",
			Target:  &c.flagPush,
			Default: true,
			Usage:   "Push the artifact to the configured registry.",
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
