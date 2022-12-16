package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/davecgh/go-spew/spew"
	"github.com/posener/complete"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/plugin"
	"github.com/hashicorp/waypoint/pkg/config/funcs"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type AppDocsCommand struct {
	*baseCommand

	flagBuiltin  bool
	flagMarkdown bool
	flagJson     bool
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

	if f.EnvVar != "" {
		fmt.Fprintf(&list, "\n- Environment Variable: **%s**", f.EnvVar)
	}

	if list.Len() != 0 {
		parts = append(parts, list.String())
	}

	nh := h + "#"
	if len(nh) > 6 {
		nh = h
	}

	if sf := f.SubFields; len(sf) > 0 {
		for _, f := range sf {
			var sub bytes.Buffer
			c.emitField(&sub, nh, name, f)
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
	// Categories and fields are both stored in FieldDocs, so if we see a category then check if any of its sub fields are not optional.
	// If so, the whole category goes in the required section of the website docs, with the optional fields still being labelled as such
	// within the category
	for _, f := range fields {
		var requiredSubfield bool
		if sf := f.SubFields; len(sf) > 0 {
			for _, fo := range sf {
				if !fo.Optional {
					requiredSubfield = true
					break
				}
			}
		}
		if requiredSubfield || (!f.Optional && !f.Category) {
			r = append(r, f)
		} else {
			o = append(o, f)
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
		nh := h + "#"
		if len(nh) > 6 {
			nh = h
		}

		c.emitField(w, nh, "", f)
		endingSpace(w, i, len(fields))
	}
}

// jsonFormat attempts to output all the data included in a a plugin's Documentation() function in the JSON file format
func (c *AppDocsCommand) jsonFormat(name, ct string, doc *docs.Documentation) {
	// we use this constant to compare to ct for some special behavior
	const csType = "configsourcer"

	w, err := os.Create(fmt.Sprintf("./embedJson/gen/%s-%s.json", ct, name))
	if err != nil {
		c.ui.Output(fmt.Sprintf("Failed to create files: %s", clierrors.Humanize(err)), terminal.StatusError)
		panic(err)
	}

	jMap := map[string]interface{}{"name": name, "type": ct}

	dets := doc.Details()
	if dets.Description != "" {
		jMap["description"] = dets.Description
	}

	if dets.Input != "" {
		jMap["input"] = dets.Input
	}

	if dets.Output != "" {
		jMap["output"] = dets.Output
	}

	if dets.Example != "" {
		jMap["example"] = strings.TrimSpace(dets.Example)
	}

	mappers := dets.Mappers
	jMap["mappers"] = mappers

	if ct == "configsourcer" {
		required, optional := splitFields(doc.RequestFields())
		jMap["requiredFields"] = required
		jMap["optionalFields"] = optional
		use := "`dynamic` for sourcing [configuration values](/waypoint/docs/app-config/dynamic) or [input variable values](/waypoint/docs/waypoint-hcl/variables/dynamic)."
		jMap["use"] = use

		if len(doc.Fields()) > 0 {
			jMap["sourceFieldsHelp"] = "Source Parameters\n" +
				"The parameters below are used with `waypoint config source-set` to configure\n" +
				"the behavior this plugin. These are _not_ used in `dynamic` calls. The\n" +
				"parameters used for `dynamic` are in the previous section.\n"

			required, optional := splitFields(doc.Fields())
			jMap["requiredSourceFields"] = required
			jMap["optionalSourceFields"] = optional
		}
	} else {
		required, optional := splitFields(doc.Fields())
		use := "the [`use` stanza](/waypoint/docs/waypoint-hcl/use) for this plugin."
		jMap["use"] = use
		jMap["requiredFields"] = required
		jMap["optionalFields"] = optional

		if fields := doc.TemplateFields(); len(fields) > 0 {
			jMap["outputAttrsHelp"] = "Output attributes can be used in your `waypoint.hcl` as [variables](/waypoint/docs/waypoint-hcl/variables) via [`artifact`](/waypoint/docs/waypoint-hcl/variables/artifact) or [`deploy`](/waypoint/docs/waypoint-hcl/variables/deploy)."
			jMap["outputAttrs"] = fields
		}
	}

	t, _ := json.MarshalIndent(jMap, "", "   ")
	fmt.Fprintf(w, "%s\n", t)
}

func (c *AppDocsCommand) mdxFormat(name, ct string, doc *docs.Documentation) {
	// we use this constant to compare to ct for some special behavior
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

	use := "the [`use` stanza](/waypoint/docs/waypoint-hcl/use) for this plugin."
	c.emitSection(w, "Required", use, "###", required)

	fmt.Fprintf(w, "\n\n")

	c.emitSection(w, "Optional", use, "###", optional)

	if fields := doc.TemplateFields(); len(fields) > 0 {
		fmt.Fprintf(w, "\n\n### Output Attributes\n")
		fmt.Fprintf(w, "\nOutput attributes can be used in your `waypoint.hcl` as [variables](/waypoint/docs/waypoint-hcl/variables) via [`artifact`](/waypoint/docs/waypoint-hcl/variables/artifact) or [`deploy`](/waypoint/docs/waypoint-hcl/variables/deploy).\n\n")
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

	use := "`dynamic` for sourcing [configuration values](/waypoint/docs/app-config/dynamic) or [input variable values](/waypoint/docs/waypoint-hcl/variables/dynamic)."
	c.emitSection(w, "Required", use, "###", required)

	fmt.Fprintf(w, "\n\n")

	c.emitSection(w, "Optional", use, "###", optional)

	if len(doc.Fields()) > 0 {
		fmt.Fprintf(w, "\n\n### Source Parameters\n\n"+
			"The parameters below are used with `waypoint config source-set` to configure\n"+
			"the behavior this plugin. These are _not_ used in `dynamic` calls. The\n"+
			"parameters used for `dynamic` are in the previous section.\n\n")

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

type pluginDocs struct {
	pluginName string
	pluginType string
	doc        *docs.Documentation
}

func getDocs(builtinPluginNames []string, log hclog.Logger) ([]*pluginDocs, error) {
	types := []component.Type{
		component.BuilderType,
		component.RegistryType,
		component.PlatformType,
		component.ReleaseManagerType,
		component.ConfigSourcerType,
		component.TaskLauncherType,
	}

	docfactories := map[component.Type]*factory.Factory{}

	for _, t := range types {
		fact, err := factory.New(component.TypeMap[t])
		if err != nil {
			return nil, err
		}

		docfactories[t] = fact
	}

	// Look for any reattach plugins
	var reattachPluginConfigs map[string]*goplugin.ReattachConfig
	reattachPluginsStr := os.Getenv("WP_REATTACH_PLUGINS")
	if reattachPluginsStr != "" {
		var err error
		reattachPluginConfigs, err = plugin.ParseReattachPlugins(reattachPluginsStr)
		if err != nil {
			return nil, err
		}
	}

	for _, name := range builtinPluginNames {
		_, ok := plugin.Builtins[name]
		if !ok {
			return nil, fmt.Errorf("Builtin plugin named %s does not exist", name)
		}
		if reattachConfig, ok := reattachPluginConfigs[name]; ok {
			log.Debug(fmt.Sprintf("plugin %s is declared as running for reattachment", name))
			for _, t := range types {
				if err := docfactories[t].Register(name, plugin.ReattachPluginFactory(reattachConfig, t)); err != nil {
					return nil, err
				}
			}
			continue
		} else {
			for _, t := range types {
				f := plugin.BuiltinFactory(name, t)
				docfactories[t].Register(name, f)
			}
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
		{docfactories[component.TaskLauncherType], "task"},
	}

	var requestedDocs []*pluginDocs
	for _, f := range factories {
		types := f.f.Registered()
		sort.Strings(types)

		for _, t := range types {

			fn := f.f.Func(t)
			res := fn.Call(argmapper.Typed(log))
			if res.Err() != nil {
				return nil, res.Err()
			}

			raw := res.Out(0)

			// If we have a plugin.Instance then we can extract other information
			// from this plugin. We accept pure factories too that don't return
			// this so we type-check here.
			cleanup := func() {}
			if pinst, ok := raw.(*plugin.Instance); ok {
				raw = pinst.Component
				cleanup = pinst.Close // must cleanup during this loop to avoid instantiating all plugins simultaneously
			}

			doc, err := component.Documentation(raw)
			if err != nil {
				continue
			}

			requestedDocs = append(requestedDocs, &pluginDocs{
				pluginName: t,
				pluginType: f.t,
				doc:        doc,
			})
			cleanup()
		}
	}
	return requestedDocs, nil
}

func (c *AppDocsCommand) builtinDocs(args []string) int {
	var pluginNames []string
	if c.flagPlugin != "" {
		pluginNames = append(pluginNames, c.flagPlugin)
	} else {
		// Use all plugins
		for pluginName := range plugin.Builtins {
			pluginNames = append(pluginNames, pluginName)
		}
	}

	pluginDocs, err := getDocs(pluginNames, c.Log)
	if err != nil {
		c.ui.Output(fmt.Sprintf("Failed to get plugin docs: %s", err), terminal.StatusError)
		return 1
	}

	for _, pluginDoc := range pluginDocs {
		if c.flagMarkdown {
			c.markdownFormat(pluginDoc.pluginName, pluginDoc.pluginType, pluginDoc.doc)
		} else if c.flagJson {
			c.jsonFormat(pluginDoc.pluginName, pluginDoc.pluginType, pluginDoc.doc)
		} else {
			c.basicFormat(pluginDoc.pluginName, pluginDoc.pluginType, pluginDoc.doc)
		}
	}

	return 0
}

func (c *AppDocsCommand) builtinJSON() int {

	var pluginNames []string
	if c.flagPlugin != "" {
		pluginNames = append(pluginNames, c.flagPlugin)
	} else {
		// Use all plugins
		for pluginName := range plugin.Builtins {
			pluginNames = append(pluginNames, pluginName)
		}
	}

	pluginDocs, err := getDocs(pluginNames, c.Log)
	if err != nil {
		c.ui.Output(fmt.Sprintf("Failed to get plugin docs: %s", clierrors.Humanize(err)), terminal.StatusError)
		return 1
	}

	for _, pluginDoc := range pluginDocs {
		c.jsonFormat(pluginDoc.pluginName, pluginDoc.pluginType, pluginDoc.doc)
	}

	return c.funcsMDX()
}

func (c *AppDocsCommand) builtinMDX() int {

	var pluginNames []string
	if c.flagPlugin != "" {
		pluginNames = append(pluginNames, c.flagPlugin)
	} else {
		// Use all plugins
		for pluginName := range plugin.Builtins {
			pluginNames = append(pluginNames, pluginName)
		}
	}

	pluginDocs, err := getDocs(pluginNames, c.Log)
	if err != nil {
		c.ui.Output(fmt.Sprintf("Failed to get plugin docs: %s", err), terminal.StatusError)
		return 1
	}

	for _, pluginDoc := range pluginDocs {
		switch pluginDoc.pluginType {
		case "configsourcer":
			c.mdxFormatConfigSourcer(pluginDoc.pluginName, pluginDoc.pluginType, pluginDoc.doc)

		default:
			c.mdxFormat(pluginDoc.pluginName, pluginDoc.pluginType, pluginDoc.doc)
		}
	}

	return c.funcsMDX()
}

func (c *AppDocsCommand) funcsMDX() int {
	var ectx hcl.EvalContext

	funcs.AddStandardFunctions(&ectx, ".")

	all := ectx.Functions

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

	if os.Getenv("WP_REATTACH_PLUGINS") != "" {
		// Currently, only waypoint runners have the logic necessary to reattach to an existing plugin.
		c.ui.Output("WP_REATTACH_PLUGINS detected, but plugin debugging is not supported with this command.", terminal.StatusError)
		// Exit immediately, as an IDE user is unlikely to notice this warning otherwise
		return 1
	}

	if c.flagBuiltin {
		return c.builtinDocs(args)
	}

	if c.flagMDX {
		return c.builtinMDX()
	}

	if c.flagJson {
		return c.builtinJSON()
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

		f.BoolVar(&flag.BoolVar{
			Name:   "json",
			Target: &c.flagJson,
			Usage:  "Generate documentation in json format",
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
