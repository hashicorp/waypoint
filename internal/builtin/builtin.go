package builtin

import (
	"github.com/mitchellh/devflow/internal"
	"github.com/mitchellh/devflow/internal/mapper"

	"github.com/mitchellh/devflow/internal/builtin/pack"
)

var (
	BuilderM = mapper.NewM((*internal.Builder)(nil))
)

func init() {
	must(BuilderM.RegisterImpl("pack", (*pack.Builder)(nil)))
	must(BuilderM.RegisterMapper("pack", pack.NewBuilderFromSource))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
