package vcsutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/stretchr/testify/require"
)

type VCSTester struct {
	repoPath       string
	testFile       *os.File
	remoteRepoPath string
	remoteUrl      string
	remoteName     string
}

func generateGitState(branchName string) (VCSTester, error) {

	log := hclog.Default()

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	if err != nil {
		return VCSTester{}, err
	}
	// cleanup is in a separate function since we can't defer it here

	// Create a temporary directory for our remote
	remote, err := ioutil.TempDir("", "test")
	if err != nil {
		return VCSTester{}, err
	}
	// cleanup is in a separate function since we can't defer it here

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
	file := td + "/testfile"
	r, err := os.OpenFile(file, os.O_CREATE, 0600)
	r.Close()

	if _, err := runGitCommand(log, td, "add", "testfile"); err != nil {
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
		r,
		remote,
		remoteUrl,
		remoteName,
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

	t.Run("Committing a change is dirty", func(t *testing.T) {
		change := []byte("I'm changing EVERYTHING")
		err := ioutil.WriteFile(vcsTester.testFile.Name(), change, 0600)
		require.NoError(err)

		_, err = runGitCommand(v.log, vcsTester.repoPath, "commit", "-am", "\"committed\"")
		require.NoError(err)

		dirty, err := v.IsDirty(vcsTester.remoteUrl, branchName)
		require.NoError(err)
		require.True(dirty)
	})
}

func TestGetRemoteName(t *testing.T) {
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

	t.Run("Get the remote name", func(t *testing.T) {
		name, err := v.getRemoteName(vcsTester.remoteUrl)
		require.NoError(err)
		require.Equal(vcsTester.remoteName, name)
	})

	t.Run("Get the remote name if multiple remotes", func(t *testing.T) {
		// create a new dir for new remote
		td, err := ioutil.TempDir("", "test")
		defer os.RemoveAll(td)
		require.NoError(err)

		// add a remote
		_, err = runGitCommand(v.log, v.path, "remote", "add", "newremote", td)
		require.NoError(err)

		// do it again for fun and profit
		tdt, err := ioutil.TempDir("", "test")
		defer os.RemoveAll(tdt)
		require.NoError(err)

		// add a remote
		_, err = runGitCommand(v.log, v.path, "remote", "add", "znewremote", tdt)
		require.NoError(err)

		name, err := v.getRemoteName(vcsTester.remoteUrl)
		require.NoError(err)
		require.Equal(vcsTester.remoteName, name)
	})

	t.Run("Fail if no match found", func(t *testing.T) {
		name, err := v.getRemoteName(vcsTester.remoteUrl + "noexist")
		require.Error(err)
		require.Empty(name)
		require.Contains(err.Error(), "no remote with url matching")
	})

	t.Run("Fail if no remotes", func(t *testing.T) {
		td, err := ioutil.TempDir("", "test")
		defer os.RemoveAll(td)
		require.NoError(err)

		_, err = runGitCommand(v.log, td, "init")
		require.NoError(err)
		v.path = td

		name, err := v.getRemoteName("irrelevant-to-test")
		require.Error(err)
		require.Empty(name)
		require.Contains(err.Error(), "no remotes found for repo at path")
	})
}

func TestRemoteHasDiff(t *testing.T) {
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

	t.Run("Initial state is same as remote", func(t *testing.T) {
		diff, err := v.remoteHasDiff(vcsTester.remoteName, branchName)
		require.NoError(err)
		require.False(diff)
	})

	t.Run("Local commits differ from remote on changes", func(t *testing.T) {
		// create branch that differs from remote branch name for cross-branch comparison
		_, err = runGitCommand(v.log, vcsTester.repoPath, "checkout", "-b", "newbranch")

		change := []byte("I'm changing EVERYTHING")
		require.NoError(err)
		err := ioutil.WriteFile(vcsTester.testFile.Name(), change, 0600)
		require.NoError(err)

		diff, err := v.remoteHasDiff(vcsTester.remoteName, branchName)
		require.NoError(err)
		require.True(diff)
	})
}

func TestRemoteFileHasDiff(t *testing.T) {
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

	t.Run("Initial state is same as remote", func(t *testing.T) {
		diff, err := v.remoteFileHasDiff(vcsTester.remoteName, branchName, vcsTester.testFile.Name())
		require.NoError(err)
		require.False(diff)
	})

	t.Run("Diff for changes to specified file", func(t *testing.T) {
		change := []byte("I'm changing EVERYTHING")
		err := ioutil.WriteFile(vcsTester.testFile.Name(), change, 0600)
		require.NoError(err)

		diff, err := v.remoteFileHasDiff(vcsTester.remoteName, branchName, vcsTester.testFile.Name())
		require.NoError(err)
		require.True(diff)
	})

	t.Run("No diff for other local changes", func(t *testing.T) {
		file := vcsTester.repoPath + "/dirtyfile"
		r, err := os.OpenFile(file, os.O_CREATE, 0600)
		r.Close()
		require.NoError(err)

		diff, err := v.remoteFileHasDiff(vcsTester.remoteName, branchName, vcsTester.testFile.Name())
		require.NoError(err)
		require.False(diff)
	})
}
