package vcsutil

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
)

var (
	sshRemoteRegexp  *regexp.Regexp
	httpRemoteRegexp *regexp.Regexp
)

func init() {
	// Regex matching http/https remotes, tokenizing the unique components for replacement.
	sshRemoteRegexp, _ = regexp.Compile(`git@(.*?):(.*\.git)`)                 // Example: git@git.test:testorg/testrepo.git
	httpRemoteRegexp, _ = regexp.Compile(`http[s]?:\/\/(.*?\..*?)\/(.*\.git)`) // Example: https://git.test/testorg/testrepo.git
}

type VcsChecker struct {
	log hclog.Logger

	// path is the local path to the cloned git repo
	path string

	// remoteName is the name of the remote that corresponds to the url
	remoteName string

	// internal go-git repo
	repo *git.Repository
}

// NewVcsChecker creates a new configured VcsChecker.
// NOTE: VcsChecker subprocceses to the `git` command. Git must be installed.
func NewVcsChecker(log hclog.Logger, path string, remoteUrl string) (*VcsChecker, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, errors.New("command 'git' not installed")
	}

	v := &VcsChecker{
		log:  log,
		path: path,
	}

	var err error
	v.repo, err = git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open git repo at path %s", path)
	}

	remoteName, err := v.getRemoteName(remoteUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get remote name from url %s", remoteUrl)
	}
	v.remoteName = remoteName

	return v, nil
}

// RepoHasDiff looks for unstaged, staged, and committed (but not pushed)
// differences between the local VcsChecker.path repo and the specified
// remote url and branch
func (v *VcsChecker) RepoHasDiff(remoteBranch string) (bool, error) {
	gitStatus, err := v.runVcsGitCommand("status", "-s")
	if len(gitStatus) > 0 {
		return true, nil
	}

	diff, err := v.remoteHasDiff(remoteBranch)
	if err != nil {
		return false, err
	}

	return diff, nil
}

// FileHasDiff checks only the specified file for unstaged, staged, and
// committed (but not pushed) differences between the local VcsChecker.path
// repo and the specified remote url and branch
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

func isSSHRemote(remote string) bool {
	return sshRemoteRegexp.MatchString(remote)
}

func isHTTPSRemote(remote string) bool {
	return httpRemoteRegexp.MatchString(remote)
}

// remoteConvertHTTPStoSSH converts an https-style remote into its corresponding ssh style remote.
// Based on regex, and may not match every possible style of remote, but tested on github and gitlab.
//    Example input: https://git.test/testorg/testrepo.git
//           output: git@git.test:testorg/testrepo.git
func remoteConvertHTTPStoSSH(httpsRemote string) (string, error) {
	if !isHTTPSRemote(httpsRemote) {
		return "", fmt.Errorf("%s is not an https remote", httpsRemote)
	}

	sshRemote := httpRemoteRegexp.ReplaceAllString(httpsRemote, "git@$1:$2")
	if !isSSHRemote(sshRemote) {
		return "", fmt.Errorf("failed to convert https remote %s to ssh remote: got %s, which is not valid", httpsRemote, sshRemote)
	}
	return sshRemote, nil
}

// remoteConvertSSHtoHTTPS converts an ssh-style remote into its corresponding https style remote.
// Based on regex, and may not match every possible style of remote, but tested on github and gitlab.
//    Example input: git@git.test:testorg/testrepo.git
//           output: https://git.test/testorg/testrepo.git
func remoteConvertSSHtoHTTPS(sshRemote string) (string, error) {
	if !isSSHRemote(sshRemote) {
		return "", fmt.Errorf("%s is not an ssh remote", sshRemote)
	}

	httpsRemote := sshRemoteRegexp.ReplaceAllString(sshRemote, "https://$1/$2")
	if !isHTTPSRemote(httpsRemote) {
		return "", fmt.Errorf("failed to convert ssh remote %s to https remote: got %s, which is not valid", sshRemote, httpsRemote)
	}
	return httpsRemote, nil
}

// getRemoteName queries the repo at VcsChecker.path for all remotes, and then
// searches for the remote that matches the provided url, returning an error if
// no remote url is found.
// It will also attempt to match against different protocols - if an https protocol is
// specified, if it can't find an exact match, it will look for an ssh-style match (and vice-versa)
func (v *VcsChecker) getRemoteName(url string) (name string, err error) {
	remotes, err := v.repo.Remotes()
	if err != nil {
		return "", errors.Wrap(err, "failed to list remotes")
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("no remotes found for repo at path %q", v.path)
	}

	var exactMatchRemoteName string
	for _, remote := range remotes {
		remoteConfig := remote.Config()
		if remoteConfig == nil {
			continue
		}
		if len(remoteConfig.Fetch) == 0 {
			// Must be able to fetch from the remote. This could happen if a remote is set up as a push mirror.
			continue
		}
		for _, remoteUrl := range remoteConfig.URLs {
			if remoteUrl == url {
				if exactMatchRemoteName != "" {
					// NOTE(izaak): I can't think of a dev setup where you'd get multiple remotes with the same url.
					// If it does though, I think it's likely that any remote will work for us for diffing purposes,
					// wo we'll warn and continue.
					v.log.Warn("Found multiple remotes with the target url. Will choose remote-1.", "url", url, "remote-1", exactMatchRemoteName, "remote-2", remoteConfig.Name)
				} else {
					exactMatchRemoteName = remoteConfig.Name
				}
			}
		}
	}

	if exactMatchRemoteName != "" {
		return exactMatchRemoteName, nil
	}

	// Try to find an alternate match
	var alternateProtocolMatch string

	for _, remote := range remotes {
		remoteConfig := remote.Config()
		if remoteConfig == nil {
			continue
		}
		if len(remoteConfig.Fetch) == 0 {
			// Must be able to fetch from the remote. This could happen if a remote is set up as a push mirror.
			continue
		}
		for _, remoteUrl := range remoteConfig.URLs {
			var convertedUrl string
			if isHTTPSRemote(url) && isSSHRemote(remoteUrl) {
				convertedUrl, err = remoteConvertHTTPStoSSH(url)
				if err != nil {
					v.log.Debug("failed to convert https remote to ssh remote", "httpsRemote", url, "error", err)
				}
			}
			if isSSHRemote(url) && isHTTPSRemote(remoteUrl) {
				convertedUrl, err = remoteConvertSSHtoHTTPS(url)
				if err != nil {
					v.log.Debug("failed to convert ssh remote to https remote", "sshRemote", url, "error", err)
				}
			}

			if convertedUrl != "" && convertedUrl == remoteUrl {
				if alternateProtocolMatch != "" {
					// NOTE(izaak): I can't think of a dev setup where you'd get multiple remotes with the same url.
					// If it does though, I think it's likely that any remote will work for us for diffing purposes,
					// wo we'll warn and continue.
					v.log.Warn("Found multiple remotes that match the target URL, albeit with a different protocol. Will choose remote-1.", "url", url, "remote-1", exactMatchRemoteName, "remote-2", remoteConfig.Name)
				} else {
					alternateProtocolMatch = remoteConfig.Name
				}
			}
		}
	}

	if alternateProtocolMatch != "" {
		v.log.Debug("found remote with an alternate protocol that matches remote url", "url", url, "matching remote name", alternateProtocolMatch)
		return alternateProtocolMatch, nil
	}

	return "", fmt.Errorf("no remote with url matching %s found", url)

}

// remoteHasDiff compares the local repo to the specified branch on the configured remote
func (v *VcsChecker) remoteHasDiff(remoteBranch string) (bool, error) {
	diff, err := v.runVcsGitCommand("diff", fmt.Sprintf("%s/%s", v.remoteName, remoteBranch))
	if err != nil {
		return false, errors.Wrapf(err, "failed to diff against remote %q on branch %q", v.remoteName, remoteBranch)
	}
	if len(diff) > 0 {
		return true, nil
	}
	return false, nil
}

// remoteFileHasDiff compares the named file in the local repo to the
// specified remote and branch
func (v *VcsChecker) remoteFileHasDiff(remoteName string, remoteBranch string, filename string) (bool, error) {
	_, err := v.runVcsGitCommand("diff", "--quiet", fmt.Sprintf("%s/%s", remoteName, remoteBranch), "--", filename)
	if err == nil {
		return false, nil
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false, errors.Wrapf(err, "failed to diff against remote %q on branch %q", remoteName, remoteBranch)
	}
	if exitErr.ExitCode() != 0 {
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
