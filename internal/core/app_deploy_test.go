// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Test that a DeployFunc that returns a DeclaredResource results in a resource saved on the deployment response
func TestAppDeploy_withDeclaredResource(t *testing.T) {
	require := require.New(t)

	// Make our factory for platforms
	mock := &componentmocks.Platform{}
	factory := TestFactory(t, component.PlatformType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testDeployConfig)),
		WithFactory(component.PlatformType, factory),
	), "test")

	// Setup our value
	mock.On("DeployFunc").Return(
		func(c context.Context, declaredResourcesResp *component.DeclaredResourcesResp) (component.Deployment, error) {
			declaredResourcesResp.DeclaredResources = append(declaredResourcesResp.DeclaredResources,
				&sdk.DeclaredResource{
					Name:                "test-instance",
					Type:                "instance",
					Platform:            "test-platform",
					CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE,
				})
			return &componentmocks.Deployment{}, nil
		},
	)

	push := &pb.PushedArtifact{
		Artifact: &pb.Artifact{},
	}

	deploy, err := app.Deploy(context.Background(), push)
	require.NoError(err)
	require.NotNil(deploy)
	require.NotEmpty(deploy.DeclaredResources)
}

// Test that we set the correct generation ID.
func TestAppDeploy_generation(t *testing.T) {
	t.Run("with no generation implementation", func(t *testing.T) {
		require := require.New(t)

		// Make our factory for platforms
		mock := &componentmocks.Platform{}
		factory := TestFactory(t, component.PlatformType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testDeployConfig)),
			WithFactory(component.PlatformType, factory),
		), "test")

		// Setup our value
		mock.On("DeployFunc").Return(func(context.Context) (component.Deployment, error) {
			return &componentmocks.Deployment{}, nil
		})

		push := &pb.PushedArtifact{
			Artifact: &pb.Artifact{},
		}

		var gen1 string
		{
			// Deploy
			deploy, err := app.Deploy(context.Background(), push)
			require.NoError(err)
			require.NotNil(deploy)
			require.NotEmpty(deploy.Generation)

			gen1 = deploy.Generation.Id
		}

		{
			// Deploy again, should be a different generation
			deploy, err := app.Deploy(context.Background(), push)
			require.NoError(err)
			require.NotNil(deploy)
			require.NotEmpty(deploy.Generation)
			require.NotEqual(deploy.Generation, gen1)
		}
	})

	t.Run("with an explicit generation ID", func(t *testing.T) {
		require := require.New(t)

		// Make our factory for platforms
		mock := &mockPlatformGen{}
		factory := TestFactory(t, component.PlatformType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testDeployConfig)),
			WithFactory(component.PlatformType, factory),
		), "test")

		// Setup our funcs
		mock.Generation.On("GenerationFunc").Return(func() []byte {
			return []byte("HELLO")
		})
		mock.Platform.On("DeployFunc").Return(func(context.Context) (component.Deployment, error) {
			return &componentmocks.Deployment{}, nil
		})

		push := &pb.PushedArtifact{
			Artifact: &pb.Artifact{},
		}

		var gen1 string
		{
			// Deploy
			deploy, err := app.Deploy(context.Background(), push)
			require.NoError(err)
			require.NotNil(deploy)
			require.NotEmpty(deploy.Generation)

			gen1 = deploy.Generation.Id
		}

		{
			// Deploy again, should be EQUAL
			deploy, err := app.Deploy(context.Background(), push)
			require.NoError(err)
			require.NotNil(deploy)
			require.NotEmpty(deploy.Generation)
			require.Equal(deploy.Generation.Id, gen1)
		}
	})

	t.Run("with an empty generation ID", func(t *testing.T) {
		require := require.New(t)

		// Make our factory for platforms
		mock := &mockPlatformGen{}
		factory := TestFactory(t, component.PlatformType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testDeployConfig)),
			WithFactory(component.PlatformType, factory),
		), "test")

		// Setup our funcs
		mock.Generation.On("GenerationFunc").Return(func() []byte {
			return []byte{}
		})
		mock.Platform.On("DeployFunc").Return(func(context.Context) (component.Deployment, error) {
			return &componentmocks.Deployment{}, nil
		})

		push := &pb.PushedArtifact{
			Artifact: &pb.Artifact{},
		}

		var gen1 string
		{
			// Deploy
			deploy, err := app.Deploy(context.Background(), push)
			require.NoError(err)
			require.NotNil(deploy)
			require.NotEmpty(deploy.Generation)

			gen1 = deploy.Generation.Id
		}

		{
			// Deploy again, should be different
			deploy, err := app.Deploy(context.Background(), push)
			require.NoError(err)
			require.NotNil(deploy)
			require.NotEmpty(deploy.Generation)
			require.NotEqual(deploy.Generation.Id, gen1)
		}
	})
}

type mockPlatformGen struct {
	componentmocks.Platform
	componentmocks.Generation
}

const testDeployConfig = `
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
