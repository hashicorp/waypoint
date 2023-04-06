// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"context"
	"testing"

	"github.com/hashicorp/go-argmapper"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	"github.com/hashicorp/waypoint/internal/config"
)

func TestAppBuild_happy(t *testing.T) {
	require := require.New(t)

	// Make our factory for platforms
	mock := &componentmocks.Builder{}
	factory := TestFactory(t, component.BuilderType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testBuildConfig)),
		WithFactory(component.BuilderType, factory),
		WithJobInfo(&component.JobInfo{Id: "hello"}),
	), "test")

	// Setup our value
	artifact := &componentmocks.Artifact{}
	artifact.On("Labels").Return(map[string]string{"foo": "foo"})
	mock.On("BuildFunc").Return(func() component.Artifact {
		return artifact
	})

	{
		// Destroy
		build, _, err := app.Build(context.Background())
		require.NoError(err)

		// Verify that we set the status properly
		require.Equal("foo", build.Labels["foo"])
		require.Contains(build.Labels, "waypoint/workspace")

		// Verify we have the ID set
		require.Equal("hello", build.JobId)
	}
}

// Test that labels are available
func TestAppBuild_labels(t *testing.T) {
	require := require.New(t)

	mockB := &componentmocks.Builder{}
	mockC := &componentmocks.Configurable{}
	mock := struct {
		component.Builder
		component.Configurable
	}{
		Builder:      mockB,
		Configurable: mockC,
	}

	// Make our factory for platforms
	factory := TestFactory(t, component.BuilderType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testBuildConfigLabels)),
		WithFactory(component.BuilderType, factory),
		WithJobInfo(&component.JobInfo{Id: "hello"}),
	), "test")

	// Setup our value
	artifact := &componentmocks.Artifact{}
	artifact.On("Labels").Return(map[string]string{"foo": "foo"})
	mockB.On("BuildFunc").Return(func() component.Artifact {
		return artifact
	})

	// Setup our config
	var config struct {
		Foo string `hcl:"foo,optional"`
	}
	mockC.On("Config").Return(&config, nil)

	// Build
	_, _, err := app.Build(context.Background())
	require.NoError(err)

	// Test our config value
	require.Equal("bar", config.Foo)
}

// Test that we have an argument that lets us know there is a registry.
func TestAppBuild_hasRegistry(t *testing.T) {
	t.Run("with registry", func(t *testing.T) {
		require := require.New(t)

		// Make our factory for platforms
		mock := &componentmocks.Builder{}
		factory := TestFactory(t, component.BuilderType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testBuildConfigWithRegistry)),
			WithFactory(component.BuilderType, factory),
			WithJobInfo(&component.JobInfo{Id: "hello"}),
		), "test")

		// Setup our value
		var actualVal *bool
		artifact := &componentmocks.Artifact{}
		artifact.On("Labels").Return(map[string]string{"foo": "foo"})
		mock.On("BuildFunc").Return(func(args struct {
			argmapper.Struct

			HasRegistry bool
		}) component.Artifact {
			actualVal = &args.HasRegistry
			return artifact
		})

		{
			// Build
			build, _, err := app.Build(context.Background(), BuildWithPush(false))
			require.NoError(err)

			// Verify that we set the status properly
			require.Equal("foo", build.Labels["foo"])
			require.Contains(build.Labels, "waypoint/workspace")

			// Verify we have the ID set
			require.Equal("hello", build.JobId)

			// Verify that we DID have a registry set
			require.NotNil(actualVal)
			require.True(*actualVal)
		}
	})

	t.Run("with no registry", func(t *testing.T) {
		require := require.New(t)

		// Make our factory for platforms
		mock := &componentmocks.Builder{}
		factory := TestFactory(t, component.BuilderType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testBuildConfig)),
			WithFactory(component.BuilderType, factory),
			WithJobInfo(&component.JobInfo{Id: "hello"}),
		), "test")

		// Setup our value
		var actualVal *bool
		artifact := &componentmocks.Artifact{}
		artifact.On("Labels").Return(map[string]string{"foo": "foo"})
		mock.On("BuildFunc").Return(func(args struct {
			argmapper.Struct

			HasRegistry bool
		}) component.Artifact {
			actualVal = &args.HasRegistry
			return artifact
		})

		{
			// Build
			build, _, err := app.Build(context.Background(), BuildWithPush(false))
			require.NoError(err)

			// Verify that we set the status properly
			require.Equal("foo", build.Labels["foo"])
			require.Contains(build.Labels, "waypoint/workspace")

			// Verify we have the ID set
			require.Equal("hello", build.JobId)

			// Verify that we DID have a registry set
			require.NotNil(actualVal)
			require.False(*actualVal)
		}
	})
}

const testBuildConfig = `
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

const testBuildConfigWithRegistry = `
project = "test"

app "test" {
	build {
		use "test" {}

		registry {
		  use "test" {}
		}
	}

	deploy {
		use "test" {}
	}
}
`
const testBuildConfigLabels = `
project = "test"

app "test" {
	build {
		labels = { "foo" = "bar" }
		use "test" {
			foo = labels["foo"]
		}
	}

	deploy {
		use "test" {}
	}
}
`
