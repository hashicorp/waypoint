package vcsutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

type VCSTester struct {
	repoPath       string
	testFiles      []*os.File
	remoteRepoPath string
	remoteUrl      string
	remoteName     string
}

var testHasGit bool

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
}

func generateGitState(branchName string, t *testing.T) (VCSTester, error) {

	log := hclog.Default()

	// Create a temporary directory for our test
	td := t.TempDir()

	// Create a temporary directory for our remote
	remote := t.TempDir()

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
	if err != nil {
		return VCSTester{}, err
	}
	r.Close()

	if _, err := runGitCommand(log, td, "add", "testfile"); err != nil {
		return VCSTester{}, err
	}

	// Add another one for testing multi-line diff output
	f := td + "/second-testfile"
	rr, err := os.OpenFile(f, os.O_CREATE, 0600)
	if err != nil {
		return VCSTester{}, err
	}
	rr.Close()

	if _, err := runGitCommand(log, td, "add", "second-testfile"); err != nil {
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
		[]*os.File{r, rr},
		remote,
		remoteUrl,
		remoteName,
	}, nil
}

func TestIsDirty(t *testing.T) {
	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	hclog.Default().SetLevel(hclog.Debug)

	require := require.New(t)
	branchName := "waypoint"

	vcsTester, err := generateGitState(branchName, t)
	require.NoError(err)

	v, err := NewVcsChecker(
		hclog.Default(),
		vcsTester.repoPath,
		vcsTester.remoteUrl,
	)

	require.NoError(err)

	t.Run("Initial state is clean", func(t *testing.T) {
		dirty, err := v.RepoHasDiff(branchName)
		require.NoError(err)
		require.False(dirty)
	})

	t.Run("Creating (but not committing) a new file is dirty", func(t *testing.T) {
		file := vcsTester.repoPath + "/dirtyfile"
		r, err := os.OpenFile(file, os.O_CREATE, 0600)
		r.Close()
		require.NoError(err)

		dirty, err := v.RepoHasDiff(branchName)
		require.NoError(err)
		require.True(dirty)
	})

	t.Run("Committing a change is dirty", func(t *testing.T) {
		change := []byte("I'm changing EVERYTHING")
		err := ioutil.WriteFile(vcsTester.testFiles[0].Name(), change, 0600)
		require.NoError(err)

		_, err = runGitCommand(v.log, vcsTester.repoPath, "commit", "-am", "\"committed\"")
		require.NoError(err)

		dirty, err := v.RepoHasDiff(branchName)
		require.NoError(err)
		require.True(dirty)
	})
}

func TestGetRemoteName(t *testing.T) {
	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	hclog.Default().SetLevel(hclog.Debug)

	require := require.New(t)

	t.Run("Fail if no remotes", func(t *testing.T) {
		td, err := ioutil.TempDir("", "test")
		defer os.RemoveAll(td)
		require.NoError(err)

		log := hclog.Default()
		_, err = runGitCommand(log, td, "init")
		require.NoError(err)

		_, err = NewVcsChecker(
			log,
			td,
			"irrelevant-to-test",
		)
		require.Error(err)
		require.Contains(err.Error(), "no remotes found for repo at path")
	})

	branchName := "waypoint"
	vcsTester, err := generateGitState(branchName, t)
	require.NoError(err)

	v, err := NewVcsChecker(
		hclog.Default(),
		vcsTester.repoPath,
		vcsTester.remoteUrl,
	)
	require.NoError(err)

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

	// This test case is especially important to waypoint, as it's not unlikely that developers will use the https
	// remote for the project definition (and user/pass auth), but ssh auth locally. We need to be able to
	// detect which ssh remote corresponds to the http remote.
	t.Run("Works regardless of protocol", func(t *testing.T) {
		httpRemote := "https://git.test/testorg/testrepo.git"
		sshRemote := "git@git.test:testorg/testrepo.git"
		_, err := runGitCommand(v.log, v.path, "remote", "add", "sshRemote", sshRemote)
		require.NoError(err)

		name, err := v.getRemoteName(httpRemote)
		require.NoError(err)
		require.Equal(name, "sshRemote")

		// Again, but detecting http locally from ssh remotely. This is less important.
		_, err = runGitCommand(v.log, v.path, "remote", "remove", "sshRemote")
		require.NoError(err)

		_, err = runGitCommand(v.log, v.path, "remote", "add", "httpRemote", httpRemote)
		require.NoError(err)

		name, err = v.getRemoteName(sshRemote)
		require.NoError(err)
		require.Equal(name, "httpRemote")
	})

	t.Run("Fail if no match found", func(t *testing.T) {
		name, err := v.getRemoteName(vcsTester.remoteUrl + "noexist")
		require.Error(err)
		require.Empty(name)
		require.Contains(err.Error(), "no remote with url matching")
	})

}

func TestRemoteHasDiff(t *testing.T) {
	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	hclog.Default().SetLevel(hclog.Debug)

	require := require.New(t)
	branchName := "waypoint"

	vcsTester, err := generateGitState(branchName, t)
	require.NoError(err)

	v, err := NewVcsChecker(
		hclog.Default(),
		vcsTester.repoPath,
		vcsTester.remoteUrl,
	)

	require.NoError(err)

	t.Run("Initial state is same as remote", func(t *testing.T) {
		diff, err := v.remoteHasDiff(branchName)
		require.NoError(err)
		require.False(diff)
	})

	t.Run("Local commits differ from remote on changes", func(t *testing.T) {
		// create branch that differs from remote branch name for cross-branch comparison
		_, err = runGitCommand(v.log, vcsTester.repoPath, "checkout", "-b", "newbranch")

		change := []byte("I'm changing EVERYTHING")
		require.NoError(err)
		err := ioutil.WriteFile(vcsTester.testFiles[0].Name(), change, 0600)
		require.NoError(err)

		diff, err := v.remoteHasDiff(branchName)
		require.NoError(err)
		require.True(diff)
	})
}

func TestRemoteFileHasDiff(t *testing.T) {
	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	hclog.Default().SetLevel(hclog.Debug)

	require := require.New(t)
	branchName := "waypoint"

	vcsTester, err := generateGitState(branchName, t)
	require.NoError(err)

	v, err := NewVcsChecker(
		hclog.Default(),
		vcsTester.repoPath,
		vcsTester.remoteUrl,
	)

	require.NoError(err)

	t.Run("Initial state is same as remote", func(t *testing.T) {
		diff, err := v.remoteFileHasDiff(vcsTester.remoteName, branchName, vcsTester.testFiles[0].Name())
		require.NoError(err)
		require.False(diff)
	})

	t.Run("No diff for changes to non-named file(s)", func(t *testing.T) {
		change := []byte("I'm changing more things")
		err := ioutil.WriteFile(vcsTester.testFiles[1].Name(), change, 0600)
		require.NoError(err)

		diff, err := v.remoteFileHasDiff(vcsTester.remoteName, branchName, vcsTester.testFiles[0].Name())
		require.NoError(err)
		require.False(diff)
	})

	t.Run("Diff for changes to named file", func(t *testing.T) {
		change := []byte("I'm changing EVERYTHING")
		err := ioutil.WriteFile(vcsTester.testFiles[0].Name(), change, 0600)
		require.NoError(err)

		diff, err := v.remoteFileHasDiff(vcsTester.remoteName, branchName, vcsTester.testFiles[0].Name())
		require.NoError(err)
		require.True(diff)
	})
}

func Test_remoteConvertSSHtoHTTPS(t *testing.T) {
	require := require.New(t)
	httpRemote := "https://git.test/testorg/testrepo.git"
	sshRemote := "git@git.test:testorg/testrepo.git"

	newHttpRemote, err := remoteConvertSSHtoHTTPS(sshRemote)
	require.NoError(err)
	require.Equal(httpRemote, newHttpRemote)
}

func Test_remoteConvertHTTPStoSSH(t *testing.T) {
	require := require.New(t)
	httpRemote := "https://git.test/testorg/testrepo.git"
	sshRemote := "git@git.test:testorg/testrepo.git"

	newSSHRemote, err := remoteConvertHTTPStoSSH(httpRemote)
	require.NoError(err)
	require.Equal(sshRemote, newSSHRemote)
}
