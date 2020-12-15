package snapshot

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func TestServerSnapshot(t *testing.T) {
	ctx := context.Background()
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
}
