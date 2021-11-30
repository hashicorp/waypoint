package gitdirty

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
	githubStyleSshRemoteRegexp  *regexp.Regexp
	githubStyleHttpRemoteRegexp *regexp.Regexp
)

func init() {
	// Regex matching http/https remotes, tokenizing the unique components for replacement.
	// Works for github, gitlab, sourcehut, and other remotes using this style.
	githubStyleSshRemoteRegexp = regexp.MustCompile(`git@(.*?\..*?):(.*)`)            // Example: git@git.test:testorg/testrepo.git
	githubStyleHttpRemoteRegexp = regexp.MustCompile(`http[s]?:\/\/(.*?\..*?)\/(.*)`) // Example: https://git.test/testorg/testrepo.git
}

// RepoIsDirty looks for unstaged, staged, and committed (but not pushed)
// changes on the local GitDirty.path repo not on the specified remote
// url and branch.
// CAVEAT: This does not fetch any remotes, and therefore will not detect if
// the local copy is behind the remote,
func RepoIsDirty(log hclog.Logger, repoPath string, remoteUrl string, remoteBranch string) (bool, error) {
	return FileIsDirty(log, repoPath, remoteUrl, remoteBranch, "")
}

// FileIsDirty checks only the specified file for unstaged, staged, and committed
// (but not pushed) changes on the local GitDirty.path repo not on the specified remote
// url and branch. If filePath is empty, this will check the entire repo.
// CAVEAT: This does not fetch any remotes, and therefore will not detect if
// the local copy is behind the remote,
func FileIsDirty(log hclog.Logger, repoPath string, remoteUrl string, remoteBranch string, filePath string) (bool, error) {
	remoteName, err := getRemoteName(log, repoPath, remoteUrl)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to get remote name for url %s", remoteUrl)
	}

	hasBranch, err := remoteHasTrackingBranch(log, repoPath, remoteName, remoteBranch)
	if err != nil {
		return false, errors.Wrapf(err, "failed to determine if remote %s has branch %s", remoteName, remoteBranch)
	}
	if !hasBranch {
		return false, fmt.Errorf(
			"remote %s does not have specified branch %s. To fix this, try running `git fetch %s`",
			remoteName, remoteBranch, remoteName,
		)
	}

	diff, err := remoteHasDiff(log, repoPath, remoteName, remoteBranch, filePath)
	if err != nil {
		return false, err
	}
	return diff, nil
}

// remoteHasTrackingBranch checks to see if the configured remote
func remoteHasTrackingBranch(log hclog.Logger, repoPath string, remoteName string, branch string) (bool, error) {
	remoteBranchOutput, err := runGitCommand(log, repoPath, "branch", "-r")
	if err != nil {
		return false, errors.Wrapf(err, "failed to list branches for repo at path %s", repoPath)
	}
	trackingBranch := fmt.Sprintf("%s/%s", remoteName, branch)
	branches := strings.Split(remoteBranchOutput, "\n")
	for _, thisBranch := range branches {
		thisBranch = strings.TrimSpace(thisBranch)
		if thisBranch == trackingBranch {
			return true, nil
		}
	}
	return false, nil
}

func isSSHRemote(remote string) bool {
	return githubStyleSshRemoteRegexp.MatchString(remote)
}

func isHTTPSRemote(remote string) bool {
	return githubStyleHttpRemoteRegexp.MatchString(remote)
}

// remoteConvertHTTPStoSSH converts an https-style remote into its corresponding ssh style remote.
// Based on regex, and may not match every possible style of remote, but tested on github and gitlab.
//    Example input: https://git.test/testorg/testrepo.git
//           output: git@git.test:testorg/testrepo.git
func remoteConvertHTTPStoSSH(httpsRemote string) (string, error) {
	if !isHTTPSRemote(httpsRemote) {
		return "", fmt.Errorf("%s is not an https remote", httpsRemote)
	}

	sshRemote := githubStyleHttpRemoteRegexp.ReplaceAllString(httpsRemote, "git@$1:$2")
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

	httpsRemote := githubStyleSshRemoteRegexp.ReplaceAllString(sshRemote, "https://$1/$2")
	if !isHTTPSRemote(httpsRemote) {
		return "", fmt.Errorf("failed to convert ssh remote %s to https remote: got %s, which is not valid", sshRemote, httpsRemote)
	}
	return httpsRemote, nil
}

// getRemoteName queries the repo at GitDirty.path for all remotes, and then
// searches for the remote that matches the provided url, returning an error if
// no remote url is found.
// It will also attempt to match against different protocols - if an https protocol is
// specified, if it can't find an exact match, it will look for an ssh-style match (and vice-versa)
func getRemoteName(log hclog.Logger, repoPath string, remoteUrl string) (name string, err error) {
	repo, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to open git repo at path %s", repoPath)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", errors.Wrap(err, "failed to list remotes")
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("no remotes found for repo at path %q", repoPath)
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
		for _, thisRemoteUrl := range remoteConfig.URLs {
			if thisRemoteUrl == remoteUrl {
				if exactMatchRemoteName != "" {
					// NOTE(izaak): I can't think of a dev setup where you'd get multiple remotes with the same url.
					// If it does though, I think it's likely that any remote will work for us for diffing purposes,
					// wo we'll warn and continue.
					log.Warn("Found multiple remotes with the target url. Will choose remote-1.", "url", thisRemoteUrl, "remote-1", exactMatchRemoteName, "remote-2", remoteConfig.Name)
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
	var alternateProtocolRemoteName string

	for _, remote := range remotes {
		remoteConfig := remote.Config()
		if remoteConfig == nil {
			continue
		}
		if len(remoteConfig.Fetch) == 0 {
			// Must be able to fetch from the remote. This could happen if a remote is set up as a push mirror.
			continue
		}
		for _, thisRemoteUrl := range remoteConfig.URLs {
			var convertedUrl string
			if isHTTPSRemote(remoteUrl) && isSSHRemote(thisRemoteUrl) {
				convertedUrl, err = remoteConvertHTTPStoSSH(remoteUrl)
				if err != nil {
					log.Debug("failed to convert https remote to ssh remote", "httpsRemote", remoteUrl, "error", err)
				}
			}
			if isSSHRemote(remoteUrl) && isHTTPSRemote(thisRemoteUrl) {
				convertedUrl, err = remoteConvertSSHtoHTTPS(remoteUrl)
				if err != nil {
					log.Debug("failed to convert ssh remote to https remote", "sshRemote", remoteUrl, "error", err)
				}
			}

			if convertedUrl != "" && convertedUrl == thisRemoteUrl {
				if alternateProtocolRemoteName != "" {
					// NOTE(izaak): I can't think of a dev setup where you'd get multiple remotes with the same url.
					// If it does though, I think it's likely that any remote will work for us for diffing purposes,
					// wo we'll warn and continue.
					log.Warn("Found multiple remotes that match the target URL, albeit with a different protocol. Will choose remote-1.", "url", remoteUrl, "remote-1", exactMatchRemoteName, "remote-2", remoteConfig.Name)
				} else {
					alternateProtocolRemoteName = remoteConfig.Name
				}
			}
		}
	}

	if alternateProtocolRemoteName != "" {
		log.Debug("found remote with an alternate protocol that matches remote url",
			"url", remoteUrl,
			"matching remote name", alternateProtocolRemoteName,
		)
		return alternateProtocolRemoteName, nil
	}

	return "", fmt.Errorf("no remote with url matching %s found", remoteUrl)
}

// remoteHasDiff compares the local repo to the specified branch on the configured remote.
// If filePath is not empty, it will check only the specified file path.
func remoteHasDiff(log hclog.Logger, repoPath string, remoteName string, remoteBranch string, filePath string) (bool, error) {
	args := []string{"diff", "--quiet", fmt.Sprintf("%s/%s", remoteName, remoteBranch)}
	if filePath != "" {
		args = append(args, "--", filePath)
	}
	_, err := runGitCommand(log, repoPath, args...)
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false, errors.Wrapf(err, "failed to diff against remote %q on branch %q", remoteName, remoteBranch)
	}
	if exitErr.ExitCode() != 0 {
		return true, nil
	}
	return false, nil
}

// runGitCommand executes a git command against the repo at the given path with
// the provided args, returning the standard output.
func runGitCommand(log hclog.Logger, path string, gitArgs ...string) (output string, err error) {
	args := append([]string{"-C", path}, gitArgs...)
	log.Debug(fmt.Sprintf("Running this command: git %s", strings.Join(args, " ")))
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
