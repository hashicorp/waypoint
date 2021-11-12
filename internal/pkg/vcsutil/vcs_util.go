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

func (v *VcsChecker) getRemoteName(url string) (name string, err error) {
	remotes, err := v.runVcsGitCommand("remote", "-v")
	if err != nil {
		return "", errors.Wrapf(err, "failed to list git remotes")
	}

	var matchingRemoteName string

	for _, line := range strings.Split(remotes, "\n") {
		if line == "" {
			// This always happens once
			continue
		}
		split1 := strings.Split(line, "\t")
		if len(split1) != 2 {
			// That's weird
			continue
		}
		remoteName := split1[0]
		theRestOfTheRemote := strings.Split(split1[1], " ")
		if len(theRestOfTheRemote) != 2 {
			// That's weird too
			continue
		}
		remoteUrl := theRestOfTheRemote[0]
		remoteType := theRestOfTheRemote[1]

		if url != remoteUrl {
			continue
		}

		// So the url matches, but can we fetch from it?
		// It would be weird if we couldn't.

		if remoteType != "(fetch)" {
			// TODO: it could be nice to warn if we find a remote with the right URI but wrong type
			continue
		}

		matchingRemoteName = remoteName
	}

	if matchingRemoteName == "" {
		return "", fmt.Errorf("no remote with url matching %s found", url)
	}
	return matchingRemoteName, nil
}

// will error if no remote url is found
func (v *VcsChecker) remoteHasDiff(remoteName string, remoteBranch string) (bool, error) {
	diff, err := v.runVcsGitCommand("diff", fmt.Sprintf("%s/%s", remoteName, remoteBranch))
	if err != nil {
		return false, errors.Wrapf(err, "failed to diff against remote %s", remoteName)
	}
	if len(diff) > 0 {
		return true, nil
	}
	return false, nil
}

func (v *VcsChecker) runVcsGitCommand(gitArgs ...string) (output string, err error) {
	return runGitCommand(v.log, v.path, gitArgs...)
}

func runGitCommand(log hclog.Logger, path string, gitArgs ...string) (output string, err error) {
	args := append([]string{"-C", path}, gitArgs...)
	log.Debug(fmt.Sprintf("Running this command: git %s", strings.Join(args, " ")))
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return string(out), err
}
