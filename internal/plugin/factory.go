// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/internal-shared/pluginclient"
)

// Indicates that the process is running inside an ondemand runner. This information
// is passed to the plugins so they can configure themselves specially for this
// unique context.
var InsideODR bool

// exePath contains the value of os.Executable. We cache the value because
// we use it a lot and subsequent calls perform syscalls.
var exePath string

func init() {
	var err error
	exp := os.Getenv("WAYPOINT_BUILTIN_PLUGIN_EXE")
	if exp != "" {
		exePath, err = filepath.Abs(exp)
		if err != nil {
			panic(err)
		}
	} else {
		exePath, err = os.Executable()
		if err != nil {
			panic(err)
		}
	}
}

// Factory returns the factory function for a plugin that is already
// represented by an *exec.Cmd. This returns an *Instance and NOT the component
// interface value directly. This instance lets you more carefully manage the
// lifecycle of the plugin as well as get additional information about the
// plugin.
func Factory(cmd *exec.Cmd, typ component.Type) interface{} {
	return func(log hclog.Logger) (interface{}, error) {
		// We have to copy the command because go-plugin will set some
		// fields on it.
		cmdCopy := *cmd

		config := pluginclient.ClientConfig(log, InsideODR)
		config.Cmd = &cmdCopy
		config.Logger = log

		// Log that we're going to launch this
		log.Info("launching plugin", "type", typ, "path", cmd.Path, "args", cmd.Args)

		// Connect to the plugin
		client := plugin.NewClient(config)
		rpcClient, err := client.Client()
		if err != nil {
			log.Error("error creating plugin client", "err", err)
			client.Kill()
			return nil, err
		}

		// Request the plugin. We don't request the mapper type because
		// we handle that below.
		var raw interface{}
		if typ != component.MapperType {
			raw, err = rpcClient.Dispense(strings.ToLower(typ.String()))
			if err != nil {
				log.Error("error requesting plugin", "type", typ, "id", strings.ToLower(typ.String()), "err", err)
				client.Kill()
				return nil, err
			}
		}

		// Request the mappers
		mappers, err := pluginclient.Mappers(client)
		if err != nil {
			log.Error("error requesting plugin mappers", "err", err)
			client.Kill()
			return nil, err
		}

		log.Debug("plugin successfully launched and connected")
		return &Instance{
			Component: raw,
			Mappers:   mappers,
			Close:     func() { client.Kill() },
		}, nil
	}
}

// BuiltinFactory creates a factory for a built-in plugin type.
func BuiltinFactory(name string, typ component.Type) interface{} {
	cmd := exec.Command(exePath, "plugin", name)

	// For non-windows systems, we attach stdout/stderr as extra fds
	// so that we can get direct access to the TTY if possible for output.
	if runtime.GOOS != "windows" {
		cmd.ExtraFiles = []*os.File{os.Stdout, os.Stderr}
	}

	return Factory(cmd, typ)
}

// ReattachPluginFactory produces a provider factory that uses the passed
// reattach information to connect to go-plugin processes that are already
// running, and implements Instance against it.
func ReattachPluginFactory(reattach *plugin.ReattachConfig, typ component.Type) interface{} {
	return func(log hclog.Logger) (interface{}, error) {
		config := pluginclient.ClientConfig(log, InsideODR)
		config.Logger = log
		config.Reattach = reattach

		// Verify that we know about the requested plugin's protocol version
		plugins, ok := config.VersionedPlugins[reattach.ProtocolVersion]
		if !ok {
			err := fmt.Errorf("no plugins found for requested protocol version number %d", reattach.ProtocolVersion)
			log.Error(err.Error())
			return Instance{}, err
		}
		config.Plugins = plugins

		// Log that we're going to launch this
		log.Info("Connecting to a reattach plugin", "type", typ, "Addr", reattach.Addr.String())

		// Connect to the plugin
		client := plugin.NewClient(config)
		rpcClient, err := client.Client()
		if err != nil {
			log.Error("error creating reattach plugin client", "err", err)
			client.Kill()
			return nil, err
		}

		// Request the plugin. We don't request the mapper type because
		// we handle that below.
		var raw interface{}
		if typ != component.MapperType {
			raw, err = rpcClient.Dispense(strings.ToLower(typ.String()))
			if err != nil {
				log.Error("error requesting reattach plugin", "type", typ, "err", err)
				client.Kill()
				return nil, err
			}
		}

		// Request the mappers
		mappers, err := pluginclient.Mappers(client)
		if err != nil {
			log.Error("error requesting reattach plugin mappers", "err", err)
			client.Kill()
			return nil, err
		}

		log.Debug("successfully reattached to a plugin")
		return &Instance{
			Component: raw,
			Mappers:   mappers,
			Close:     func() { client.Kill() },
		}, nil
	}
}

// Instance is the result generated by the factory. This lets us pack
// a bit more information into plugin-launched components.
type Instance struct {
	// Component is the dispensed component
	Component interface{}

	// Mappers is the list of mappers that this plugin is providing.
	Mappers []*argmapper.Func

	// Closer is a function that should be called to clean up resources
	// associated with this plugin.
	Close func()
}
