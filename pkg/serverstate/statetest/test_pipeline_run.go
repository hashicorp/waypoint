package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["pipeline_run"] = []testFunc{
		TestPipelineRun,
	}
}

func TestPipelineRun(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set Pipeline
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)
		pipeline := &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: &pb.Ref_PipelineId{Id: p.Id}}}

		// Set Pipeline Run
		r := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Get run by pipeline and sequence
		{
			resp, err := s.PipelineRunGet(pipeline, 1)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
			require.Equal(pipeline.Ref.(*pb.Ref_Pipeline_Id).Id.Id, resp.Pipeline.Ref.(*pb.Ref_Pipeline_Id).Id.Id)
		}

		// Get run by Id
		{
			resp, err := s.PipelineRunGetById(r.Id)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
		}

		// Set another pipeline run
		r2 := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		err = s.PipelineRunPut(r2)
		require.NoError(err)

		// Get run by pipeline and sequence, should auto increment
		{
			resp, err := s.PipelineRunGet(pipeline, 2)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r2.Id, resp.Id)
			require.Equal(r2.Sequence, resp.Sequence)
		}

		// Set another pipeline run
		latest_r := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		err = s.PipelineRunPut(latest_r)
		require.NoError(err)

		// Get run by pipeline and sequence, should auto increment
		{
			resp, err := s.PipelineRunGet(pipeline, 3)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(latest_r.Id, resp.Id)
			require.Equal(latest_r.Sequence, resp.Sequence)
		}

		// Get latest run by pipeline ID
		{
			resp, err := s.PipelineRunGetLatest(pipeline.Ref.(*pb.Ref_Pipeline_Id).Id.Id)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(latest_r.Id, resp.Id)
			require.Equal(latest_r.Sequence, resp.Sequence)
		}
	})

	t.Run("Update existing run", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set Pipeline
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)
		pipeline := &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: &pb.Ref_PipelineId{Id: p.Id}}}

		// Set Pipeline Run
		pr := &pb.PipelineRun{
			Id:       "test-pr",
			Pipeline: pipeline,
			Sequence: 1,
		}
		r := ptypes.TestPipelineRun(t, pr)
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Get latest run by pipeline
		{
			resp, err := s.PipelineRunGetLatest(pipeline.Ref.(*pb.Ref_Pipeline_Id).Id.Id)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
			require.Equal(pb.PipelineRun_PENDING, resp.State)
		}

		// Update existing pipeline run
		r.State = pb.PipelineRun_ERROR
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Get pipeline run, ID and sequence should not change
		{
			resp, err := s.PipelineRunGetLatest(pipeline.Ref.(*pb.Ref_Pipeline_Id).Id.Id)
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

		// Set Pipeline
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)

		pipeline := &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: &pb.Ref_PipelineId{Id: p.Id}}}

		// Set Pipeline Run
		r := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Set Another Pipeline Run
		r2 := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		err = s.PipelineRunPut(r2)
		require.NoError(err)

		// Set Another Pipeline Run
		r3 := ptypes.TestPipelineRun(t, &pb.PipelineRun{Pipeline: pipeline})
		err = s.PipelineRunPut(r3)
		require.NoError(err)

		// List all runs, check that sequence increments
		{
			resp, err := s.PipelineRunList(pipeline)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 3)
			require.NotEqual(r.Id, r2.Id, r3.Id)
			require.Equal(uint64(1), resp[0].Sequence)
			require.Equal(uint64(2), resp[1].Sequence)
			require.Equal(uint64(3), resp[2].Sequence)
		}
	})
}
