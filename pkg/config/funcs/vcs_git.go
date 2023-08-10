// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package funcs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func VCSGitFuncs(path string) map[string]function.Function {
	state := &VCSGit{Path: path}

	return map[string]function.Function{
		"gitrefpretty": state.RefPrettyFunc(),
		"gitrefhash":   state.RefHashFunc(),
		"gitreftag":    state.RefTagFunc(),
		"gitremoteurl": state.RemoteUrlFunc(),
	}
}

type VCSGit struct {
	// Path of the git repository. Parent directories will be searched for
	// a ".git" folder automatically.
	Path string

	// GoGitOnly forces go-git usage and disables all subprocessing to "git"
	// even if it exists.
	GoGitOnly bool

	initErr error
	repo    *git.Repository
}

// RefPrettyFunc returns a string format of the current Git ref. This function
// takes some liberties to humanize the output: it will use a tag if the
// ref matches a tag, it will append "+CHANGES" to the commit if there are
// uncommitted changed files, etc.
//
// You may use direct functions such as `gitrefhash` if you want the direct
// hash. Or `gitreftag` to get the current tag.
//
// waypoint:gitrefpretty
func (s *VCSGit) RefPrettyFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{},
		Type:   function.StaticReturnType(cty.String),
		Impl:   s.refPrettyFunc,
	})
}

func (s *VCSGit) refPrettyFunc(args []cty.Value, retType cty.Type) (cty.Value, error) {
	if err := s.init(); err != nil {
		return cty.UnknownVal(cty.String), err
	}

	ref, err := s.repo.Head()
	if err != nil {
		return cty.UnknownVal(cty.String), fmt.Errorf("error getting repo HEAD reference - this repo may have no commits: %w", err)
	}
	result := ref.Hash().String()

	// Get the tags
	iter, err := s.repo.Tags()
	if err != nil {
		return cty.UnknownVal(cty.String), fmt.Errorf("error getting repo tags: %w", err)
	}
	defer iter.Close()
	for {
		tagRef, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return cty.UnknownVal(cty.String), err
		}
		if tagRef.Hash() == ref.Hash() {
			result = tagRef.Name().Short()
			break
		}
	}

	// To determine if there are changes we subprocess because go-git's Status
	// function is really, really slow sadly. On the waypoint repo at the time
	// of this commit, go-git took 12s on my machine vs. 50ms for `git`.
	goGitChanges := s.GoGitOnly
	if !goGitChanges {
		cmd := exec.Command("git", "diff", "--quiet")
		cmd.Stdout = ioutil.Discard
		cmd.Stderr = ioutil.Discard
		cmd.Dir = s.Path
		err := cmd.Run()

		// If git isn't available, we fall back to using go-git. This can
		// take a very long time and that is sad but we want to give consistent
		// results from this func.
		if errors.Is(err, exec.ErrNotFound) {
			goGitChanges = true
		}

		if !goGitChanges && err != nil {
			exitError, ok := err.(*exec.ExitError)
			if !ok {
				return cty.UnknownVal(cty.String), fmt.Errorf("error executing git: %s", err)
			}

			if exitError.ExitCode() != 0 {
				result += fmt.Sprintf("_CHANGES_%d", time.Now().Unix())
			}
		}
	}

	// If we want to use go-git for change detection, then do that now.
	if goGitChanges {
		wt, err := s.repo.Worktree()
		if err != nil {
			return cty.UnknownVal(cty.String), err
		}

		st, err := wt.Status()
		if err != nil {
			return cty.UnknownVal(cty.String), err
		}

		if !st.IsClean() {
			result += fmt.Sprintf("_CHANGES_%d", time.Now().Unix())
		}
	}

	return cty.StringVal(result), nil
}

// RefHashFunc returns the full hash of the HEAD ref.
//
// waypoint:gitrefhash
func (s *VCSGit) RefHashFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{},
		Type:   function.StaticReturnType(cty.String),
		Impl:   s.refHashFunc,
	})
}

func (s *VCSGit) refHashFunc(args []cty.Value, retType cty.Type) (cty.Value, error) {
	if err := s.init(); err != nil {
		return cty.UnknownVal(cty.String), err
	}

	ref, err := s.repo.Head()
	if err != nil {
		return cty.UnknownVal(cty.String), fmt.Errorf("error getting repo HEAD reference - this repo may have no commits: %w", err)
	}

	return cty.StringVal(ref.Hash().String()), nil
}

// RefTagFunc returns the tag of the HEAD ref or empty if not tag is found.
//
// waypoint:gitreftag
func (s *VCSGit) RefTagFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{},
		Type:   function.StaticReturnType(cty.String),
		Impl:   s.refTagFunc,
	})
}

func (s *VCSGit) refTagFunc(args []cty.Value, retType cty.Type) (cty.Value, error) {
	if err := s.init(); err != nil {
		return cty.UnknownVal(cty.String), err
	}

	ref, err := s.repo.Head()
	if err != nil {
		return cty.UnknownVal(cty.String), fmt.Errorf("error getting repo HEAD reference - this repo may have no commits: %w", err)
	}

	// Get the tags
	iter, err := s.repo.Tags()
	if err != nil {
		return cty.UnknownVal(cty.String), fmt.Errorf("error getting repo tags: %w", err)
	}

	var tagRefStr string
	err = iter.ForEach(func(t *plumbing.Reference) error {
		if t.Hash() == ref.Hash() {
			tagRefStr = t.Name().Short()
		}
		return nil
	})
	if err != nil {
		return cty.UnknownVal(cty.String), err
	}

	if tagRefStr != "" {
		return cty.StringVal(tagRefStr), nil
	}

	return cty.StringVal(""), nil
}

// RemoteUrlFunc returns the URL for the matching remote or unknown
// if it can't be found.
//
// waypoint:gitremoteurl
func (s *VCSGit) RemoteUrlFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "name",
				Type: cty.String,
			},
		},
		Type: function.StaticReturnType(cty.String),
		Impl: s.remoteUrlFunc,
	})
}

func (s *VCSGit) remoteUrlFunc(args []cty.Value, retType cty.Type) (cty.Value, error) {
	if err := s.init(); err != nil {
		return cty.UnknownVal(cty.String), err
	}

	name := args[0].AsString()

	remote, err := s.repo.Remote(name)
	if err != nil {
		if err == git.ErrRemoteNotFound {
			err = nil
		}

		return cty.UnknownVal(cty.String), err
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return cty.UnknownVal(cty.String), nil
	}

	return cty.StringVal(urls[0]), nil
}

func (s *VCSGit) init() error {
	// If we initialized already return
	if s.initErr != nil {
		return s.initErr
	}
	if s.repo != nil {
		return nil
	}

	// Open the repo
	repo, err := git.PlainOpenWithOptions(s.Path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		s.initErr = err
		return err
	}
	s.repo = repo
	return nil
}
