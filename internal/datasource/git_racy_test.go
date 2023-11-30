// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// These tests fail the race detector, and should eventually be fixed.
//go:build !race

package datasource

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

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
