package plugin

import (
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/sdk"
	"github.com/hashicorp/waypoint/sdk/component"

	"github.com/hashicorp/waypoint/builtin/aws/alb"
	"github.com/hashicorp/waypoint/builtin/aws/ami"
	"github.com/hashicorp/waypoint/builtin/aws/ec2"
	"github.com/hashicorp/waypoint/builtin/aws/ecr"
	"github.com/hashicorp/waypoint/builtin/aws/ecs"
	"github.com/hashicorp/waypoint/builtin/azure/aci"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/builtin/files"
	"github.com/hashicorp/waypoint/builtin/google/cloudrun"
	"github.com/hashicorp/waypoint/builtin/k8s"
	"github.com/hashicorp/waypoint/builtin/netlify"
	"github.com/hashicorp/waypoint/builtin/nomad"
	"github.com/hashicorp/waypoint/builtin/pack"
)

var (
	// Builtins is the map of all available builtin plugins and their
	// options for launching them.
	Builtins = map[string][]sdk.Option{
		"files":                    files.Options,
		"pack":                     pack.Options,
		"docker":                   docker.Options,
		"google-cloud-run":         cloudrun.Options,
		"azure-container-instance": aci.Options,
		"kubernetes":               k8s.Options,
		"netlify":                  netlify.Options,
		"aws-ecs":                  ecs.Options,
		"aws-ecr":                  ecr.Options,
		"nomad":                    nomad.Options,
		"aws-ami":                  ami.Options,
		"aws-ec2":                  ec2.Options,
		"aws-alb":                  alb.Options,
	}

	// BaseFactories is the set of base plugin factories. This will include any
	// built-in or well-known plugins by default. This should be used as the base
	// for building any set of factories.
	BaseFactories = map[component.Type]*factory.Factory{
		component.BuilderType:        mustFactory(factory.New(component.TypeMap[component.BuilderType])),
		component.RegistryType:       mustFactory(factory.New(component.TypeMap[component.RegistryType])),
		component.PlatformType:       mustFactory(factory.New(component.TypeMap[component.PlatformType])),
		component.ReleaseManagerType: mustFactory(factory.New(component.TypeMap[component.ReleaseManagerType])),
	}
)

func init() {
	b := BaseFactories[component.BuilderType]
	b.Register("docker", BuiltinFactory("docker", component.BuilderType))
	b.Register("files", BuiltinFactory("files", component.BuilderType))
	b.Register("pack", BuiltinFactory("pack", component.BuilderType))
	b.Register("aws-ami", BuiltinFactory("aws-ami", component.BuilderType))

	b = BaseFactories[component.RegistryType]
	b.Register("docker", BuiltinFactory("docker", component.RegistryType))
	b.Register("files", BuiltinFactory("files", component.RegistryType))
	b.Register("aws-ecr", BuiltinFactory("aws-ecr", component.RegistryType))

	b = BaseFactories[component.PlatformType]
	b.Register("google-cloud-run", BuiltinFactory("google-cloud-run", component.PlatformType))
	b.Register("kubernetes", BuiltinFactory("kubernetes", component.PlatformType))
	b.Register("azure-container-instance", BuiltinFactory("azure-container-instance", component.PlatformType))
	b.Register("netlify", BuiltinFactory("netlify", component.PlatformType))
	b.Register("docker", BuiltinFactory("docker", component.PlatformType))
	b.Register("aws-ecs", BuiltinFactory("aws-ecs", component.PlatformType))
	b.Register("nomad", BuiltinFactory("nomad", component.PlatformType))
	b.Register("aws-ec2", BuiltinFactory("aws-ec2", component.PlatformType))

	b = BaseFactories[component.ReleaseManagerType]
	b.Register("google-cloud-run", BuiltinFactory("google-cloud-run", component.ReleaseManagerType))
	b.Register("kubernetes", BuiltinFactory("kubernetes", component.ReleaseManagerType))
	b.Register("aws-alb", BuiltinFactory("aws-alb", component.ReleaseManagerType))
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
