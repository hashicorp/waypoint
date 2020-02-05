package core

import (
	"context"
	"testing"

	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/config"

	"github.com/mitchellh/devflow/internal/component"
	componentmocks "github.com/mitchellh/devflow/internal/component/mocks"
	"github.com/mitchellh/devflow/internal/mapper"
)

func TestNewProject(t *testing.T) {
	require := require.New(t)

	builderF := testFactory(t, component.BuilderType)
	builderM := componentmocks.ForType(component.BuilderType)
	testFactorySingle(t, builderF, "test", builderM)

	p, err := NewProject(context.Background(),
		WithConfig(config.TestConfig(t, testConfig)),
		WithFactory(component.BuilderType, builderF),
	)
	require.NoError(err)

	app, err := p.App("test")
	require.NoError(err)
	require.NotNil(app.Builder)
}

func testFactory(t *testing.T, typ component.Type) *mapper.Factory {
	f, err := mapper.NewFactory(component.TypeMap[typ])
	require.NoError(t, err)
	return f
}

func testFactorySingle(t *testing.T, f *mapper.Factory, n string, v interface{}) {
	require.NoError(t, f.Register(n, func() interface{} { return v }))
}

const testConfig = `
app "test" {
	build "test" {}
}
`
