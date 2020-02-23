package builtin

import (
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/pkg/mapper"

	"github.com/mitchellh/devflow/internal/builtin/docker"
	"github.com/mitchellh/devflow/internal/builtin/google"
	"github.com/mitchellh/devflow/internal/builtin/lambda"
	"github.com/mitchellh/devflow/internal/builtin/pack"
)

var (
	Builders   = mustFactory(mapper.NewFactory((*component.Builder)(nil)))
	Registries = mustFactory(mapper.NewFactory((*component.Registry)(nil)))
	Platforms  = mustFactory(mapper.NewFactory((*component.Platform)(nil)))
	Mappers    = []*mapper.Func{
		mustFunc(mapper.NewFunc(docker.PackImageMapper)),
	}
)

func init() {
	must(Builders.Register("pack", pack.NewBuilder))
	must(Builders.Register("lambda", lambda.NewBuilder))

	must(Registries.Register("docker", docker.NewRegistry))

	must(Platforms.Register("google-cloud-run", google.NewCloudRunPlatform))
	must(Platforms.Register("lambda", lambda.NewDeployer))
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

func mustFunc(f *mapper.Func, err error) *mapper.Func {
	must(err)
	return f
}
