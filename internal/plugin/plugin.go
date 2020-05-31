package plugin

import (
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/sdk"
	"github.com/hashicorp/waypoint/sdk/component"

	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/builtin/google/cloudrun"
	"github.com/hashicorp/waypoint/builtin/k8s"
	"github.com/hashicorp/waypoint/builtin/lambda"
	"github.com/hashicorp/waypoint/builtin/pack"
)

var (
	Builders   = mustFactory(factory.New((*component.Builder)(nil)))
	Registries = mustFactory(factory.New((*component.Registry)(nil)))
	Platforms  = mustFactory(factory.New((*component.Platform)(nil)))
	Releasers  = mustFactory(factory.New((*component.ReleaseManager)(nil)))

	// Builtins is the map of all available builtin plugins and their
	// options for launching them.
	Builtins = map[string][]sdk.Option{
		"pack":             pack.Options,
		"docker":           docker.Options,
		"google-cloud-run": cloudrun.Options,
		"kubernetes":       k8s.Options,
		"lambda":           lambda.Options,
	}
)

func init() {
	Builders.Register("docker", BuiltinFactory("docker", component.BuilderType))
	Builders.Register("pack", BuiltinFactory("pack", component.BuilderType))
	Builders.Register("lambda", BuiltinFactory("lambda", component.BuilderType))

	Registries.Register("docker", BuiltinFactory("docker", component.RegistryType))

	Platforms.Register("google-cloud-run", BuiltinFactory("google-cloud-run", component.PlatformType))
	Platforms.Register("kubernetes", BuiltinFactory("kubernetes", component.PlatformType))
	Platforms.Register("lambda", BuiltinFactory("lambda", component.PlatformType))

	Releasers.Register("google-cloud-run", BuiltinFactory("google-cloud-run", component.ReleaseManagerType))
	Releasers.Register("kubernetes", BuiltinFactory("kubernetes", component.ReleaseManagerType))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustFactory(f *factory.Factory, err error) *factory.Factory {
	must(err)
	return f
}
