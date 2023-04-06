// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/opaqueany"
	mockpkg "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestAppCanDestroyDeploy(t *testing.T) {
	t.Run("can", func(t *testing.T) {
		require := require.New(t)

		// Our mock platform, which must also implement Destroyer
		mock := struct {
			*componentmocks.Platform
			*componentmocks.Destroyer
		}{
			&componentmocks.Platform{},
			&componentmocks.Destroyer{},
		}

		// Make our factory for platforms
		factory := TestFactory(t, component.PlatformType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testPlatformConfig)),
			WithFactory(component.PlatformType, factory),
		), "test")

		require.True(app.CanDestroyDeploy())
	})

	t.Run("cannot", func(t *testing.T) {
		require := require.New(t)

		// Our mock platform, which must also implement Destroyer
		mock := &componentmocks.Platform{}

		// Make our factory for platforms
		factory := TestFactory(t, component.PlatformType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testPlatformConfig)),
			WithFactory(component.PlatformType, factory),
		), "test")

		require.False(app.CanDestroyDeploy())
	})
}

func TestAppDestroyDeploy_happy(t *testing.T) {
	require := require.New(t)

	// Our mock platform, which must also implement Destroyer
	mock := struct {
		*componentmocks.Platform
		*componentmocks.Destroyer
	}{
		&componentmocks.Platform{},
		&componentmocks.Destroyer{},
	}

	// Make our factory for platforms
	factory := TestFactory(t, component.PlatformType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testPlatformConfig)),
		WithFactory(component.PlatformType, factory),
	), "test")

	// Expect to have the destroy function called
	deployment, err := opaqueany.New(&empty.Empty{})
	require.NoError(err)
	mock.Destroyer.On("DestroyFunc").Return(func(v *opaqueany.Any) error {
		if v == nil || v != deployment {
			return fmt.Errorf("value didn't match")
		}

		return nil
	})

	{
		// Destroy
		require.NoError(app.DestroyDeploy(context.Background(), &pb.Deployment{
			Application: app.ref,
			Workspace:   app.workspace,
			Deployment:  deployment,
		}))

		// Verify that we set the status properly
		resp, err := app.client.ListDeployments(context.Background(), &pb.ListDeploymentsRequest{
			Application: app.ref,
			Workspace:   app.workspace,
		})
		require.NoError(err)
		require.Equal(pb.Operation_DESTROYED, resp.Deployments[0].State)
	}

	{
		// Try with an error
		mock.Destroyer.Mock = mockpkg.Mock{}
		mock.Destroyer.On("DestroyFunc").Return(func() error {
			return fmt.Errorf("error!")
		})

		err := app.DestroyDeploy(context.Background(), &pb.Deployment{
			Application: app.ref,
			Workspace:   app.workspace,
			Deployment:  deployment,
		})
		require.Error(err)
		require.Contains(err.Error(), "error")

		// Verify that we set the status properly
		resp, err := app.client.ListDeployments(context.Background(), &pb.ListDeploymentsRequest{
			Application: app.ref,
			Workspace:   app.workspace,
			Order: &pb.OperationOrder{
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		require.NoError(err)
		require.Equal(pb.Operation_DESTROYED, resp.Deployments[0].State)
	}
}

const testPlatformConfig = `
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
