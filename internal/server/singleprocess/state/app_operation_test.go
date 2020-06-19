package state

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
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
		require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "A",
		})))

		// Read it back
		raw, err := op.Get(s, "A")
		require.NoError(err)
		require.NotNil(raw)

		b, ok := raw.(*pb.Build)
		require.True(ok)
		require.NotNil(b.Application)
		require.Equal("A", b.Id)
	})

	t.Run("latest basic", func(t *testing.T) {
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
		for _, timeVal := range times {
			pt, err := ptypes.TimestampProto(timeVal)
			require.NoError(err)

			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: strconv.FormatInt(timeVal.Unix(), 10),
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
		raw, err := op.Latest(s, ref, nil)
		require.NoError(err)
		b := raw.(*pb.Build)
		require.Equal(strconv.FormatInt(latest.Unix(), 10), b.Id)

		// Try listing
		builds, err := op.List(s, &listOperationsOptions{
			Application: ref,
		})
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

	t.Run("returns error if none are completed", func(t *testing.T) {
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

		require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id:          strconv.FormatInt(ts.Unix(), 10),
			Application: ref,
			Status: &pb.Status{
				State:     pb.Status_RUNNING,
				StartTime: pt,
			},
		})))

		// Get the latest
		b, err := op.Latest(s, ref, nil)
		require.Error(err)
		require.Nil(b)
	})

	t.Run("list without application returns error", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		_, err := op.List(s, &listOperationsOptions{})
		require.Error(err)
	})

	t.Run("list with filter", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ref := &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		}

		{
			ts := time.Now().Add(5 * time.Hour)
			pt, err := ptypes.TimestampProto(ts)
			require.NoError(err)

			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id:          "A",
				Application: ref,
				Status: &pb.Status{
					State:     pb.Status_RUNNING,
					StartTime: pt,
				},
			})))
		}
		{
			ts := time.Now().Add(6 * time.Hour)
			pt, err := ptypes.TimestampProto(ts)
			require.NoError(err)

			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id:          "B",
				Application: ref,
				Status: &pb.Status{
					State:     pb.Status_ERROR,
					StartTime: pt,
				},
			})))
		}
		{
			ts := time.Now().Add(7 * time.Hour)
			pt, err := ptypes.TimestampProto(ts)
			require.NoError(err)

			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id:          "C",
				Application: ref,
				Status: &pb.Status{
					State:     pb.Status_ERROR,
					StartTime: pt,
				},
			})))
		}

		// List with a filter
		results, err := op.List(s, &listOperationsOptions{
			Application: ref,
			Status: []*pb.StatusFilter{
				&pb.StatusFilter{
					Filters: []*pb.StatusFilter_Filter{
						&pb.StatusFilter_Filter{
							Filter: &pb.StatusFilter_Filter_State{
								State: pb.Status_ERROR,
							},
						},
					},
				},
			},
		})
		require.NoError(err)
		require.Len(results, 2)
	})

	t.Run("list by workspace specified", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		{
			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: "A",
				Workspace: &pb.Ref_Workspace{
					Workspace: "WS_A",
				},
			})))
		}
		{
			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: "B",
				Workspace: &pb.Ref_Workspace{
					Workspace: "WS_B",
				},
			})))
		}
		{
			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: "C",
				Workspace: &pb.Ref_Workspace{
					Workspace: "WS_A",
				},
			})))
		}

		// List with a filter
		build := serverptypes.TestValidBuild(t, nil)
		results, err := op.List(s, &listOperationsOptions{
			Application: build.Application,
			Workspace:   &pb.Ref_Workspace{Workspace: "WS_A"},
		})
		require.NoError(err)
		require.Len(results, 2)

		var ids []string
		for _, result := range results {
			ids = append(ids, result.(*pb.Build).Id)
		}
		sort.Strings(ids)
		require.Equal("A", ids[0])
		require.Equal("C", ids[1])
	})

	t.Run("list by workspace unspecified", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		{
			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: "A",
				Workspace: &pb.Ref_Workspace{
					Workspace: "WS_A",
				},
			})))
		}
		{
			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: "B",
				Workspace: &pb.Ref_Workspace{
					Workspace: "WS_B",
				},
			})))
		}
		{
			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: "C",
				Workspace: &pb.Ref_Workspace{
					Workspace: "WS_A",
				},
			})))
		}

		// List with a filter
		build := serverptypes.TestValidBuild(t, nil)
		results, err := op.List(s, &listOperationsOptions{
			Application: build.Application,
		})
		require.NoError(err)
		require.Len(results, 3)

		var ids []string
		for _, result := range results {
			ids = append(ids, result.(*pb.Build).Id)
		}
		sort.Strings(ids)
		require.Equal("A", ids[0])
		require.Equal("B", ids[1])
		require.Equal("C", ids[2])
	})
}
