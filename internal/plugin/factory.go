package plugin

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/mitchellh/devflow/sdk/component"
	sdkplugin "github.com/mitchellh/devflow/sdk/plugin"
)

// exePath contains the value of os.Executable. We cache the value because
// we use it a lot and subsequent calls perform syscalls.
var exePath string

func init() {
	var err error
	exePath, err = os.Executable()
	if err != nil {
		panic(err)
	}
}

// Factory returns the factory function for a plugin that is already
// represented by an *exec.Cmd.
func Factory(cmd *exec.Cmd, typ component.Type) interface{} {
	return func(log hclog.Logger) (interface{}, error) {
		config := sdkplugin.ClientConfig()
		config.Cmd = cmd
		config.Logger = log
		config.AutoMTLS = true

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

		// Request the plugin
		raw, err := rpcClient.Dispense(strings.ToLower(typ.String()))
		if err != nil {
			log.Error("error requesting plugin", "type", typ, "err", err)
			client.Kill()
			return nil, err
		}

		log.Debug("plugin successfully launched and connected")
		return raw, nil
	}
}

// BuiltinFactory creates a factory for a built-in plugin type.
func BuiltinFactory(name string, typ component.Type) interface{} {
	cmd := exec.Command(exePath, "plugin", name)
	return Factory(cmd, typ)
}
