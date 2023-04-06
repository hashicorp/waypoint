// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/json"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/internal-shared/protomappers"
)

// PluginRequest describes a plugin that should be setup to have its functions
// invoked. Config is used to identify the plugin by name, and Type is used to identify
// which one of the plugin types should be addressed within the plugin process.
type PluginRequest struct {
	// Config contains the information about the plugin itself. This will
	// be used to locate the plugin so it can be started.
	Config

	// The different components
	Type component.Type
	Name string

	Dir string

	ConfigData []byte
	JsonConfig bool
}

func startInstance(
	ctx context.Context, log hclog.Logger, req *PluginRequest,
) (*Instance, error) {
	// Figure out the command to execute to start to bring the plugin
	// into existance
	pluginPaths, err := DefaultPaths(req.Dir)
	if err != nil {
		return nil, err
	}

	// Look for any reattach plugins to allow debugging task launcher plugins
	var reattachPluginConfigs map[string]*goplugin.ReattachConfig
	reattachPluginsStr := os.Getenv("WP_REATTACH_PLUGINS")
	if reattachPluginsStr != "" {
		var err error
		reattachPluginConfigs, err = ParseReattachPlugins(reattachPluginsStr)
		if err != nil {
			return nil, err
		}
	}

	var fn interface{}

	// If we have a debug plugin setup use that
	if reattachConfig, ok := reattachPluginConfigs[req.Config.Name]; ok {
		log.Debug(fmt.Sprintf("plugin %s is declared as running for reattachment", req.Config.Name))
		fn = ReattachPluginFactory(reattachConfig, req.Type)
	} else {
		// Otherwise discover and launch the plugin
		cmd, err := Discover(&req.Config, pluginPaths)
		if err != nil {
			return nil, err
		}

		if cmd != nil {
			fn = Factory(cmd, req.Type)
		} else {
			// If the plugin was not found, it is only an error if
			// we don't have that plugin as a builtin plugin (for builtin plugins
			// we relaunch the existing process)

			if _, ok := Builtins[req.Config.Name]; !ok {
				return nil, fmt.Errorf("plugin %q not found", req.Config.Name)
			} else {
				log.Debug("plugin found as builtin")
				fn = BuiltinFactory(req.Config.Name, req.Type)
			}
		}
	}

	// Use argmapper to call our plugin creation function with
	// the standard context and log args. This is just to create
	// the plugin process via go-plugin, not to actually invoke
	// anything within the plugin yet.

	afn, err := argmapper.NewFunc(fn)
	if err != nil {
		return nil, err
	}

	fnResult := afn.Call(argmapper.Typed(ctx, log))
	if err := fnResult.Err(); err != nil {
		return nil, err
	}

	rawComponent := fnResult.Out(0)

	// If we have a plugin.Instance then we can extract other information
	// from this plugin. We accept pure factories too that don't return
	// this so we type-check here. This would only happen if the above code
	// used a creation function other than the *Factory functions.
	pi, ok := rawComponent.(*Instance)
	if !ok {
		pi = &Instance{
			Component: rawComponent,
			Close:     func() {},
		}
	}

	return pi, nil
}

// Plugin is the state of a created plugin to be invoked, returned by Open.
type Plugin struct {
	req      *PluginRequest
	Instance *Instance
}

// Open resolves the plugin information in req and returns a Plugin value
// to have Invoke called upon it. Open also returns the raw component
// interface, which is used to get the function value that will be invoked
func Open(
	ctx context.Context, log hclog.Logger, req *PluginRequest,
) (*Plugin, interface{}, error) {
	pi, err := startInstance(ctx, log, req)
	if err != nil {
		return nil, nil, err
	}

	var hclCtx *hcl.EvalContext

	log.Debug("decoding task plugin config")

	var (
		file  *hcl.File
		diags hcl.Diagnostics
	)

	if req.JsonConfig {
		file, diags = json.Parse(req.ConfigData, "plugin.json")
	} else {
		file, diags = hclsyntax.ParseConfig(req.ConfigData, "plugin.hcl", hcl.Pos{Line: 1, Column: 1})
	}

	if diags.HasErrors() {
		return nil, nil, diags
	}

	// Only configure plugin config if it exists
	if file.Body != nil {
		diag := component.Configure(pi.Component, file.Body, hclCtx.NewChild())
		if diag.HasErrors() {
			return nil, nil, diag
		}
	}

	return &Plugin{req: req, Instance: pi}, pi.Component, nil
}

// Close must be called to cleanup the plugin process when it is no longer
// needed.
func (p *Plugin) Close() error {
	p.Instance.Close()
	return nil
}

// Invoke calls the given fn interface{} value as an argmapper function.
// The additional args are passed to the function on invocation.
// The fn value is obtained by casting the component returned by OpenPlugin
// to a component interface (ie component.TaskLauncher)
// and then one of the *Func() functions is called on the specific type.
func (p *Plugin) Invoke(
	ctx context.Context, log hclog.Logger, fn interface{}, args ...interface{},
) (interface{}, error) {

	// We allow f to be a *mapper.Func because our plugin system creates
	// a argmapper.Func directly to manage translating the arguments in the
	// plugin to the host (which is likely executing this code).
	// TODO: test
	rawFunc, ok := fn.(*argmapper.Func)
	if !ok {
		var err error
		rawFunc, err = argmapper.NewFunc(fn, argmapper.Logger(log))
		if err != nil {
			return nil, err
		}
	}

	// Injecting protomappers is what makes our host <- argmapper -> plugin boundary
	// system work. It allows us to translate from higher level types like hclog.Logger
	// into types that can be passed over the plugin boundary (as protobuf messages),
	// and then back again into higher level types on the plugin side.
	mappers, err := argmapper.NewFuncList(protomappers.All,
		argmapper.Logger(log),
	)
	if err != nil {
		return nil, err
	}

	var amArgs []argmapper.Arg

	for _, a := range args {
		switch a := a.(type) {
		case argmapper.Arg:
			amArgs = append(amArgs, a)
		default:
			amArgs = append(amArgs, argmapper.Typed(a))
		}
	}

	// Make sure we have access to our context and logger and default args
	amArgs = append(amArgs,
		argmapper.ConverterFunc(mappers...),
		argmapper.Typed(
			ctx,
			log,
		),
	)

	log.Debug("invoking plugin function")

	// Build the chain and call it
	callResult := rawFunc.Call(amArgs...)
	if err := callResult.Err(); err != nil {
		return nil, err
	}
	raw := callResult.Out(0)

	log.Debug("invoked plugin function", "raw", raw)

	return raw, nil
}
