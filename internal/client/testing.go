package client

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

// TestProject returns an initialized client pointing to an in-memory test
// server. This will close automatically on test completion.
//
// This will also change the working directory to a temporary directory
// so that any side effect file creation doesn't impact the real working
// directory. If you need to use your working directory, query it before
// calling this.
func TestProject(t testing.T, opts ...Option) *Project {
	require := require.New(t)
	client := singleprocess.TestServer(t)

	ctx := context.Background()

	// Initialize our client
	result, err := New(ctx, append([]Option{
		WithClient(client),
		WithLocal(),
		WithProjectRef(&pb.Ref_Project{Project: "test_p"}),
	}, opts...)...)
	require.NoError(err)

	// Move into a temporary directory
	td := testTempDir(t)
	testChdir(t, td)

	// Create a valid waypoint configuration file
	configpkg.TestConfigFile(t, configpkg.TestSource(t))

	return result
}

// TestApp returns an app reference that can be used for testing.
func TestApp(t testing.T, c *Project) string {
	// This function doesn't do much right now, but I've kept it as a
	// function in case in the future we need to create the app in the
	// server or something.
	return "test_a"
}

func testChdir(t testing.T, dir string) {
	pwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { require.NoError(t, os.Chdir(pwd)) })
}

func testTempDir(t testing.T) string {
	dir, err := ioutil.TempDir("", "waypoint-test")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}
