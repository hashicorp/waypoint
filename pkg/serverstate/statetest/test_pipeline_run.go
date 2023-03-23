// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statetest

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["pipeline_run"] = []testFunc{
		TestPipelineRun,
	}
}

func TestPipelineRun(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "project"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		// Set Pipeline
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(ctx, p)
		require.NoError(err)
		pipeline := &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: p.Id}}

		// Set Pipeline Run
		r := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		err = s.PipelineRunPut(ctx, r)
		require.NoError(err)

		// We manually add job ids to a PipelineRun over in RunPipeline, so this
		// replicates how job ids are added to a run messasge
		for i := 1; i < 4; i++ {
			is := strconv.Itoa(i)
			require.NoError(s.JobCreate(ctx, serverptypes.TestJobNew(t, &pb.Job{
				Id: is,
				Pipeline: &pb.Ref_PipelineStep{
					PipelineId:  p.Id,
					RunSequence: r.Sequence,
				},
			})))

			r.Jobs = append(r.Jobs, &pb.Ref_Job{Id: string(is)})
		}
		err = s.PipelineRunPut(ctx, r)
		require.NoError(err)

		// Get run by pipeline and sequence
		{
			resp, err := s.PipelineRunGet(ctx, pipeline, 1)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
			require.Equal(pipeline.Ref.(*pb.Ref_Pipeline_Id).Id, resp.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id)
			require.NotEmpty(resp.Jobs)
			require.Equal("1", resp.Jobs[0].Id)
			require.Equal("2", resp.Jobs[1].Id)
			require.Equal("3", resp.Jobs[2].Id)
		}

		// Get run by Id
		{
			resp, err := s.PipelineRunGetById(ctx, r.Id)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
		}

		// Set another pipeline run
		r2 := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		r2.Id = "ypr2"
		err = s.PipelineRunPut(ctx, r2)
		require.NoError(err)

		// Get run by pipeline and sequence, should auto increment
		{
			resp, err := s.PipelineRunGet(ctx, pipeline, 2)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r2.Id, resp.Id)
			require.Equal(r2.Sequence, resp.Sequence)
		}

		// Set another pipeline run
		r3 := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		r3.Id = "pr3"
		err = s.PipelineRunPut(ctx, r3)
		require.NoError(err)

		// Get latest run by pipeline ID
		{
			resp, err := s.PipelineRunGetLatest(ctx, pipeline.Ref.(*pb.Ref_Pipeline_Id).Id)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r3.Id, resp.Id)
			require.Equal(r3.Sequence, resp.Sequence)
		}

		// Get run by pipeline and sequence, should auto increment
		{
			resp, err := s.PipelineRunGet(ctx, pipeline, 3)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r3.Id, resp.Id)
			require.Equal(r3.Sequence, resp.Sequence)
		}
	})

	t.Run("Update existing run", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "project"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		// Set Pipeline
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(ctx, p)
		require.NoError(err)
		pipeline := &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: p.Id}}

		// Set Pipeline Run
		pr := &pb.PipelineRun{
			Id:       "test-pr",
			Pipeline: pipeline,
			Sequence: 1,
		}
		r := ptypes.TestPipelineRun(t, pr)
		err = s.PipelineRunPut(ctx, r)
		require.NoError(err)

		// Get latest run by pipeline
		{
			resp, err := s.PipelineRunGetLatest(ctx, pipeline.Ref.(*pb.Ref_Pipeline_Id).Id)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
			require.Equal(pb.PipelineRun_PENDING, resp.State)
		}

		// Update existing pipeline run
		r.State = pb.PipelineRun_ERROR
		err = s.PipelineRunPut(ctx, r)
		require.NoError(err)

		// Get pipeline run, ID and sequence should not change
		{
			resp, err := s.PipelineRunGetLatest(ctx, pipeline.Ref.(*pb.Ref_Pipeline_Id).Id)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(uint64(1), resp.Sequence)
			require.Equal(r.State, resp.State)
		}
	})

	t.Run("List", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "project"}
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		// Set Pipeline
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(ctx, p)
		require.NoError(err)

		pipeline := &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: p.Id}}

		// Set Pipeline Run
		r := ptypes.TestPipelineRun(t, &pb.PipelineRun{Id: "test1", Pipeline: pipeline})
		err = s.PipelineRunPut(ctx, r)
		require.NoError(err)

		// Set Another Pipeline Run
		r2 := ptypes.TestPipelineRun(t, &pb.PipelineRun{Id: "test2", Pipeline: pipeline})
		err = s.PipelineRunPut(ctx, r2)
		require.NoError(err)

		// Set Another Pipeline Run
		r3 := ptypes.TestPipelineRun(t, &pb.PipelineRun{Id: "test3", Pipeline: pipeline})
		err = s.PipelineRunPut(ctx, r3)
		require.NoError(err)

		// List all runs, check that sequence increments
		{
			resp, err := s.PipelineRunList(ctx, pipeline)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 3)
			require.Equal(r.Id, resp[0].Id)
			require.Equal(uint64(1), resp[0].Sequence)
			require.Equal(r2.Id, resp[1].Id)
			require.Equal(uint64(2), resp[1].Sequence)
			require.Equal(r3.Id, resp[2].Id)
			require.Equal(uint64(3), resp[2].Sequence)
		}
	})
}
