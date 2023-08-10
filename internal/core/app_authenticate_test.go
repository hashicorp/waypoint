// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	"github.com/hashicorp/waypoint/internal/config"
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

	// Expect to have the validateAuth function called
	mock.Authenticator.On("ValidateAuthFunc").Return(func() error {
		return errors.New("foo")
	})

	// Expect to have the auth function called
	mock.Authenticator.On("AuthFunc").Return(func() error {
		return errors.New("foo")
	})

	{
		// Authenticate
		_, err := app.Auth(context.Background(), &Component{Value: mock})
		require.Contains(err.Error(), "foo")
	}
}

func TestAppAuthenticate_noAuth(t *testing.T) {
	require := require.New(t)

	// Our mock builder, which doesn't implement auth
	mock := struct {
		*componentmocks.Builder
	}{
		&componentmocks.Builder{},
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
		_, err := app.Auth(context.Background(), &Component{Value: mock.Builder})
		require.Error(err)
		require.Contains(err.Error(), "does not implement")
	}
}

const testAuthPlatformConfig = `
project = "test"

app "test" {
	build {
		use "test" {}
	}

	deploy {
		use "test" {}
	}
}
`

const testAuthBuilderConfig = `
project = "test"

app "test" {
	build {
		use "test" {}
	}

	deploy {
		use "test" {}
	}
}
`
