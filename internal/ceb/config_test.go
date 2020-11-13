package ceb

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

// Test that our child process is restarted with an env var change.
func TestConfig_envVarChange(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start up the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
		server.TestWithRestart(restartCh),
	)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{
		Client: client,
		Helper: "write-env",
		HelperEnv: map[string]string{
			"HELPER_PATH": path,
			"TEST_VALUE":  "",
		},
	})

	// The child should still start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)

	// Get our deployment
	deployment, err := client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: ceb.DeploymentId(),
			},
		},
	})
	require.NoError(err)

	// Change our config
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Scope: &pb.ConfigVar_Application{
					Application: deployment.Application,
				},
				Name:  "TEST_VALUE",
				Value: &pb.ConfigVar_Static{Static: "hello"},
			},
		},
	})
	require.NoError(err)

	// The child should still start up
	var data []byte
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "hello")
	}, 5*time.Second, 10*time.Millisecond)

	// Set our config again but to the same value
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Scope: &pb.ConfigVar_Application{
					Application: deployment.Application,
				},
				Name:  "TEST_VALUE",
				Value: &pb.ConfigVar_Static{Static: "hello"},
			},
		},
	})
	require.NoError(err)

	// The child should still start up
	time.Sleep(1 * time.Second)
	data2, err := ioutil.ReadFile(path)
	require.NoError(err)
	require.Equal(data, data2)
}
