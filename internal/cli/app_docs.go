package cli

import (
	"context"
	"fmt"
	"os"
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
	"github.com/hashicorp/waypoint/sdk/docs"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

type AppDocsCommand struct {
	*baseCommand

	flagBuiltin  bool
	flagMarkdown bool
	flagType     string
	flagPlugin   string
	flagMDX      bool
}

func (c *AppDocsCommand) basicFormat(name, ct string, doc *docs.Documentation) {
	c.ui.Output("%s (%s)", name, ct, terminal.WithHeaderStyle())

	dets := doc.Details()

	if dets.Description != "" {
		c.ui.Output("\n%s\n", dets.Description)
	}

	c.ui.Output("Interface:")

	if dets.Input != "" {
		c.ui.Output("  Input: **%s**", dets.Input)
	}

	if dets.Output != "" {
		c.ui.Output("  Output: **%s**", dets.Output)
	}

	mappers := dets.Mappers

	if len(mappers) > 0 {
		c.ui.Output("Mappers:")

		for _, m := range mappers {
			c.ui.Output("  %s\n", m.Description)
			c.ui.Output("  Input: **%s**\n", m.Input)
			c.ui.Output("  Output: **%s**\n", m.Output)
		}
	}

	for _, f := range doc.Fields() {
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
				Value: f.Summary,
			},
		})
	}
}

func (c *AppDocsCommand) markdownFormat(name, ct string, doc *docs.Documentation) {
	c.ui.Output("## %s (%s)", name, ct)

	dets := doc.Details()

	if dets.Description != "" {
		c.ui.Output("\n%s\n", dets.Description)
	}

	c.ui.Output("### Interface")

	if dets.Input != "" {
		c.ui.Output("* Input: **%s**", dets.Input)
	}

	if dets.Output != "" {
		c.ui.Output("* Output: **%s**", dets.Output)
	}

	mappers := dets.Mappers

	if len(mappers) > 0 {
		c.ui.Output("")
		c.ui.Output("### Mappers\n\n")

		for _, m := range mappers {
			c.ui.Output("#### %s\n", m.Description)
			c.ui.Output("* Input: **%s**\n", m.Input)
			c.ui.Output("* Output: **%s**\n", m.Output)
		}
	}

	c.ui.Output("### Variables")

	for _, f := range doc.Fields() {
		c.ui.Output("\n#### %s", f.Field)

		c.ui.Output("%s\n%s", f.Synopsis, f.Summary)

		c.ui.Output("\n* Type: **%s**", f.Type)

		if f.Optional {
			c.ui.Output("* __Optional__")

			if f.Default != "" {
				c.ui.Output("* Default: %s", f.Default)
			}
		}
	}
}

func (c *AppDocsCommand) humanize(s string) string {
	s = strings.TrimLeft(s, " \n\t")
	s = strings.TrimRight(s, ". \n\t")

	if s == "" {
		return s
	}

	s = strings.ToUpper(string(s[0])) + s[1:]

	s += "."

	return s
}

func (c *AppDocsCommand) mdxFormat(name, ct string, doc *docs.Documentation) {
	w, err := os.Create(fmt.Sprintf("./website/content/partials/components/%s-%s.mdx", ct, name))
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "## %s (%s)\n\n", name, ct)

	dets := doc.Details()

	if dets.Description != "" {
		fmt.Fprintf(w, "%s\n\n", c.humanize(dets.Description))
	}

	fmt.Fprintf(w, "### Interface\n\n")

	if dets.Input != "" {
		fmt.Fprintf(w, "- Input: **%s**\n", dets.Input)
	}

	if dets.Output != "" {
		fmt.Fprintf(w, "- Output: **%s**\n", dets.Output)
	}

	mappers := dets.Mappers

	if len(mappers) > 0 {
		fmt.Fprintf(w, "\n### Mappers\n\n")

		for _, m := range mappers {
			fmt.Fprintf(w, "#### %s\n", m.Description)
			fmt.Fprintf(w, "* Input: **%s**\n", m.Input)
			fmt.Fprintf(w, "* Output: **%s**\n", m.Output)
		}
	}

	fmt.Fprintf(w, "\n### Variables\n")

	for _, f := range doc.Fields() {
		fmt.Fprintf(w, "\n#### %s\n", f.Field)

		if f.Summary != "" {
			fmt.Fprintf(w, "%s\n\n%s\n", c.humanize(f.Synopsis), c.humanize(f.Summary))
		} else {
			fmt.Fprintf(w, "%s\n", c.humanize(f.Synopsis))
		}

		if f.Type != "" {
			fmt.Fprintf(w, "\n- Type: **%s**\n", f.Type)
		}

		if f.Optional {
			fmt.Fprintf(w, "- __Optional__\n")

			if f.Default != "" {
				fmt.Fprintf(w, "- Default: %s\n", f.Default)
			}
		}
	}

	if dets.Example != "" {
		fmt.Fprintf(w, "\n\n### Examples\n```\n%s\n```\n", dets.Example)
	}
}

func (c *AppDocsCommand) markdownFormatPB(name, ct string, doc *pb.Documentation) {
	c.ui.Output("## %s (%s)", name, ct)

	if doc.Description != "" {
		c.ui.Output("\n%s\n", doc.Description)
	}

	c.ui.Output("### Variables")

	var keys []string
	for k := range doc.Fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		f := doc.Fields[k]
		c.ui.Output("\n#### %s", f.Name)

		c.ui.Output("%s\n%s", f.Synopsis, f.Summary)

		c.ui.Output("\n* Type: **%s**", f.Type)

		if f.Optional {
			c.ui.Output("* __Optional__")

			if f.Default != "" {
				c.ui.Output("* Default: %s", f.Default)
			}
		}
	}
}

func (c *AppDocsCommand) builtinDocs(args []string) int {
	factories := []struct {
		f *factory.Factory
		t string
	}{
		{plugin.Builders, "builder"},
		{plugin.Registries, "registry"},
		{plugin.Platforms, "platform"},
		{plugin.Releasers, "releasemanager"},
	}

	for _, f := range factories {
		if c.flagType != "" && c.flagType != f.t {
			continue
		}

		types := f.f.Registered()
		sort.Strings(types)

		for _, t := range types {
			if c.flagPlugin != "" && c.flagPlugin != t {
				continue
			}

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

			doc, err := component.Documentation(raw)
			if err != nil {
				panic(err)
			}

			if c.flagMarkdown {
				c.markdownFormat(t, f.t, doc)
			} else {
				c.basicFormat(t, f.t, doc)
			}
		}
	}

	return 0
}

func (c *AppDocsCommand) builtinMDX(args []string) int {
	factories := []struct {
		f *factory.Factory
		t string
	}{
		{plugin.Builders, "builder"},
		{plugin.Registries, "registry"},
		{plugin.Platforms, "platform"},
		{plugin.Releasers, "releasemanager"},
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

			doc, err := component.Documentation(raw)
			if err != nil {
				panic(err)
			}

			c.mdxFormat(t, f.t, doc)
		}
	}

	return 0
}

func (c *AppDocsCommand) Run(args []string) int {
	opts := []Option{
		WithArgs(args),
		WithFlags(c.Flags()),
	}

	needCfg := true

	for _, s := range args {
		if s == "-website-mdx" {
			needCfg = false
		}

		if s == "-builtin" {
			needCfg = false
		}
	}

	if !needCfg {
		opts = append(opts, WithNoConfig())
	}

	// Initialize. If we fail, we just exit since Init handles the UI.
	err := c.Init(opts...)
	if err != nil {
		return 1
	}

	if c.flagBuiltin {
		return c.builtinDocs(args)
	}

	if c.flagMDX {
		return c.builtinMDX(args)
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

			ct := strings.ToLower(r.Component.Type.String())

			if c.flagType != "" && c.flagType != ct {
				continue
			}

			if c.flagPlugin != "" && c.flagPlugin != r.Component.Name {
				continue
			}

			if c.flagMarkdown {
				c.markdownFormatPB(r.Component.Name, ct, r.Docs)
				continue
			}

			c.ui.Output("%s (%s)", r.Component.Name, ct, terminal.WithHeaderStyle())

			if r.Docs.Description != "" {
				c.ui.Output("\n%s\n", r.Docs.Description)
			}

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

		f.BoolVar(&flag.BoolVar{
			Name:   "markdown",
			Target: &c.flagMarkdown,
			Usage:  "Show documentation in markdown format",
		})

		f.StringVar(&flag.StringVar{
			Name:   "type",
			Target: &c.flagType,
			Usage:  "Only show documentation for this type of plugin",
		})

		f.StringVar(&flag.StringVar{
			Name:   "plugin",
			Target: &c.flagPlugin,
			Usage:  "Only show documentation for plugins with this name",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "website-mdx",
			Target: &c.flagMDX,
			Usage:  "Write out builtin docs inclusion on the waypoint website",
			Hidden: true,
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
	return "Show documentation for components"
}

func (c *AppDocsCommand) Help() string {
	helpText := `
Usage: waypoint docs [options]

  Output documentation about the plugins. By default, it lists the documentation
	for the plugins configured by this project.

	The flags can change which plugins are listed and in which format.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}
