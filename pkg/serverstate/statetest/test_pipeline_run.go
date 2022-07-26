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
		require.NoError(err) // no job id set

		// Set Pipeline Run
		pr := &pb.PipelineRun{
			Pipeline: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: &pb.Ref_PipelineId{
						Id: p.Id,
					},
				},
			},
		}
		r := ptypes.TestPipelineRun(t, pr)
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// Get run by pipeline and sequence
		{
			resp, err := s.PipelineRunGet(pr.Pipeline, r.Sequence)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Id, pr.Id)
			require.Equal(resp.Sequence, r.Sequence)
		}

		pr = &pb.PipelineRun{
			Id:       "test-run2",
			Sequence: 2,
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
		r := ptypes.TestPipelineRun(t, nil)
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// List
		{
			resp, err := s.PipelineRunList(pipeline)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
			require.Equal(resp[0].Id, r.Id)
		}

		pr := &pb.PipelineRun{
			Id:       "test-run2",
			Sequence: 2,
		}

		// Set Another Pipeline Run
		r = ptypes.TestPipelineRun(t, pr)
		err = s.PipelineRunPut(r)
		require.NoError(err)

		// List multiple runs
		{
			resp, err := s.PipelineRunList(pipeline)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 2)
			require.Equal(resp[1].Id, pr.Id)
		}
	})
}
