package clisnapshot

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

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "fancyserver")

	w, err := os.Create(path)
	defer w.Close()
	require.NoError(err)

	err = WriteSnapshot(ctx, client, w)
	require.NoError(err)

	require.FileExists(path)

	r, err := os.Open(path)
	defer r.Close()
	require.NoError(err)

	err = ReadSnapshot(ctx, client, r, false)
	require.NoError(err)
}
