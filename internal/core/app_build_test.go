package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/sdk/component"
	componentmocks "github.com/hashicorp/waypoint/sdk/component/mocks"
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
	}
}

const testBuildConfig = `
project = "test"

app "test" {
	build "test" {}
}
`
