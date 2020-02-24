package plugin

import (
	"github.com/mitchellh/devflow/sdk"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/pkg/mapper"

	"github.com/mitchellh/devflow/builtin/pack"
)

var (
	Builders   = mustFactory(mapper.NewFactory((*component.Builder)(nil)))
	Registries = mustFactory(mapper.NewFactory((*component.Registry)(nil)))
	Platforms  = mustFactory(mapper.NewFactory((*component.Platform)(nil)))

	// Builtins is the map of all available builtin plugins and their
	// options for launching them.
	Builtins = map[string][]sdk.Option{
		"pack": pack.Options,
	}
)

func init() {
	Builders.Register("pack", BuiltinFactory("pack", component.BuilderType))
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
