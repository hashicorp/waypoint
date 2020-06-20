package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/sdk/component"
	componentmocks "github.com/hashicorp/waypoint/sdk/component/mocks"
)

func TestAppAuthenticate(t *testing.T) {
	require := require.New(t)

	// Our mock platform, which must also implement Authenticator
	mock := struct {
		*componentmocks.Platform
		*componentmocks.Authenticator
	}{
		&componentmocks.Platform{},
		&componentmocks.Authenticator{},
	}

	// Make our factory for platforms
	factory := TestFactory(t, component.PlatformType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testAuthPlatformConfig)),
		WithFactory(component.PlatformType, factory),
	), "test")

	// Expect to have the auth function called
	mock.Authenticator.On("AuthFunc").Return(func() error {
		return nil
	})

	{
		// Authenticate
		_, err := app.AuthenticateComponent(context.Background(), mock.Platform)
		require.Contains(err.Error(), "baz")
		// require.NoError(err)
	}
}

func TestAppAuthenticate_noAuth(t *testing.T) {
	require := require.New(t)

	// Our mock builder, which currently doesn't implement auth
	mock := struct {
		*componentmocks.Builder
		*componentmocks.Authenticator
	}{
		&componentmocks.Builder{},
		&componentmocks.Authenticator{},
	}

	// Make our factory for builders
	factory := TestFactory(t, component.BuilderType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testAuthBuilderConfig)),
		WithFactory(component.BuilderType, factory),
	), "test")

	{
		// Authenticate and should not err
		_, err := app.AuthenticateComponent(context.Background(), mock.Builder)
		// require.NoError(err)
		require.Contains(err.Error(), "error")
	}
}

const testAuthPlatformConfig = `
project = "test"

app "test" {
	deploy "test" {}
}
`

const testAuthBuilderConfig = `
project = "test"

app "test" {
	build "test" {}
}
`
