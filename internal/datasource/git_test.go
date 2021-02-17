package datasource

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

		ts, err := ptypes.Timestamp(ref.Timestamp)
		require.NoError(err)
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
		newRef, err := s.Changes(
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
		)
		require.NoError(err)
		require.NotNil(newRef)

		latestRef = newRef
	})

	t.Run("with latest ref", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, err := s.Changes(
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
		)
		require.NoError(err)
		require.Nil(newRef)
	})

	t.Run("with old ref", func(t *testing.T) {
		require := require.New(t)

		var s GitSource
		newRef, err := s.Changes(
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
		)
		require.NoError(err)
		require.NotNil(newRef)
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
