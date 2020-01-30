package builtin

import (
	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/mapper"

	"github.com/mitchellh/devflow/internal/builtin/pack"
)

var (
	Builders = mustFactory(mapper.NewFactory((*component.Builder)(nil)))
)

func init() {
	must(Builders.Register("pack", pack.NewBuilder))
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
