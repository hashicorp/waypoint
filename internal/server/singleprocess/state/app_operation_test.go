package state

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestAppOperation(t *testing.T) {
	op := &appOperation{
		Struct: (*pb.Build)(nil),
		Bucket: buildOp.Bucket,
	}

	op.Test(t)

	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(op.Put(s, false, &pb.Build{
			Id: "A",
			Application: &pb.Ref_Application{
				Application: "a_test",
				Project:     "p_test",
			},
		}))

		// Read it back
		raw, err := op.Get(s, "A")
		require.NoError(err)
		require.NotNil(raw)

		b, ok := raw.(*pb.Build)
		require.True(ok)
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
		for _, t := range times {
			pt, err := ptypes.TimestampProto(t)
			require.NoError(err)

			require.NoError(op.Put(s, false, &pb.Build{
				Id: strconv.FormatInt(t.Unix(), 10),
				Application: &pb.Ref_Application{
					Application: "a_test",
					Project:     "p_test",
				},

				Status: &pb.Status{
					State:        pb.Status_SUCCESS,
					StartTime:    pt,
					CompleteTime: pt,
				},
			}))
		}

		ref := &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		}

		// Get the latest
		raw, err := op.Latest(s, ref)
		require.NoError(err)
		b := raw.(*pb.Build)
		require.Equal(strconv.FormatInt(latest.Unix(), 10), b.Id)

		// Try listing
		builds, err := op.List(s, ref)
		require.NoError(err)
		require.Len(builds, len(times))

		// Lists should be in descending order by completion time
		var lastTime time.Time
		for _, raw := range builds {
			build := raw.(*pb.Build)
			timeVal, err := ptypes.Timestamp(build.Status.CompleteTime)
			require.NoError(err)

			if !lastTime.IsZero() && timeVal.After(lastTime) {
				t.Fatal("timestamp should be descending")
			}

			lastTime = timeVal
		}
	})

	t.Run("latest nil if none are completed", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ref := &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		}

		ts := time.Now().Add(5 * time.Hour)
		pt, err := ptypes.TimestampProto(ts)
		require.NoError(err)

		require.NoError(op.Put(s, false, &pb.Build{
			Id:          strconv.FormatInt(ts.Unix(), 10),
			Application: ref,
			Status: &pb.Status{
				State:     pb.Status_RUNNING,
				StartTime: pt,
			},
		}))

		// Get the latest
		b, err := op.Latest(s, ref)
		require.NoError(err)
		require.Nil(b)
	})
}
