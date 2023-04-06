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

func TestAppCanDestroyRelease(t *testing.T) {
	t.Run("can", func(t *testing.T) {
		require := require.New(t)

		// Our mock platform, which must also implement Destroyer
		mock := struct {
			*componentmocks.ReleaseManager
			*componentmocks.Destroyer
		}{
			&componentmocks.ReleaseManager{},
			&componentmocks.Destroyer{},
		}
		mock.Destroyer.On("DestroyFunc").Return(42)

		// Make our factory for platforms
		factory := TestFactory(t, component.ReleaseManagerType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testReleaseManagerConfig)),
			WithFactory(component.ReleaseManagerType, factory),
		), "test")

		require.True(app.CanDestroyRelease())
	})

	t.Run("cannot", func(t *testing.T) {
		require := require.New(t)

		// Our mock platform, which must also implement Destroyer
		mock := &componentmocks.ReleaseManager{}

		// Make our factory for platforms
		factory := TestFactory(t, component.ReleaseManagerType)
		TestFactoryRegister(t, factory, "test", mock)

		// Make our app
		app := TestApp(t, TestProject(t,
			WithConfig(config.TestConfig(t, testReleaseManagerConfig)),
			WithFactory(component.ReleaseManagerType, factory),
		), "test")

		require.False(app.CanDestroyRelease())
	})
}

func TestAppDestroyRelease_happy(t *testing.T) {
	require := require.New(t)

	// Our mock platform, which must also implement Destroyer
	mock := struct {
		*componentmocks.ReleaseManager
		*componentmocks.Destroyer
	}{
		&componentmocks.ReleaseManager{},
		&componentmocks.Destroyer{},
	}

	// Make our factory for platforms
	factory := TestFactory(t, component.ReleaseManagerType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testReleaseManagerConfig)),
		WithFactory(component.ReleaseManagerType, factory),
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
		require.NoError(app.DestroyRelease(context.Background(), &pb.Release{
			Application: app.ref,
			Workspace:   app.workspace,
			Release:     deployment,
		}))

		// Verify that we set the status properly
		resp, err := app.client.ListReleases(context.Background(), &pb.ListReleasesRequest{
			Application: app.ref,
			Workspace:   app.workspace,
		})
		require.NoError(err)
		require.Equal(pb.Operation_DESTROYED, resp.Releases[0].State)
	}

	{
		// Try with an error
		mock.Destroyer.Mock = mockpkg.Mock{}
		mock.Destroyer.On("DestroyFunc").Return(func() error {
			return fmt.Errorf("error!")
		})

		err := app.DestroyRelease(context.Background(), &pb.Release{
			Application: app.ref,
			Workspace:   app.workspace,
			Release:     deployment,
		})
		require.Error(err)
		require.Contains(err.Error(), "error")

		// Verify that we set the status properly
		resp, err := app.client.ListReleases(context.Background(), &pb.ListReleasesRequest{
			Application: app.ref,
			Workspace:   app.workspace,
			Order: &pb.OperationOrder{
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		require.NoError(err)
		require.Equal(pb.Operation_DESTROYED, resp.Releases[0].State)
	}
}

const testReleaseManagerConfig = `
project = "test"

app "test" {
	build {
		use "test" {}
	}
	deploy {
		use "test" {}
	}
	release {
		use "test" {}
	}
}
`
