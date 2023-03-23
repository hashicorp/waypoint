// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestGitProjectSource(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected *pb.Job_Git
	}{
		{
			"minimum",
			`
url = "foo"
`,
			&pb.Job_Git{
				Url: "foo",
			},
		},

		{
			"basic auth",
			`
url = "foo"
username = "alice"
password = "giraffe"
`,
			&pb.Job_Git{
				Url: "foo",
				Auth: &pb.Job_Git_Basic_{
					Basic: &pb.Job_Git_Basic{
						Username: "alice",
						Password: "giraffe",
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			// Parse the input
			f, diag := hclsyntax.ParseConfig([]byte(tt.Input), "<test>", hcl.Pos{})
			require.False(diag.HasErrors())

			// Get the project source value
			var s GitSource
			result, err := s.ProjectSource(f.Body, &hcl.EvalContext{})
			require.NoError(err)
			actual := result.Source.(*pb.Job_DataSource_Git).Git
			require.Equal(actual, tt.Expected)
		})
	}
}

func TestGitSourceOverride(t *testing.T) {
	cases := []struct {
		Name     string
		Input    *pb.Job_DataSource
		M        map[string]string
		Expected *pb.Job_DataSource
		Error    string
	}{
		{
			"nothing",
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "foo",
					},
				},
			},
			nil,
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "foo",
					},
				},
			},
			"",
		},

		{
			"ref",
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "foo",
					},
				},
			},
			map[string]string{"ref": "bar"},
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "foo",
						Ref: "bar",
					},
				},
			},
			"",
		},

		{
			"ref with auth",
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "foo",
						Auth: &pb.Job_Git_Basic_{
							Basic: &pb.Job_Git_Basic{
								Username: "foo",
								Password: "bar",
							},
						},
					},
				},
			},
			map[string]string{"ref": "bar"},
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "foo",
						Ref: "bar",
						Auth: &pb.Job_Git_Basic_{
							Basic: &pb.Job_Git_Basic{
								Username: "foo",
								Password: "bar",
							},
						},
					},
				},
			},
			"",
		},

		{
			"invalid",
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "foo",
					},
				},
			},
			map[string]string{"other": "bar"},
			nil,
			"other",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			var s GitSource
			err := s.Override(tt.Input, tt.M)
			if tt.Error != "" {
				require.Error(err)
				require.Contains(err.Error(), tt.Error)
				return
			}

			require.NoError(err)
			require.Equal(tt.Expected, tt.Input)
		})
	}
}

func TestGitSourceGet(t *testing.T) {
	t.Run("basic clone", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, refRaw, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: testGitFixture(t, "git-noop"),
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "waypoint.hcl"))
		require.NoError(err)

		// Verify ref
		ref := refRaw.Ref.(*pb.Job_DataSource_Ref_Git).Git
		require.Equal("b6bf15100c570f2be6a231a095d395ed16dfed81", ref.Commit)
		require.Equal("Fixes\n", ref.CommitMessage)

		ts := ref.Timestamp.AsTime()
		require.False(ts.IsZero())
	})

	t.Run("branch ref", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, _, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: testGitFixture(t, "git-refs"),
						Ref: "branch",
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "waypoint.hcl"))
		require.NoError(err)
		_, err = os.Stat(filepath.Join(dir, "branchfile"))
		require.NoError(err)
	})

	t.Run("commit", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, _, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: testGitFixture(t, "git-refs"),
						Ref: "29758b9",
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "waypoint.hcl"))
		require.NoError(err)
		_, err = os.Stat(filepath.Join(dir, "two"))
		require.Error(err)
	})

	t.Run("submodule: clone disabled (default)", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, _, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: testGitFixture(t, "git-submodule"),
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "hello.txt"))
		require.NoError(err)
		_, err = os.Stat(filepath.Join(dir, "examples", "README.md"))
		require.Error(err)
	})

	t.Run("submodule: recursion on clone", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, _, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:               testGitFixture(t, "git-submodule"),
						RecurseSubmodules: 10,
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "hello.txt"))
		require.NoError(err)
		_, err = os.Stat(filepath.Join(dir, "examples", "README.md"))
		require.NoError(err)
	})

	t.Run("submodule: ref with no submodules", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, _, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:               testGitFixture(t, "git-submodule"),
						Ref:               "758c263",
						RecurseSubmodules: 10,
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "hello.txt"))
		require.NoError(err)
		_, err = os.Stat(filepath.Join(dir, "examples", "README.md"))
		require.Error(err)
	})

	t.Run("submodule: removed HEAD", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, _, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:               testGitFixture(t, "git-submodule-rm"),
						RecurseSubmodules: 10,
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "hello.txt"))
		require.NoError(err)
		_, err = os.Stat(filepath.Join(dir, "examples", "README.md"))
		require.Error(err)
	})

	t.Run("submodule: ref with submodule", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		dir, _, closer, err := s.Get(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:               testGitFixture(t, "git-submodule-rm"),
						Ref:               "27e97ef4f312fe37588f84209e4d056825dee614",
						RecurseSubmodules: 10,
					},
				},
			},
			"",
		)
		require.NoError(err)
		if closer != nil {
			defer closer()
		}

		// Verify files
		_, err = os.Stat(filepath.Join(dir, "hello.txt"))
		require.NoError(err)
		_, err = os.Stat(filepath.Join(dir, "examples", "README.md"))
		require.NoError(err)
	})
}

func TestGitSourceChanges(t *testing.T) {
	var latestRef *pb.Job_DataSource_Ref

	// NOTE(mitchellh): Most of the tests in this test use the Waypoint
	// GitHub repo and require an internet connection. There is some brittleness
	// to that if the internet is down, GitHub is down, etc. but we feel this
	// is a fairly stable dependency for tests and the benefits of testing
	// this functionality outweigh the risks.

	t.Run("nil current ref", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/hashicorp/waypoint.git",
					},
				},
			},
			nil,
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)

		latestRef = newRef
	})

	// Note this test depends on the 'nil current ref' test to run first.
	// t.Run tests are run sequentially which makes this work.
	t.Run("with latest ref", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/hashicorp/waypoint.git",
					},
				},
			},
			latestRef,
			"",
		)
		require.NoError(err)
		require.Nil(newRef)
		require.False(ignore)
	})

	t.Run("with old ref", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/hashicorp/waypoint.git",
					},
				},
			},
			&pb.Job_DataSource_Ref{
				Ref: &pb.Job_DataSource_Ref_Git{
					Git: &pb.Job_Git_Ref{
						Commit: "old",
					},
				},
			},
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)
	})

	// Test a specific tag ref, we expect our public Waypoint repo tags
	// to never change for the purpose of this test.
	t.Run("tag ref shorthand", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/hashicorp/waypoint.git",
						Ref: "v0.1.0",
					},
				},
			},
			nil,
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)

		hash := newRef.Ref.(*pb.Job_DataSource_Ref_Git).Git.Commit
		require.Equal(hash, "66d19f02c5da9e628998d688cbc0d1755eeabf62")
	})

	t.Run("tag ref full", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/hashicorp/waypoint.git",
						Ref: "refs/tags/v0.1.0",
					},
				},
			},
			nil,
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)

		hash := newRef.Ref.(*pb.Job_DataSource_Ref_Git).Git.Commit
		require.Equal(hash, "66d19f02c5da9e628998d688cbc0d1755eeabf62")
	})

	// This assumes release/0.1.0 won't change again, which it probably won't.
	t.Run("branch ref shorthand", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/hashicorp/waypoint.git",
						Ref: "release/0.1.0",
					},
				},
			},
			nil,
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)

		hash := newRef.Ref.(*pb.Job_DataSource_Ref_Git).Git.Commit
		require.Equal(hash, "a71a259607c26e93037aee9f2496a1da83dea6f2")
	})

	t.Run("branch ref full", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/hashicorp/waypoint.git",
						Ref: "refs/heads/release/0.1.0",
					},
				},
			},
			nil,
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)

		hash := newRef.Ref.(*pb.Job_DataSource_Ref_Git).Git.Commit
		require.Equal(hash, "a71a259607c26e93037aee9f2496a1da83dea6f2")
	})

	t.Run("no changes in specific path without ignore setting", func(t *testing.T) {
		require := require.New(t)

		hclog.L().SetLevel(hclog.Trace)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:  "https://github.com/hashicorp/waypoint.git",
						Ref:  "release/0.1.0",
						Path: "idontexist",
					},
				},
			},
			&pb.Job_DataSource_Ref{
				Ref: &pb.Job_DataSource_Ref_Git{
					Git: &pb.Job_Git_Ref{
						Commit: "38a28ec5af18265189fc6d55fd5970fd5e48544d",
					},
				},
			},
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)
	})

	// Detect no changes in a specific path
	t.Run("no changes in specific path", func(t *testing.T) {
		require := require.New(t)

		hclog.L().SetLevel(hclog.Trace)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:                      "https://github.com/hashicorp/waypoint.git",
						Ref:                      "release/0.1.0",
						Path:                     "idontexist",
						IgnoreChangesOutsidePath: true,
					},
				},
			},
			&pb.Job_DataSource_Ref{
				Ref: &pb.Job_DataSource_Ref_Git{
					Git: &pb.Job_Git_Ref{
						Commit: "38a28ec5af18265189fc6d55fd5970fd5e48544d",
					},
				},
			},
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.True(ignore)
	})

	t.Run("changes in specific path", func(t *testing.T) {
		require := require.New(t)

		hclog.L().SetLevel(hclog.Trace)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:                      "https://github.com/hashicorp/waypoint.git",
						Ref:                      "release/0.1.0",
						Path:                     "./internal",
						IgnoreChangesOutsidePath: true,
					},
				},
			},
			&pb.Job_DataSource_Ref{
				Ref: &pb.Job_DataSource_Ref_Git{
					Git: &pb.Job_Git_Ref{
						Commit: "38a28ec5af18265189fc6d55fd5970fd5e48544d",
					},
				},
			},
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)
	})

	// If there is no current commit, specific path always returns changes
	t.Run("specific path with no current commit", func(t *testing.T) {
		require := require.New(t)

		hclog.L().SetLevel(hclog.Trace)

		var s GitSource
		newRef, ignore, err := s.Changes(
			context.Background(),
			hclog.L(),
			terminal.ConsoleUI(context.Background()),
			&pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url:                      "https://github.com/hashicorp/waypoint.git",
						Ref:                      "release/0.1.0",
						Path:                     "idontexist",
						IgnoreChangesOutsidePath: true,
					},
				},
			},
			nil,
			"",
		)
		require.NoError(err)
		require.NotNil(newRef)
		require.False(ignore)
	})
}

// testGitFixture MUST be called before TestRunner since TestRunner
// changes our working directory.
func testGitFixture(t *testing.T, n string) string {
	t.Helper()

	// Get our full path
	wd, err := os.Getwd()
	require.NoError(t, err)
	wd, err = filepath.Abs(wd)
	require.NoError(t, err)
	path := filepath.Join(wd, "testdata", n)

	// Look for a DOTgit
	original := filepath.Join(path, "DOTgit")
	_, err = os.Stat(original)
	require.NoError(t, err)

	// Rename it
	newPath := filepath.Join(path, ".git")
	require.NoError(t, os.Rename(original, newPath))
	t.Cleanup(func() { os.Rename(newPath, original) })

	return path
}
