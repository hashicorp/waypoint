// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func TestBuild(t *testing.T) {
	ctx := context.Background()
	buildOp.Test(t)

	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "A",
		})))

		// Read it back
		b, err := s.BuildGet(ctx, appOpById("A"))
		require.NoError(err)
		require.NotNil(b.Application)
		require.Equal("A", b.Id)
	})

	t.Run("latest", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Build a bunch of times
		var times []time.Time
		var latest time.Time
		for i := 0; i < 50; i++ {
			latest = time.Now().Add(time.Duration(10*i) * time.Hour)
			times = append(times, latest)
		}

		// Shuffle the times so that we insert them randomly. We do this
		// to test our index is really doing the right thing.
		rand.Shuffle(len(times), func(i, j int) { times[i], times[j] = times[j], times[i] })

		// Create a build for each time
		for _, ts := range times {
			pt := timestamppb.New(ts)

			require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: strconv.FormatInt(ts.Unix(), 10),
				Application: &pb.Ref_Application{
					Application: "a_test",
					Project:     "p_test",
				},

				Status: &pb.Status{
					State:        pb.Status_SUCCESS,
					StartTime:    pt,
					CompleteTime: pt,
				},
			})))
		}

		ref := &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		}

		// Get the latest
		b, err := s.BuildLatest(ctx, ref, nil)
		require.NoError(err)
		require.Equal(strconv.FormatInt(latest.Unix(), 10), b.Id)

		// Try listing
		builds, err := s.BuildList(ctx, ref)
		require.NoError(err)
		require.Len(builds, len(times))

		// Lists should be in descending order by completion time
		var lastTime time.Time
		for _, build := range builds {
			timeVal := build.Status.CompleteTime.AsTime()

			if !lastTime.IsZero() && timeVal.After(lastTime) {
				t.Fatal("timestamp should be descending")
			}

			lastTime = timeVal
		}
	})

	t.Run("returns error if none are completed", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ref := &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		}

		ts := time.Now().Add(5 * time.Hour)
		pt := timestamppb.New(ts)

		require.NoError(s.BuildPut(ctx, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id:          strconv.FormatInt(ts.Unix(), 10),
			Application: ref,
			Status: &pb.Status{
				State:     pb.Status_RUNNING,
				StartTime: pt,
			},
		})))

		// Get the latest
		b, err := s.BuildLatest(ctx, ref, nil)
		require.Error(err)
		require.Nil(b)
	})
}
