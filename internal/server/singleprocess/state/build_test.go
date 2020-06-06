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

func TestBuild(t *testing.T) {
	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.BuildPut(false, &pb.Build{
			Id: "A",
			Application: &pb.Ref_Application{
				Application: "a_test",
				Project:     "p_test",
			},
		}))

		// Read it back
		b, err := s.BuildGet("A")
		require.NoError(err)
		require.NotNil(b.Application)
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

			require.NoError(s.BuildPut(false, &pb.Build{
				Id: strconv.FormatInt(t.Unix(), 10),
				Application: &pb.Ref_Application{
					Application: "a_test",
					Project:     "p_test",
				},

				Status: &pb.Status{
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
		b, err := s.BuildLatest(ref)
		require.NoError(err)
		require.Equal(strconv.FormatInt(latest.Unix(), 10), b.Id)

		// Try listing
		builds, err := s.BuildList(ref)
		require.NoError(err)
		require.Len(builds, len(times))

		// Lists should be in descending order by completion time
		var lastTime time.Time
		for _, build := range builds {
			timeVal, err := ptypes.Timestamp(build.Status.CompleteTime)
			require.NoError(err)

			if !lastTime.IsZero() && timeVal.After(lastTime) {
				t.Fatal("timestamp should be descending")
			}

			lastTime = timeVal
		}
	})
}
