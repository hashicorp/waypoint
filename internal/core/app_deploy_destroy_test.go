package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	mockpkg "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	componentmocks "github.com/hashicorp/waypoint/sdk/component/mocks"
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
	deployment := &any.Any{}
	mock.Destroyer.On("DestroyFunc").Return(func(v *any.Any) error {
		if v == nil || v != deployment {
			return fmt.Errorf("value didn't match")
		}

		return nil
	})

	// Destroy
	require.NoError(app.DestroyDeploy(context.Background(), &pb.Deployment{
		Deployment: deployment,
	}))

	// Try with an error
	mock.Destroyer.Mock = mockpkg.Mock{}
	mock.Destroyer.On("DestroyFunc").Return(func() error {
		return fmt.Errorf("error!")
	})

	err := app.DestroyDeploy(context.Background(), &pb.Deployment{})
	require.Error(err)
	require.Contains(err.Error(), "error")
}

const testPlatformConfig = `
app "test" {
	deploy "test" {}
}
`
