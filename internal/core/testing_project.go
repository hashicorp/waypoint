package core

import (
	"context"
	"io/ioutil"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/sdk/component"
	componentmocks "github.com/mitchellh/devflow/sdk/component/mocks"
	"github.com/mitchellh/devflow/sdk/datadir"
	"github.com/mitchellh/devflow/sdk/internal-shared/mapper"
)

// TestProject returns a fully in-memory and side-effect free Project that
// can be used for testing. Additional options can be given to provide your own
// factories, configuration, etc.
func TestProject(t testing.T, opts ...Option) *Project {
	td, err := ioutil.TempDir("", "core")
	require.NoError(t, err)

	projDir, err := datadir.NewProject(td)
	require.NoError(t, err)

	defaultOpts := []Option{
		WithConfig(config.TestConfig(t, testProjectConfig)),
		WithDataDir(projDir),
	}

	// Create the default factory for all component types
	for typ := range component.TypeMap {
		f, _ := TestFactorySingle(t, typ, "test")
		defaultOpts = append(defaultOpts, WithFactory(typ, f))
	}

	p, err := NewProject(context.Background(), append(defaultOpts, opts...)...)
	require.NoError(t, err)

	return p
}

// TestFactorySingle creates a factory for the given component type and
// registers a single implementation and returns that mock. This is useful
// to create a factory for the WithFactory option that returns a mocked value
// that can be tested against.
func TestFactorySingle(t testing.T, typ component.Type, n string) (*mapper.Factory, *mock.Mock) {
	f := TestFactory(t, typ)
	c := componentmocks.ForType(typ)
	require.NotNil(t, c)
	TestFactoryRegister(t, f, n, c)

	return f, componentmocks.Mock(c)
}

// TestFactory creates a factory for the given component type.
func TestFactory(t testing.T, typ component.Type) *mapper.Factory {
	f, err := mapper.NewFactory(component.TypeMap[typ])
	require.NoError(t, err)
	return f
}

// TestFactoryRegister registers a singleton value to be returned for the
// factory for the name n.
func TestFactoryRegister(t testing.T, f *mapper.Factory, n string, v interface{}) {
	require.NoError(t, f.Register(n, func() interface{} { return v }))
}

// testProjectConfig is the default config for TestProject
const testProjectConfig = `
app "test" {
	build "test" {}
}
`
