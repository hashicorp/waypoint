package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/posener/complete"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/config/funcs"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

func (c *AppDocsCommand) emitField(w io.Writer, h, out string, f *docs.FieldDocs) {
	name := f.Field

	if out != "" {
		name = out + "." + name
	}

	var parts []string

	if f.Category {
		parts = append(parts, fmt.Sprintf("%s %s (category)", h, name))
	} else {
		parts = append(parts, fmt.Sprintf("%s %s", h, name))
	}

	if f.Summary != "" {
		parts = append(parts, fmt.Sprintf("%s\n\n%s", c.humanize(f.Synopsis), c.humanize(f.Summary)))
	} else if f.Synopsis != "" {
		parts = append(parts, c.humanize(f.Synopsis))
	}

	var list bytes.Buffer
	if !f.Category && f.Type != "" {
		fmt.Fprintf(&list, "- Type: **%s**", f.Type)
	}

	if f.Optional {
		fmt.Fprintf(&list, "\n- **Optional**")

		if f.Default != "" {
			fmt.Fprintf(&list, "\n- Default: %s", f.Default)
		}
	}

	if list.Len() != 0 {
		parts = append(parts, list.String())
	}

	if sf := f.SubFields; len(sf) > 0 {
		for _, f := range sf {
			var sub bytes.Buffer
			c.emitField(&sub, h+"#", name, f)
			parts = append(parts, sub.String())
		}
	}

	for i, part := range parts {
		fmt.Fprintf(w, "%s", part)
		endingSpace(w, i, len(parts))
	}
}

func endingSpace(w io.Writer, i, tot int) {
	if i < tot-1 {
		fmt.Fprintf(w, "\n\n")
	}
}

func splitFields(fields []*docs.FieldDocs) (required, optional []*docs.FieldDocs) {
	var o, r []*docs.FieldDocs

	for _, f := range fields {
		if f.Optional {
			o = append(o, f)
		} else {
			r = append(r, f)
		}
	}

	return r, o
}

func (c *AppDocsCommand) emitSection(w io.Writer, name, use, h string, fields []*docs.FieldDocs) {
	fmt.Fprintf(w, "%s %s Parameters\n", h, name)

	if len(fields) == 0 {
		fmt.Fprintf(w, "\nThis plugin has no %s parameters.", strings.ToLower(name))
		return
	}

	if use != "" {
		fmt.Fprintf(w, "\nThese parameters are used in %s\n\n", use)
	} else {
		fmt.Fprintln(w)
	}

	for i, f := range fields {
		c.emitField(w, h+"#", "", f)
		endingSpace(w, i, len(fields))
	}
}

func (c *AppDocsCommand) mdxFormat(name, ct string, doc *docs.Documentation) {
	// we use this constnat to compare to ct for some special behavior
	const csType = "configsourcer"

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

	space := false
	if dets.Input != "" {
		fmt.Fprintf(w, "- Input: **%s**\n", dets.Input)
		space = true
	}

	if dets.Output != "" {
		fmt.Fprintf(w, "- Output: **%s**\n", dets.Output)
		space = true
	}

	if space {
		fmt.Fprintf(w, "\n")
	}

	if dets.Example != "" {
		fmt.Fprintf(w, "### Examples\n\n```hcl\n%s\n```\n\n", strings.TrimSpace(dets.Example))
	}

	mappers := dets.Mappers
	if len(mappers) > 0 {
		fmt.Fprintf(w, "### Mappers\n\n")

		for _, m := range mappers {
			fmt.Fprintf(w, "#### %s\n\n", m.Description)
			fmt.Fprintf(w, "- Input: **%s**\n", m.Input)
			fmt.Fprintf(w, "- Output: **%s**\n", m.Output)
		}

		fmt.Fprintf(w, "\n")
	}

	required, optional := splitFields(doc.Fields())

	use := "the [`use` stanza](/docs/waypoint-hcl/use) for this plugin."
	c.emitSection(w, "Required", use, "###", required)

	fmt.Fprintf(w, "\n\n")

	c.emitSection(w, "Optional", use, "###", optional)

	if fields := doc.TemplateFields(); len(fields) > 0 {
		fmt.Fprintf(w, "\n\n### Output Attributes\n")
		fmt.Fprintf(w, "\nOutput attributes can be used in your `waypoint.hcl` as [variables](/docs/waypoint-hcl/variables) via [`artifact`](/docs/waypoint-hcl/variables/artifact) or [`deploy`](/docs/waypoint-hcl/variables/deploy).\n\n")
		for i, f := range fields {
			c.emitField(w, "####", "", f)
			endingSpace(w, i, len(fields))
		}
	}

	fmt.Fprintln(w)
}

func (c *AppDocsCommand) mdxFormatConfigSourcer(name, ct string, doc *docs.Documentation) {
	w, err := os.Create(fmt.Sprintf("./website/content/partials/components/%s-%s.mdx", ct, name))
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "## %s (%s)\n\n", name, ct)

	dets := doc.Details()

	if dets.Description != "" {
		fmt.Fprintf(w, "%s\n\n", c.humanize(dets.Description))
	}

	if dets.Example != "" {
		fmt.Fprintf(w, "### Examples\n\n```hcl\n%s\n```\n\n", strings.TrimSpace(dets.Example))
	}

	mappers := dets.Mappers
	if len(mappers) > 0 {
		fmt.Fprintf(w, "### Mappers\n\n")

		for _, m := range mappers {
			fmt.Fprintf(w, "#### %s\n\n", m.Description)
			fmt.Fprintf(w, "- Input: **%s**\n", m.Input)
			fmt.Fprintf(w, "- Output: **%s**\n", m.Output)
		}

		fmt.Fprintln(w)
	}

	required, optional := splitFields(doc.RequestFields())

	use := "`configdynamic` for [dynamic configuration syncing](/docs/app-config/dynamic)."
	c.emitSection(w, "Required", use, "###", required)

	fmt.Fprintf(w, "\n\n")

	c.emitSection(w, "Optional", use, "###", optional)

	if len(doc.Fields()) > 0 {
		fmt.Fprintf(w, "\n\n### Source Parameters\n\n"+
			"The parameters below are used with `waypoint config set-source` to configure\n"+
			"the behavior this plugin. These are _not_ used in `configdynamic` calls. The\n"+
			"parameters used for `configdynamic` are in the previous section.\n\n")

		required, optional := splitFields(doc.Fields())

		c.emitSection(w, "Required Source", "", "####", required)

		fmt.Fprintf(w, "\n\n")

		c.emitSection(w, "Optional Source", "", "####", optional)
	}

	fmt.Fprintln(w)
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
	types := []component.Type{
		component.BuilderType,
		component.RegistryType,
		component.PlatformType,
		component.ReleaseManagerType,
		component.ConfigSourcerType,
	}

	docfactories := map[component.Type]*factory.Factory{}

	for _, t := range types {
		fact, err := factory.New(component.TypeMap[t])
		if err != nil {
			panic(err)
		}

		docfactories[t] = fact
	}

	for name := range plugin.Builtins {
		for _, t := range types {
			f := plugin.BuiltinFactory(name, t)
			docfactories[t].Register(name, f)
		}
	}

	factories := []struct {
		f *factory.Factory
		t string
	}{
		{docfactories[component.BuilderType], "builder"},
		{docfactories[component.RegistryType], "registry"},
		{docfactories[component.PlatformType], "platform"},
		{docfactories[component.ReleaseManagerType], "releasemanager"},
		{docfactories[component.ConfigSourcerType], "configsourcer"},
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
				continue
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
	types := []component.Type{
		component.BuilderType,
		component.RegistryType,
		component.PlatformType,
		component.ReleaseManagerType,
		component.ConfigSourcerType,
	}

	docfactories := map[component.Type]*factory.Factory{}

	for _, t := range types {
		fact, err := factory.New(component.TypeMap[t])
		if err != nil {
			panic(err)
		}

		docfactories[t] = fact
	}

	for name := range plugin.Builtins {
		for _, t := range types {
			f := plugin.BuiltinFactory(name, t)
			docfactories[t].Register(name, f)
		}
	}

	factories := []struct {
		f *factory.Factory
		t string
	}{
		{docfactories[component.BuilderType], "builder"},
		{docfactories[component.RegistryType], "registry"},
		{docfactories[component.PlatformType], "platform"},
		{docfactories[component.ReleaseManagerType], "releasemanager"},
		{docfactories[component.ConfigSourcerType], "configsourcer"},
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
				panic(err.Error())
			}

			switch f.t {
			case "configsourcer":
				c.mdxFormatConfigSourcer(t, f.t, doc)

			default:
				c.mdxFormat(t, f.t, doc)
			}
		}
	}

	return c.funcsMDX()
}

func (c *AppDocsCommand) funcsMDX() int {
	// Start with our HCL stdlib
	all := funcs.Stdlib()

	// add functions to our context
	addFuncs := func(fs map[string]function.Function) {
		for k, v := range fs {
			all[k] = v
		}
	}

	// Add some of our functions
	addFuncs(funcs.VCSGitFuncs("."))
	addFuncs(funcs.Filesystem())
	addFuncs(funcs.Encoding())
	addFuncs(funcs.Datetime())

	docs := funcs.Docs()

	var keys []string

	for k := range all {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	w, err := os.Create("./website/content/partials/funcs.mdx")
	if err != nil {
		panic(err)
	}

	defer w.Close()

	for _, k := range keys {
		fn := all[k]

		fmt.Fprintf(w, "## `%s`\n\n", k)

		var (
			args     []string
			argTypes []cty.Type
		)

		for _, p := range fn.Params() {
			if p.Name != "" {
				args = append(args, p.Name)
			} else {
				args = append(args, p.Type.FriendlyName())
			}

			argTypes = append(argTypes, p.Type)
		}

		if v := fn.VarParam(); v != nil {
			args = append(args, v.Name)
			argTypes = append(argTypes, v.Type)
		}

		rt, err := fn.ReturnType(argTypes)
		if err != nil {
			spew.Dump(argTypes)
			spew.Dump(fn)
			panic(err)
		}

		fmt.Fprintf(w, "```hcl\n%s(%s) %s\n```\n\n", k, strings.Join(args, ", "), rt.FriendlyName())
		if d, ok := docs[k]; ok {
			fmt.Fprintf(w, "%s\n\n", d)
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

	err = c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
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
	if err != nil {
		return 1
	}

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
