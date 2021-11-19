package vcsutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/stretchr/testify/require"
)

// TODO(izaak): clean up afterwards
type VCSTester struct {
	repoPath       string
	remoteRepoPath string
	remoteUrl      string
}

func generateGitState(branchName string) (VCSTester, error) {

	log := hclog.Default()

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	if err != nil {
		return VCSTester{}, err
	}

	// Create a temporary directory for our remote
	remote, err := ioutil.TempDir("", "test")
	if err != nil {
		return VCSTester{}, err
	}

	if _, err := runGitCommand(log, td, "init", "-b", branchName); err != nil {
		return VCSTester{}, err
	}

	if _, err := runGitCommand(log, remote, "init"); err != nil {
		return VCSTester{}, err
	}

	remoteUrl := remote + "/.git"

	remoteName := "remote"
	// add remote
	if _, err := runGitCommand(log, td, "remote", "add", remoteName, remoteUrl); err != nil {
		return VCSTester{}, err
	}

	// Create a test file and commit
	file := td + "/main"
	r, err := os.OpenFile(file, os.O_CREATE, 0600)
	r.Close()

	if _, err := runGitCommand(log, td, "add", "main"); err != nil {
		return VCSTester{}, err
	}

	if _, err := runGitCommand(log, td, "commit", "-m", "'first commit'"); err != nil {
		return VCSTester{}, err
	}

	if _, err := runGitCommand(log, td, "push", remoteName, branchName); err != nil {
		return VCSTester{}, err
	}

	return VCSTester{
		td,
		remote,
		remoteUrl,
	}, nil
}

func cleanupGeneratedDirs(vcsTester VCSTester) {
	os.RemoveAll(vcsTester.repoPath)
	os.RemoveAll(vcsTester.remoteRepoPath)
}

func TestIsDirty(t *testing.T) {
	hclog.Default().SetLevel(hclog.Debug)

	require := require.New(t)
	branchName := "waypoint"

	vcsTester, err := generateGitState(branchName)
	require.NoError(err)
	defer cleanupGeneratedDirs(vcsTester)

	v := NewVcsChecker(
		hclog.Default(),
		vcsTester.repoPath,
	)

	t.Run("Initial state is clean", func(t *testing.T) {
		dirty, err := v.IsDirty(vcsTester.remoteUrl, branchName)
		require.NoError(err)
		require.False(dirty)
	})

	t.Run("Creating (but not commiting) a new file is dirty", func(t *testing.T) {
		file := vcsTester.repoPath + "/dirtyfile"
		r, err := os.OpenFile(file, os.O_CREATE, 0600)
		r.Close()
		require.NoError(err)

		dirty, err := v.IsDirty(vcsTester.remoteUrl, branchName)
		require.NoError(err)
		require.True(dirty)
	})
}

//func TestRemotesMatchCommitted(t *testing.T) {
//	require := require.New(t)
//
//	// Create a temporary directory for our test
//	td, err := ioutil.TempDir("", "test")
//	require.NoError(err)
//	defer os.RemoveAll(td)
//
//	// Create a temporary directory for our remote
//	remote, err := ioutil.TempDir("", "test")
//	require.NoError(err)
//	defer os.RemoveAll(remote)
//
//	cmd := exec.Command("git", "-C", td, "init")
//	err = cmd.Run()
//	require.NoError(err)
//
//	// add remote
//	cmd = exec.Command("git", "-C", td, "remote", "add", "remote", remote+"/.git")
//	err = cmd.Run()
//	require.NoError(err)
//
//	// check if local commits differ from remote
//	match, err := remoteHasDiff(td, remote)
//	require.NoError(err)
//	require.False(match)
//
//	// commit file
//	file := td + "/dirtyfile"
//	r, err := os.OpenFile(file, os.O_CREATE, 0600)
//	r.Close()
//	require.NoError(err)
//	cmd = exec.Command("git", "-C", td, "add", ".")
//	err = cmd.Run()
//	require.NoError(err)
//	cmd = exec.Command("git", "-C", td, "commit", "-m", "'nothing much'")
//	err = cmd.Run()
//	require.NoError(err)
//
//	// check if local commits differ from remote
//	match, err = remoteHasDiff(td, remote)
//	require.NoError(err)
//	require.True(match)
//}
