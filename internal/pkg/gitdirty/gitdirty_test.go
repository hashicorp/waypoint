// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gitdirty

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/pkg/copy"
	"github.com/stretchr/testify/require"
)

var testHasGit bool

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
}

func TestPathIsDirty(t *testing.T) {
	tests := []struct {
		Name         string
		Fixture      string
		RemoteUrl    string
		RemoteBranch string
		FilePath     string
		Expected     bool
		ExpectedErr  string
	}{
		{
			"clean",
			"clean",
			"origin",
			"main",
			"a.txt",
			false,
			"",
		},
		{
			"uncommited change to a file is dirty",
			"uncommited-change",
			"origin",
			"main",
			"project-a/a.txt",
			true,
			"",
		},
		{
			"uncommited change to a different file is clean",
			"uncommited-change",
			"origin",
			"main",
			"README.txt",
			false,
			"",
		},
		{
			"unpushed change to a dirty directory is dirty",
			"committed-unpushed-change",
			"origin",
			"main",
			"project-a",
			true,
			"",
		},
		{
			"unpushed change to a different directory is clean",
			"committed-unpushed-change",
			"origin",
			"main",
			"project-b",
			false,
			"",
		},
	}

	log := hclog.Default()
	hclog.Default().SetLevel(hclog.Debug)

	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	for _, tt := range tests {
		name := tt.Name

		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			td, err := ioutil.TempDir("", "git")
			require.NoError(err)
			defer os.RemoveAll(td)

			// Copy our test fixture so we don't have any side effects
			path := filepath.Join("testdata", tt.Fixture)
			dstPath := filepath.Join(td, "fixture")
			require.NoError(copy.CopyDir(path, dstPath))
			path = dstPath

			testGitFixture(t, path)

			result, err := PathIsDirty(log, path, tt.RemoteUrl, tt.RemoteBranch, tt.FilePath)
			if tt.ExpectedErr != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.ExpectedErr)
			} else {
				require.Equal(tt.Expected, result)
				require.NoError(err)
			}
		})
	}
}

func TestRepoIsDirty(t *testing.T) {
	tests := []struct {
		Name         string
		Fixture      string
		RemoteUrl    string
		RemoteBranch string
		Expected     bool
		ExpectedErr  string
	}{
		{
			"clean",
			"clean",
			"origin",
			"main",
			false,
			"",
		},
		{
			"uncommited change is dirty",
			"uncommited-change",
			"origin",
			"main",
			true,
			"",
		},
		{
			"committed unpushed change is dirty",
			"committed-unpushed-change",
			"origin",
			"main",
			true,
			"",
		},
		{
			"unstaged change is dirty",
			"unstaged",
			"origin",
			"main",
			true,
			"",
		},
		{
			"no matching remote url",
			"clean",
			"git@strange-format.git.test/test",
			"main",
			true,
			"no remote with url matching",
		},
		{
			"no matching remote branch",
			"clean",
			"origin",
			"dne",
			true,
			"remote origin does not have specified branch",
		},
		{
			"No branch specified, diffs on default branch",
			"committed-unpushed-change",
			"origin",
			"",
			true,
			"",
		},
	}

	log := hclog.Default()
	hclog.Default().SetLevel(hclog.Debug)

	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	for _, tt := range tests {
		name := tt.Name

		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			td, err := ioutil.TempDir("", "git")
			require.NoError(err)
			defer os.RemoveAll(td)

			// Copy our test fixture so we don't have any side effects
			path := filepath.Join("testdata", tt.Fixture)
			dstPath := filepath.Join(td, "fixture")
			require.NoError(copy.CopyDir(path, dstPath))
			path = dstPath

			testGitFixture(t, path)

			result, err := RepoIsDirty(log, path, tt.RemoteUrl, tt.RemoteBranch)
			if tt.ExpectedErr != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.ExpectedErr)
			} else {
				require.Equal(tt.Expected, result)
				require.NoError(err)
			}
		})
	}
}

func Test_getRemoteName(t *testing.T) {
	tests := []struct {
		Name        string
		Fixture     string
		RemoteUrl   string
		Expected    string
		ExpectedErr string
	}{
		{
			"finds ssh remote name with exact match",
			"ssh-remote", // this repo has a remote url of "git@git.test:testorg/testrepo.git"
			"git@git.test:testorg/testrepo.git",
			"origin",
			"",
		},
		// It's not unlikely that developers will use the https
		// remote for the project definition (and user/pass auth),
		// but ssh auth locally. We need to be able to
		// detect which ssh remote corresponds to the http remote.
		{
			"finds remote name with ssh/https mismatch",
			"ssh-remote", // this repo has a remote url of "git@git.test:testorg/testrepo.git"
			"https://git.test/testorg/testrepo.git",
			"origin",
			"",
		},
		{
			"finds remote name with ssh/https mismatch and no .git extension on https",
			"ssh-remote", // this repo has a remote url of "git@git.test:testorg/testrepo.git"
			"https://git.test/testorg/testrepo",
			"origin",
			"",
		},
		{
			"finds https remote name with exact match",
			"https-remote", // this repo has a remote url of "git@git.test:testorg/testrepo.git"
			"https://git.test/testorg/testrepo.git",
			"origin",
			"",
		},
		{
			"finds remote name with https/ssh mismatch",
			"https-remote", // this repo has a remote url of "https://git.test/testorg/testrepo.git"
			"git@git.test:testorg/testrepo.git",
			"origin",
			"",
		},
		{
			"fails to find if remote url does not match",
			"ssh-remote", // this repo has a remote url of "https://git.test/testorg/testrepo.git"
			"git@git.test:unexpected/unexpected.git",
			"origin",
			"no remote with url matching",
		},
		{
			"finds remote name with multiple remotes",
			"multiple-remotes", // this repo has a remote url of "git@git.test:testorg/testrepo.git"
			"upstream-url",
			"upstream",
			"",
		},
		{
			"ssh-remote repo test no git@",
			"ssh-remote",
			"git.test:testorg/testrepo.git",
			"origin",
			"",
		},
		{
			"https-remote repo test no https:// so ssh",
			"https-remote",
			"git.test:testorg/testrepo.git",
			"origin",
			"",
		},
	}
	log := hclog.Default()
	hclog.Default().SetLevel(hclog.Debug)

	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	for _, tt := range tests {
		name := tt.Name

		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			td, err := ioutil.TempDir("", "git")
			require.NoError(err)
			defer os.RemoveAll(td)

			// Copy our test fixture, so we don't have any side effects
			path := filepath.Join("testdata", tt.Fixture)
			dstPath := filepath.Join(td, "fixture")
			require.NoError(copy.CopyDir(path, dstPath))
			path = dstPath

			testGitFixture(t, path)

			result, err := getRemoteName(log, path, tt.RemoteUrl)
			if tt.ExpectedErr != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.ExpectedErr)
			} else {
				require.Equal(tt.Expected, result)
				require.NoError(err)
			}
		})
	}
}

func testGitFixture(t *testing.T, path string) {
	t.Helper()

	// Look for a DOTgit
	originalGit := filepath.Join(path, "DOTgit")
	_, err := os.Stat(originalGit)
	require.NoError(t, err)

	// Rename it
	newGit := filepath.Join(path, ".git")
	require.NoError(t, os.Rename(originalGit, newGit))
	t.Cleanup(func() { os.Rename(newGit, originalGit) })

	// Look for a DOTgitignore and rename it if it exists
	originalGitignore := filepath.Join(path, "DOTgitignore")
	_, err = os.Stat(originalGitignore)
	if err == nil {
		// Rename it
		newGitignore := filepath.Join(path, ".gitignore")
		require.NoError(t, os.Rename(originalGitignore, newGitignore))
		t.Cleanup(func() { os.Rename(newGitignore, originalGitignore) })
	}
}

func Test_remoteConvertSSHtoHTTPS(t *testing.T) {
	tests := []struct {
		name        string
		httpsRemote string
		sshRemote   string
		wantErr     bool
	}{
		{
			"both normal",
			"https://git.test/testorg/testrepo.git",
			"git@git.test:testorg/testrepo.git",
			false,
		},
		{
			"no git@ for ssh",
			"https://git.test/testorg/testrepo.git",
			"git.test:testorg/testrepo.git",
			false,
		},
	}
	for _, tt := range tests {
		name := tt.name

		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			newHttpRemote, err := remoteConvertSSHtoHTTPS(tt.sshRemote)
			require.NoError(err)
			require.Equal(tt.httpsRemote, newHttpRemote)
		})
	}
}

func Test_remoteConverters(t *testing.T) {
	tests := []struct {
		name        string
		httpsRemote string
		sshRemote   string
		wantErr     bool
	}{
		{
			"standard pattern",
			"https://git.test/testorg/testrepo.git",
			"git@git.test:testorg/testrepo.git",
			false,
		},
		{
			"No trailing .git",
			"https://git.test/testorg/testrepo",
			"git@git.test:testorg/testrepo",
			false,
		},
		{
			"Invalid",
			"invalid",
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := remoteConvertHTTPStoSSH(tt.httpsRemote)
			if tt.wantErr {
				if err == nil {
					t.Errorf("remoteConvertHTTPStoSSH() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if got != tt.sshRemote {
				t.Errorf("remoteConvertHTTPStoSSH() got = %v, want %v", got, tt.sshRemote)
			}

			got, err = remoteConvertSSHtoHTTPS(tt.sshRemote)
			if tt.wantErr {
				if err == nil {
					t.Errorf("remoteConvertSSHtoHTTPS() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if got != tt.httpsRemote {
				t.Errorf("remoteConvertSSHtoHTTPS() got = %v, want %v", got, tt.httpsRemote)
			}
		})
	}
}

func Test_getDefaultBranch(t *testing.T) {
	getDefaultBranch(hclog.Default(), "/Users/izaak/dev/waypoint", "origin")
}
