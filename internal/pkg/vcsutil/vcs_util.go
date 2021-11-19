package vcsutil

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/pkg/errors"
)

type VcsChecker struct {
	log  hclog.Logger
	path string
}

func NewVcsChecker(log hclog.Logger, path string) *VcsChecker {
	return &VcsChecker{
		log:  log,
		path: path,
	}
}

// TODO(izaak): think about http vs ssh urls
func (v *VcsChecker) IsDirty(remoteUrl string, remoteBranch string) (bool, error) {
	gitStatus, err := v.runVcsGitCommand("status", "-s")
	if len(gitStatus) > 0 {
		return true, nil
	}

	remoteName, err := v.getRemoteName(remoteUrl)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to get remote name for url %s", remoteUrl)
	}

	diff, err := v.remoteHasDiff(remoteName, remoteBranch)
	if err != nil {
		return false, err
	}

	return diff, nil
}

func (v *VcsChecker) FileHasDiff(remoteUrl string, remoteBranch string, filename string) (bool, error) {
	remoteName, err := v.getRemoteName(remoteUrl)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to get remote name for url %s", remoteUrl)
	}

	diff, err := v.remoteFileHasDiff(remoteName, remoteBranch, filename)
	if err != nil {
		return false, err
	}
	return diff, nil
}

// getRemoteName queries the repo at VcsChecker.path for all remotes, and then
// searches for the remote that matches the provided url, returning an error if
// no remote url is found
func (v *VcsChecker) getRemoteName(url string) (name string, err error) {
	remotes, err := v.runVcsGitCommand("remote", "-v")
	if err != nil {
		return "", errors.Wrapf(err, "failed to list git remotes")
	}
	if len(remotes) == 0 {
		return "", fmt.Errorf("no remotes found for repo at path %q", v.path)
	}

	var matchingRemoteName string

	for _, line := range strings.Split(remotes, "\n") {
		if line == "" {
			// This always happens once
			continue
		}
		split := strings.Split(line, "\t")
		if len(split) != 2 {
			// That's weird
			continue
		}
		remoteName := split[0]
		remoteInfo := strings.Split(split[1], " ")
		if len(remoteInfo) != 2 {
			// That's weird too
			continue
		}
		remoteUrl := remoteInfo[0]
		remoteType := remoteInfo[1]

		if url != remoteUrl {
			continue
		}

		// So the url matches, but can we fetch from it?
		// If we can't, then a lot of other things in a gitops setup
		// will fail; we'll double check though
		if remoteType != "(fetch)" {
			v.log.Warn("The git remote %q is not linked as a `fetch` source. Please tell us how you did this, we thought it was impossible.")
			continue
		}

		matchingRemoteName = remoteName
	}
	if matchingRemoteName == "" {
		return "", fmt.Errorf("no remote with url matching %s found", url)
	}

	return matchingRemoteName, nil
}

// remoteHasDiff compares the local repo to the specified remote and branch
func (v *VcsChecker) remoteHasDiff(remoteName string, remoteBranch string) (bool, error) {
	diff, err := v.runVcsGitCommand("diff", fmt.Sprintf("%s/%s", remoteName, remoteBranch))
	if err != nil {
		return false, errors.Wrapf(err, "failed to diff against remote %q on branch %q", remoteName, remoteBranch)
	}
	if len(diff) > 0 {
		return true, nil
	}
	return false, nil
}

// remoteFileHasDiff compares the named file in the local repo to the 
// specified remote and branch
func (v *VcsChecker) remoteFileHasDiff(remoteName string, remoteBranch string, filename string) (bool, error) {
	diff, err := v.runVcsGitCommand("diff", fmt.Sprintf("%s/%s", remoteName, remoteBranch), "--", filename)
	if err != nil {
		return false, errors.Wrapf(err, "failed to diff against remote %q on branch %q", remoteName, remoteBranch)
	}
	if len(diff) > 0 {
		return true, nil
	}
	return false, nil
}

// runVcsGitCommand runs a git command given the provided args against the repo
// found at VcsChecker.path
func (v *VcsChecker) runVcsGitCommand(gitArgs ...string) (output string, err error) {
	return runGitCommand(v.log, v.path, gitArgs...)
}

// runGitCommand executes a git command against the repo at the given path with
// the provided args, returning the standard output.
func runGitCommand(log hclog.Logger, path string, gitArgs ...string) (output string, err error) {
	args := append([]string{"-C", path}, gitArgs...)
	log.Debug(fmt.Sprintf("Running this command: git %s", strings.Join(args, " ")))
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return string(out), err
}
