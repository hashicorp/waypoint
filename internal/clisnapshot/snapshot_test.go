package snapshot

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func TestServerSnapshot(t *testing.T) {
	ctx := context.Background()
	t.Run("base create and restore", func(t *testing.T) {
		require := require.New(t)
		// create our server
		client := singleprocess.TestServer(t)

		config := Config{
			Client: client,
		}

		// Create a temporary directory for our test
		td, err := ioutil.TempDir("", "test")
		require.NoError(err)
		defer os.RemoveAll(td)
		path := filepath.Join(td, "fancyserver")

		w, err := os.Create(path)
		require.NoError(err)

		err = config.WriteSnapshot(ctx, w)
		require.NoError(err)

		require.FileExists(path)

		r, err := os.Open(path)
		require.NoError(err)

		err = config.ReadSnapshot(ctx, r, false)
		require.NoError(err)

	})

	// probably need to further abstract ReadSnapshot as this test fails out
	// immediately when it hits `[WARN]  grpc: restore requested exit, closing database and exiting NOW`
	t.Run("restore with server exit", func(t *testing.T) {
		require := require.New(t)
		// start the server with restart channel
		restartCh := make(chan struct{})
		impl := singleprocess.TestImpl(t)
		client := server.TestServer(t, impl,
			server.TestWithContext(ctx),
			server.TestWithRestart(restartCh),
		)

		config := Config{
			Client: client,
		}

		// Create a temporary directory for our test
		td, err := ioutil.TempDir("", "test")
		require.NoError(err)
		defer os.RemoveAll(td)
		path := filepath.Join(td, "fancyserver")

		w, err := os.Create(path)
		require.NoError(err)

		err = config.WriteSnapshot(ctx, w)
		require.NoError(err)

		require.FileExists(path)

		r, err := os.Open(path)
		require.NoError(err)

		err = config.ReadSnapshot(ctx, r, true)

		// require.NoError(err)

	})

}
