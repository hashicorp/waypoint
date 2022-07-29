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
		pr := &pb.PipelineRun{
			Id:       "test-pr",
			Pipeline: pipeline,
		}
		r := ptypes.TestPipelineRun(t, pr)
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Get run by pipeline and sequence
		{
			resp, err := s.PipelineRunGet(pipeline, 1)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(pr.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
		}

		// Set another pipeline run
		r = ptypes.TestPipelineRun(t, &pb.PipelineRun{
			Id:       "test-pr2",
			Pipeline: pipeline,
		})
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Get run by pipeline and sequence, should auto increment
		{
			resp, err := s.PipelineRunGet(pipeline, 2)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
		}

		// Set another pipeline run
		r = ptypes.TestPipelineRun(t, &pb.PipelineRun{
			Id:       "test-pr3",
			Pipeline: pipeline,
		})
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Get run by pipeline and sequence, should auto increment
		{
			resp, err := s.PipelineRunGet(pipeline, 3)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(r.Id, resp.Id)
			require.Equal(r.Sequence, resp.Sequence)
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
		r := ptypes.TestPipelineRun(t, &pb.PipelineRun{
			Id: "1",
		})
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Set Another Pipeline Run
		r2 := ptypes.TestPipelineRun(t, &pb.PipelineRun{
			Id: "2",
		})
		err = s.PipelineRunPut(r2)
		require.NoError(err)

		// Set Another Pipeline Run
		r3 := ptypes.TestPipelineRun(t, &pb.PipelineRun{
			Id: "3",
		})
		err = s.PipelineRunPut(r3)
		require.NoError(err)

		// List multiple runs, check sequence increments
		{
			resp, err := s.PipelineRunList(pipeline)
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
