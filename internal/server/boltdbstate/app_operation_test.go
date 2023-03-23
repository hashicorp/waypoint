// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func TestAppOperation(t *testing.T) {
	ctx := context.Background()
	op := &appOperation{
		Struct: (*pb.Build)(nil),
		Bucket: buildOp.Bucket,
	}

	op.Test(t)

	t.Run("basic put and get and delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "A",
		})))

		// Read it back
		raw, err := op.Get(s, appOpById("A"))
		require.NoError(err)
		require.NotNil(raw)

		b, ok := raw.(*pb.Build)
		require.True(ok)
		require.NotNil(b.Application)
		require.Equal("A", b.Id)
		require.Equal(uint64(1), b.Sequence)

		// Create another, try to change the sequence number
		require.NoError(op.Put(s, true, serverptypes.TestValidBuild(t, &pb.Build{
			Id:       "A",
			Sequence: 2,
		})))

		// Read it back
		raw, err = op.Get(s, appOpById("A"))
		require.NoError(err)
		require.NotNil(raw)

		b, ok = raw.(*pb.Build)
		require.True(ok)
		require.NotNil(b.Application)
		require.Equal("A", b.Id)
		require.Equal(uint64(1), b.Sequence)

		// Get it by sequence
		raw, err = op.Get(s, &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Sequence{
				Sequence: &pb.Ref_OperationSeq{
					Application: b.Application,
					Number:      b.Sequence,
				},
			},
		})
		require.NoError(err)
		require.NotNil(raw)

		b, ok = raw.(*pb.Build)
		require.True(ok)
		require.NotNil(b.Application)
		require.Equal("A", b.Id)
		require.Equal(uint64(1), b.Sequence)

		// Delete it by ID
		err = op.Delete(s, b)
		require.Nil(err)
	})

	t.Run("get with data source ref", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a job with a data source ref
		ref := &pb.Job_DataSource_Ref{
			Ref: &pb.Job_DataSource_Ref_Git{
				Git: &pb.Job_Git_Ref{
					Commit: "hello!",
				},
			},
		}
		require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
			Id:            "jobA",
			DataSourceRef: ref,
		})))

		// Create a build
		require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id:    "A",
			JobId: "jobA",
		})))

		// Read it back
		raw, err := op.Get(s, appOpById("A"))
		require.NoError(err)
		require.NotNil(raw)

		b, ok := raw.(*pb.Build)
		require.True(ok)
		require.NotNil(b.Application)
		require.Equal("A", b.Id)
		require.NotNil(b.Preload)
		require.NotNil(b.Preload.JobDataSourceRef)

		actualRef := b.Preload.JobDataSourceRef
		gitRef, ok := actualRef.Ref.(*pb.Job_DataSource_Ref_Git)
		require.True(ok)
		require.Equal("hello!", gitRef.Git.Commit)
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

		// Since we search by latest sequence later, we don't want our latest
		// to be sequence #1. This fixes a flaky test that would happen in
		// that case cause we'd search for sequence 0 which doesn't exist.
		if times[0] == latest {
			times[0], times[1] = times[1], times[0]
		}

		// Create a build for each time
		for _, timeVal := range times {
			pt := timestamppb.New(timeVal)

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

		// Try getting the latest prior to the latest using a filter
		raw, err = op.LatestFilter(s, ref, nil, func(raw interface{}) (bool, error) {
			cand := raw.(*pb.Build)
			return cand.Sequence == b.Sequence-1, nil
		})
		require.NoError(err)
		b2 := raw.(*pb.Build)
		require.Equal(b2.Sequence, b.Sequence-1)

		// Try listing
		builds, err := op.List(s, &serverstate.ListOperationOptions{
			Application: ref,
		})
		require.NoError(err)
		require.Len(builds, len(times))

		// Lists should be in descending order by completion time
		var lastTime time.Time
		for _, raw := range builds {
			build := raw.(*pb.Build)
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

		_, err := op.List(s, &serverstate.ListOperationOptions{})
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
			pt := timestamppb.New(ts)

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
			pt := timestamppb.New(ts)

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
			pt := timestamppb.New(ts)

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
		results, err := op.List(s, &serverstate.ListOperationOptions{
			Application: ref,
			Status: []*pb.StatusFilter{
				{
					Filters: []*pb.StatusFilter_Filter{
						{
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
		results, err := op.List(s, &serverstate.ListOperationOptions{
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
		results, err := op.List(s, &serverstate.ListOperationOptions{
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

	t.Run("list with physical state filter on unsupported struct", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		{
			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id: "A",
			})))
		}

		// List with a filter
		build := serverptypes.TestValidBuild(t, nil)
		results, err := op.List(s, &serverstate.ListOperationOptions{
			Application:   build.Application,
			PhysicalState: pb.Operation_CREATED,
		})
		require.NoError(err)
		require.Len(results, 1)
	})

	t.Run("list with memwatch", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ref := &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		}

		{
			ts := time.Now().Add(5 * time.Hour)
			pt := timestamppb.New(ts)

			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id:          "A",
				Application: ref,
				Status: &pb.Status{
					State:     pb.Status_SUCCESS,
					StartTime: pt,
				},
			})))
		}
		{
			ts := time.Now().Add(6 * time.Hour)
			pt := timestamppb.New(ts)

			require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
				Id:          "B",
				Application: ref,
				Status: &pb.Status{
					State:     pb.Status_SUCCESS,
					StartTime: pt,
				},
			})))
		}

		ws := memdb.NewWatchSet()

		// make sure the watchset was populated
		require.Equal(0, len(ws))

		// List with a filter
		results, err := op.List(s, &serverstate.ListOperationOptions{
			Application: ref,
			WatchSet:    ws,
		})
		require.NoError(err)
		require.Len(results, 2)

		// make sure the watchset was populated
		require.Equal(1, len(ws))

		// Now add new item to fire the watch channel
		ts := time.Now().Add(8 * time.Hour)
		pt := timestamppb.New(ts)

		require.NoError(op.Put(s, false, serverptypes.TestValidBuild(t, &pb.Build{
			Id:          "D",
			Application: ref,
			Status: &pb.Status{
				State:     pb.Status_SUCCESS,
				StartTime: pt,
			},
		})))

		// Observe that the watch fires
		require.False(ws.Watch(time.After(1 * time.Second)))
	})

	t.Run("attempt deletion of non-existent operation id", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Attempt to delete an operation that doesn't exist
		err := op.Delete(s, &pb.Build{
			Id: "id",
			Application: &pb.Ref_Application{
				Application: "app123",
				Project:     "project456",
			},
		})
		require.Error(err)
	})
}

func TestAppOperation_deploy(t *testing.T) {
	op := &appOperation{
		Struct: (*pb.Deployment)(nil),
		Bucket: buildOp.Bucket,
	}

	op.Test(t)

	t.Run("basic put and get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create with preload set
		require.NoError(op.Put(s, false, serverptypes.TestValidDeployment(t, &pb.Deployment{
			Id: "A",
			Preload: &pb.Deployment_Preload{
				Build: serverptypes.TestValidBuild(t, nil),
			},
		})))

		// Read it back
		raw, err := op.Get(s, appOpById("A"))
		require.NoError(err)
		require.NotNil(raw)

		b, ok := raw.(*pb.Deployment)
		require.True(ok)
		require.Equal("A", b.Id)
		require.Equal(uint64(1), b.Sequence)
		require.NotEmpty(b.Generation)
		require.Equal(b.Sequence, b.Generation.InitialSequence)
		require.NotNil(b.Preload)
		require.Nil(b.Preload.Build)
		aSeq := b.Sequence

		// Inserting it again will preserve the sequence number
		b.Id = "C"
		require.NoError(op.Put(s, false, b))

		raw, err = op.Get(s, appOpById("C"))
		require.NoError(err)
		require.NotNil(raw)
		c, ok := raw.(*pb.Deployment)
		require.True(ok)
		require.NotEmpty(b.Generation)
		require.Equal(b.Generation.Id, c.Generation.Id)
		require.Equal(aSeq, c.Generation.InitialSequence)

		// Insert it again with a different generation
		b.Id = "D"
		b.Generation.Id = "other"
		require.NoError(op.Put(s, false, b))

		raw, err = op.Get(s, appOpById("D"))
		require.NoError(err)
		require.NotNil(raw)
		d, ok := raw.(*pb.Deployment)
		require.True(ok)
		require.Equal(uint64(3), d.Sequence)
		require.NotEmpty(b.Generation)
		require.Equal("other", d.Generation.Id)
		require.Equal(uint64(3), d.Generation.InitialSequence)

		// Cannot update the generation sequence
		d.Generation.InitialSequence = 42
		require.NoError(op.Put(s, true, d))

		raw, err = op.Get(s, appOpById("D"))
		require.NoError(err)
		require.NotNil(raw)
		d, ok = raw.(*pb.Deployment)
		require.True(ok)
		require.Equal(uint64(3), d.Sequence)
		require.NotEmpty(b.Generation)
		require.Equal("other", d.Generation.Id)
		require.Equal(uint64(3), d.Generation.InitialSequence)
	})

	t.Run("does not change generation if set", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Nothing special about this except it is properly formatted.
		expectedId := "f0866585-2aab-498b-8842-27706865acda"

		// Create with preload set
		require.NoError(op.Put(s, false, serverptypes.TestValidDeployment(t, &pb.Deployment{
			Id:         "A",
			Generation: &pb.Generation{Id: expectedId},
			Preload: &pb.Deployment_Preload{
				Build: serverptypes.TestValidBuild(t, nil),
			},
		})))

		// Read it back
		raw, err := op.Get(s, appOpById("A"))
		require.NoError(err)
		require.NotNil(raw)

		b, ok := raw.(*pb.Deployment)
		require.True(ok)
		require.Equal("A", b.Id)
		require.Equal(expectedId, b.Generation.Id)
		require.NotNil(b.Preload)
		require.Nil(b.Preload.Build)
	})

	t.Run("list with physical state filter supported", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		{
			require.NoError(op.Put(s, false, serverptypes.TestValidDeployment(t, &pb.Deployment{
				Id:    "A",
				State: pb.Operation_CREATED,
			})))
		}
		{
			require.NoError(op.Put(s, false, serverptypes.TestValidDeployment(t, &pb.Deployment{
				Id:    "B",
				State: pb.Operation_PENDING,
			})))
		}

		// List with a filter
		deploy := serverptypes.TestValidDeployment(t, nil)
		results, err := op.List(s, &serverstate.ListOperationOptions{
			Application:   deploy.Application,
			PhysicalState: pb.Operation_CREATED,
		})
		require.NoError(err)
		require.Len(results, 1)

		var ids []string
		for _, result := range results {
			ids = append(ids, result.(*pb.Deployment).Id)
		}
		sort.Strings(ids)
		require.Equal("A", ids[0])
	})
}

func TestAppOperation_workspaceResource(t *testing.T) {
	op := &appOperation{
		Struct: (*pb.Build)(nil),
		Bucket: buildOp.Bucket,
	}

	actual := op.workspaceResource()
	require.Equal(t, "Build", actual)
}

func appOpById(id string) *pb.Ref_Operation {
	return &pb.Ref_Operation{
		Target: &pb.Ref_Operation_Id{Id: id},
	}
}
