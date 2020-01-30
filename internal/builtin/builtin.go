package builtin

import (
	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/mapper"

	"github.com/mitchellh/devflow/internal/builtin/docker"
	"github.com/mitchellh/devflow/internal/builtin/pack"
)

var (
	Builders   = mustFactory(mapper.NewFactory((*component.Builder)(nil)))
	Registries = mustFactory(mapper.NewFactory((*component.Registry)(nil)))
	Mappers    = []*mapper.Func{
		mustFunc(mapper.NewFunc(docker.PackImageMapper)),
	}
)

func init() {
	must(Builders.Register("pack", pack.NewBuilder))
	must(Registries.Register("docker", docker.NewRegistry))
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
