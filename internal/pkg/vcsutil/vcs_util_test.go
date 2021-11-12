package vcsutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDirty(t *testing.T) {
	require := require.New(t)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)

	// Create a temporary directory for our remote
	remote, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(remote)

	cmd := exec.Command("git", "-C", td, "init")
	err = cmd.Run()
	require.NoError(err)

	// add remote
	cmd = exec.Command("git", "-C", td, "remote", "add", "remote", remote+"/.git")
	err = cmd.Run()
	require.NoError(err)

	dirty, err := IsDirty(td)
	require.NoError(err)
	require.False(dirty)

	file := td + "/dirtyfile"
	r, err := os.OpenFile(file, os.O_CREATE, 0600)
	r.Close()
	require.NoError(err)

	dirty, err = IsDirty(td)
	require.NoError(err)
	require.True(dirty)
}

func TestRemotesMatchCommitted(t *testing.T) {
	require := require.New(t)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)

	// Create a temporary directory for our remote
	remote, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(remote)

	cmd := exec.Command("git", "-C", td, "init")
	err = cmd.Run()
	require.NoError(err)

	// add remote
	cmd = exec.Command("git", "-C", td, "remote", "add", "remote", remote+"/.git")
	err = cmd.Run()
	require.NoError(err)

	// check if local commits differ from remote
	match, err := remoteHasDiff(td)
	require.NoError(err)
	require.False(match)

	// commit file
	file := td + "/dirtyfile"
	r, err := os.OpenFile(file, os.O_CREATE, 0600)
	r.Close()
	require.NoError(err)
	cmd = exec.Command("git", "-C", td, "add", ".")
	err = cmd.Run()
	require.NoError(err)
	cmd = exec.Command("git", "-C", td, "commit", "-m", "'nothing much'")
	err = cmd.Run()
	require.NoError(err)

	// check if local commits differ from remote
	match, err = remoteHasDiff(td)
	require.NoError(err)
	require.True(match)
}
