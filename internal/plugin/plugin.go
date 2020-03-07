package plugin

import (
	"github.com/mitchellh/devflow/sdk"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"

	"github.com/mitchellh/devflow/builtin/docker"
	"github.com/mitchellh/devflow/builtin/google"
	"github.com/mitchellh/devflow/builtin/lambda"
	"github.com/mitchellh/devflow/builtin/pack"
)

var (
	Builders   = mustFactory(mapper.NewFactory((*component.Builder)(nil)))
	Registries = mustFactory(mapper.NewFactory((*component.Registry)(nil)))
	Platforms  = mustFactory(mapper.NewFactory((*component.Platform)(nil)))

	// Builtins is the map of all available builtin plugins and their
	// options for launching them.
	Builtins = map[string][]sdk.Option{
		"pack":             pack.Options,
		"docker":           docker.Options,
		"google-cloud-run": google.CloudRunOptions,
		"lambda":           lambda.Options,
	}
)

func init() {
	Builders.Register("pack", BuiltinFactory("pack", component.BuilderType))
	Builders.Register("lambda", BuiltinFactory("lambda", component.BuilderType))

	Registries.Register("docker", BuiltinFactory("docker", component.RegistryType))

	Platforms.Register("google-cloud-run", BuiltinFactory("google-cloud-run", component.PlatformType))
	Platforms.Register("lambda", BuiltinFactory("lambda", component.PlatformType))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustFactory(f *mapper.Factory, err error) *mapper.Factory {
	must(err)
	return f
}
