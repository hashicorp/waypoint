package plugin

import (
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/factory"

	"github.com/hashicorp/waypoint/builtin/aws/alb"
	"github.com/hashicorp/waypoint/builtin/aws/ami"
	"github.com/hashicorp/waypoint/builtin/aws/ec2"
	"github.com/hashicorp/waypoint/builtin/aws/ecr"
	"github.com/hashicorp/waypoint/builtin/aws/ecs"
	"github.com/hashicorp/waypoint/builtin/aws/lambda"
	"github.com/hashicorp/waypoint/builtin/aws/ssm"
	pluginAWSSSM "github.com/hashicorp/waypoint/builtin/aws/ssm"
	"github.com/hashicorp/waypoint/builtin/azure/aci"
	"github.com/hashicorp/waypoint/builtin/docker"
	dockerpull "github.com/hashicorp/waypoint/builtin/docker/pull"
	"github.com/hashicorp/waypoint/builtin/exec"
	"github.com/hashicorp/waypoint/builtin/files"
	"github.com/hashicorp/waypoint/builtin/google/cloudrun"
	"github.com/hashicorp/waypoint/builtin/k8s"
	k8sapply "github.com/hashicorp/waypoint/builtin/k8s/apply"
	"github.com/hashicorp/waypoint/builtin/nomad"
	"github.com/hashicorp/waypoint/builtin/nomad/jobspec"
	"github.com/hashicorp/waypoint/builtin/pack"
	"github.com/hashicorp/waypoint/builtin/tfc"
	"github.com/hashicorp/waypoint/builtin/vault"
	pluginVault "github.com/hashicorp/waypoint/builtin/vault"
)

var (
	// Builtins is the map of all available builtin plugins and their
	// options for launching them.
	Builtins = map[string][]sdk.Option{
		"files":                    files.Options,
		"pack":                     pack.Options,
		"docker":                   docker.Options,
		"docker-pull":              dockerpull.Options,
		"exec":                     exec.Options,
		"google-cloud-run":         cloudrun.Options,
		"azure-container-instance": aci.Options,
		"kubernetes":               k8s.Options,
		"kubernetes-apply":         k8sapply.Options,
		"aws-ecs":                  ecs.Options,
		"aws-ecr":                  ecr.Options,
		"nomad":                    nomad.Options,
		"nomad-jobspec":            jobspec.Options,
		"aws-ami":                  ami.Options,
		"aws-ec2":                  ec2.Options,
		"aws-alb":                  alb.Options,
		"aws-ssm":                  ssm.Options,
		"aws-lambda":               lambda.Options,
		"vault":                    vault.Options,
		"terraform-cloud":          tfc.Options,
	}

	// BaseFactories is the set of base plugin factories. This will include any
	// built-in or well-known plugins by default. This should be used as the base
	// for building any set of factories.
	BaseFactories = map[component.Type]*factory.Factory{
		component.MapperType:         mustFactory(factory.New((*interface{})(nil))),
		component.BuilderType:        mustFactory(factory.New(component.TypeMap[component.BuilderType])),
		component.RegistryType:       mustFactory(factory.New(component.TypeMap[component.RegistryType])),
		component.PlatformType:       mustFactory(factory.New(component.TypeMap[component.PlatformType])),
		component.ReleaseManagerType: mustFactory(factory.New(component.TypeMap[component.ReleaseManagerType])),
		component.ConfigSourcerType:  mustFactory(factory.New(component.TypeMap[component.ConfigSourcerType])),
		component.TaskLauncherType:   mustFactory(factory.New(component.TypeMap[component.TaskLauncherType])),
	}

	// ConfigSourcers are the list of built-in config sourcers. These will
	// eventually be moved out to exec-based plugins but for now we just
	// hardcode them. This is used by the CEB.
	ConfigSourcers = map[string]*Instance{
		"aws-ssm": {
			Component: &pluginAWSSSM.ConfigSourcer{},
		},
		"kubernetes": {
			Component: &k8s.ConfigSourcer{},
		},
		"vault": {
			Component: &pluginVault.ConfigSourcer{},
		},
		"terraform-cloud": {
			Component: &tfc.ConfigSourcer{},
		},
	}
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustFactory(f *factory.Factory, err error) *factory.Factory {
	must(err)
	return f
}
