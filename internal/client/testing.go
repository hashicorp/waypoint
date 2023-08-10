// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
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
	// Initialize our app
	singleprocess.TestApp(t, c.Client(), c.App("test_a").Ref())

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
