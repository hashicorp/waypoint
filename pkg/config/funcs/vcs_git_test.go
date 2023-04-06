// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcs

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/waypoint/internal/pkg/copy"
)

var testHasGit bool

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
}

func TestVCSGit(t *testing.T) {
	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	cases := []struct {
		Name     string
		Fixture  string
		Subdir   string
		Func     func(*VCSGit, []cty.Value, cty.Type) (cty.Value, error)
		Args     []cty.Value
		Expected cty.Value
		Error    string
	}{
		{
			"hash: HEAD commit",
			"git-commits",
			"",
			(*VCSGit).refHashFunc,
			nil,
			cty.StringVal("380afd697abe993b89bfa08d8dd8724d6a513ba1"),
			"",
		},

		{
			"tag",
			"git-tag",
			"",
			(*VCSGit).refTagFunc,
			nil,
			cty.StringVal("hello"),
			"",
		},

		{
			"tag no tags",
			"git-commits",
			"",
			(*VCSGit).refTagFunc,
			nil,
			cty.StringVal(""),
			"",
		},

		{
			"remote doesn't exist",
			"git-commits",
			"",
			(*VCSGit).remoteUrlFunc,
			[]cty.Value{cty.StringVal("origin")},
			cty.UnknownVal(cty.String),
			"",
		},

		{
			"remote exists",
			"git-remote",
			"",
			(*VCSGit).remoteUrlFunc,
			[]cty.Value{cty.StringVal("origin")},
			cty.StringVal("https://github.com/hashicorp/example.git"),
			"",
		},

		{
			"refpretty with no changes",
			"git-commits",
			"",
			(*VCSGit).refPrettyFunc,
			nil,
			cty.StringVal("380afd697abe993b89bfa08d8dd8724d6a513ba1"),
			"",
		},

		{
			"refpretty with changes",
			"git-commits-changes",
			"",
			(*VCSGit).refPrettyFunc,
			nil,
			cty.StringVal("380afd697abe993b89bfa08d8dd8724d6a513ba1_CHANGES_*"),
			"",
		},

		{
			"refpretty with tags",
			"git-tag",
			"",
			(*VCSGit).refPrettyFunc,
			nil,
			cty.StringVal("hello"),
			"",
		},
	}

	for _, tt := range cases {
		for i := 0; i < 2; i++ {
			gogit := i == 1
			name := tt.Name
			if gogit {
				name += " (go-git)"
			}

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
				if tt.Subdir != "" {
					path = filepath.Join(path, tt.Subdir)
				}

				s := &VCSGit{Path: path, GoGitOnly: gogit}
				result, err := tt.Func(s, tt.Args, cty.String)
				if tt.Error != "" {
					require.Error(err)
					require.Contains(err.Error(), tt.Error)
					return
				}
				require.NoError(err)

				// If our expected value ends in _* then we do a prefix check
				// instead. This is so we can test the dynamic parts of the timestamp.
				if tt.Expected.Type() == cty.String && tt.Expected.IsKnown() {
					expected := tt.Expected.AsString()
					if strings.HasSuffix(expected, "_*") {
						require.True(strings.HasPrefix(
							result.AsString(), expected[:len(expected)-1]))
						return
					}
				}

				require.True(tt.Expected.RawEquals(result), result.GoString())
			})
		}
	}
}

func testGitFixture(t *testing.T, path string) {
	t.Helper()

	// Look for a DOTgit
	original := filepath.Join(path, "DOTgit")
	_, err := os.Stat(original)
	require.NoError(t, err)

	// Rename it
	newPath := filepath.Join(path, ".git")
	require.NoError(t, os.Rename(original, newPath))
	t.Cleanup(func() { os.Rename(newPath, original) })
}
