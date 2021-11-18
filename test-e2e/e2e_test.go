package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDockerE2E(t *testing.T) {

	// TODO(izaak): fix
	wpBinary = "waypoint"
	examplesRootDir = "/Users/izaak/dev/waypoint-examples"

	require := require.New(t)
	dockerTemplate := fmt.Sprintf("%s/docker/static", examplesRootDir)

	projectName, projectDir, err := SetupTestProject(dockerTemplate)
	require.NoError(err)
	fmt.Printf("Using project dir %s\n", projectDir)
	fmt.Printf("Using project name %s\n", projectName)

	// Leaving the project dir around for now

	wp := NewBinary(t, wpBinary, projectDir)

	wp.Run("init")

	t.Cleanup(func() {
		// TODO: clean up lingering resources (docker containers, db files, etc)
		//os.RemoveAll(projectDir)
	})

	// Doesn't matter what the output is here, just trying to get to a clean slate

	t.Run("Deploy workflow", func(t *testing.T) {
		wp.Run("build")
		wp.RunTableExpectLength(1, "artifact list")

		wp.Run("deploy -release=false")
		wp.RunTableExpectLength(1, "deployment list")

		wp.Run("release")
		// The docker plugin has no releaser
		wp.RunTableExpectLength(0, "release list")

		wp.Run("up")
		wp.RunTableExpectLength(2, "artifact list")
		wp.RunTableExpectLength(2, "deployment list")
	})

	t.Run("Cleanup", func(t *testing.T) {
		wp.Run("destroy -auto-approve")

		wp.RunWithOutput("No deployments found", "deployment list") // Should probably be empty table
		// Artifacts aren't currently destroyed during project destroy
		wp.RunTableExpectLength(0, "release list")
	})
}
