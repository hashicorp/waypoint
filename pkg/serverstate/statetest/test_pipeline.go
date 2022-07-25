package statetest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["pipeline"] = []testFunc{
		TestPipeline,
	}
}

func TestPipeline(t *testing.T, factory Factory, restartF RestartFactory) {
	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err) // no job id set

		// Set again, should overwrite and not error
		err = s.PipelinePut(p)
		require.NoError(err)

		// Get exact by id
		{
			resp, err := s.PipelineGet(&pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Id{
					Id: &pb.Ref_PipelineId{Id: p.Id},
				},
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Delete
		require.NoError(s.PipelineDelete(&pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: &pb.Ref_PipelineId{Id: p.Id},
			},
		}))
	})

	t.Run("Put: no steps", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		p := ptypes.TestPipeline(t, nil)
		p.Steps = nil
		err := s.PipelinePut(p)
		require.Error(err)
		require.Equal(codes.FailedPrecondition, status.Code(err))
	})

	t.Run("Put: no root step", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		p := ptypes.TestPipeline(t, nil)
		p.Steps = map[string]*pb.Pipeline_Step{
			"A": {
				Name:      "A",
				DependsOn: []string{"B"},
			},
			"B": {
				Name:      "B",
				DependsOn: []string{"C"},
			},
		}
		err := s.PipelinePut(p)
		require.Error(err)
		require.Equal(codes.FailedPrecondition, status.Code(err))
	})

	t.Run("Put: cycle", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		p := ptypes.TestPipeline(t, nil)
		p.Steps = map[string]*pb.Pipeline_Step{
			"root": {
				Name: "root",
			},
			"A": {
				Name:      "A",
				DependsOn: []string{"root", "B"},
			},
			"B": {
				Name:      "B",
				DependsOn: []string{"A"},
			},
		}
		err := s.PipelinePut(p)
		require.Error(err)
		require.Equal(codes.FailedPrecondition, status.Code(err))
	})

	t.Run("Put: non-existent dependency", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		p := ptypes.TestPipeline(t, nil)
		p.Steps = map[string]*pb.Pipeline_Step{
			"A": {
				Name: "A",
			},
			"B": {
				Name:      "B",
				DependsOn: []string{"C"},
			},
		}
		err := s.PipelinePut(p)
		require.Error(err)
		require.Equal(codes.FailedPrecondition, status.Code(err))
	})

	t.Run("Put: with no Id by Pipeline Owner and Pipeline Name Ref updates existing field", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set a few pipelines
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)

		// Update the pipeline
		p.Steps["second"] = &pb.Pipeline_Step{
			Name:      "second",
			DependsOn: []string{"root"},
			Kind: &pb.Pipeline_Step_Exec_{
				Exec: &pb.Pipeline_Step_Exec{
					Image: "hashicorp/waypoint",
				},
			},
		}
		p.Id = ""
		err = s.PipelinePut(p)
		require.NoError(err)

		// Should only be 1 pipeline
		// List should return one pipeline
		{
			resp, err := s.PipelineList(&pb.Ref_Project{
				Project: "project",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
			require.Equal(resp[0].Id, "test")
		}
	})

	t.Run("Put: with existing Id by Pipeline Owner and Pipeline Name Ref updates existing field", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set a few pipelines
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)

		// Update the pipeline
		p.Steps["second"] = &pb.Pipeline_Step{
			Name:      "second",
			DependsOn: []string{"root"},
			Kind: &pb.Pipeline_Step_Exec_{
				Exec: &pb.Pipeline_Step_Exec{
					Image: "hashicorp/waypoint",
				},
			},
		}
		p.Id = "test"

		err = s.PipelinePut(p)
		require.NoError(err)

		// Should only be 1 pipeline
		// List should return one pipeline
		{
			resp, err := s.PipelineList(&pb.Ref_Project{
				Project: "project",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
			require.Equal(resp[0].Id, p.Id)
		}
	})

	t.Run("Put: with no Id by Pipeline Owner and Pipeline Name Ref inserts a new pipeline", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set a few pipelines
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)

		// Update the pipeline
		p.Steps["second"] = &pb.Pipeline_Step{
			Name:      "second",
			DependsOn: []string{"root"},
			Kind: &pb.Pipeline_Step_Exec_{
				Exec: &pb.Pipeline_Step_Exec{
					Image: "hashicorp/waypoint",
				},
			},
		}
		p.Id = "two"
		p.Name = "two"

		err = s.PipelinePut(p)
		require.NoError(err)

		// Should only be 2 pipelines
		{
			resp, err := s.PipelineList(&pb.Ref_Project{
				Project: "project",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 2)
			require.Equal(resp[0].Id, "test")
			require.Equal(resp[1].Id, "two")
		}

		// Another new one
		p2 := ptypes.TestPipeline(t, &pb.Pipeline{
			Id:   "mario",
			Name: "mario",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
		})
		err = s.PipelinePut(p2)
		require.NoError(err)

		// Should be three
		{
			resp, err := s.PipelineList(&pb.Ref_Project{
				Project: "project",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 3)
		}
	})

	t.Run("Get: by Pipeline Owner and Pipeline Name Ref", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set a few pipelines
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)

		// Another one
		p2 := ptypes.TestPipeline(t, &pb.Pipeline{
			Id:   "mario",
			Name: "mario",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
		})
		err = s.PipelinePut(p2)
		require.NoError(err)

		// a third, same pipeline name but different project
		p3 := ptypes.TestPipeline(t, &pb.Pipeline{
			Id:   "testtest",
			Name: "mario",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "nintendo",
				},
			},
		})
		err = s.PipelinePut(p3)
		require.NoError(err)

		// Get pipeline by Owner Ref
		{
			resp, err := s.PipelineGet(&pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{Project: &pb.Ref_Project{Project: "project"}, PipelineName: "mario"},
				},
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Name, "mario")
			require.Equal(resp.Id, "mario")
		}

		// Get pipeline by Owner Ref
		{
			resp, err := s.PipelineGet(&pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{Project: &pb.Ref_Project{Project: "nintendo"}, PipelineName: "mario"},
				},
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Name, "mario")
			require.Equal(resp.Id, "testtest")
		}
	})

	t.Run("Get: by missing Pipeline Owner and Pipeline Name Ref", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)

		// Get pipeline by Owner Ref should be nothing
		{
			resp, err := s.PipelineGet(&pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{Project: &pb.Ref_Project{Project: "nope"}, PipelineName: "nope"},
				},
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
			require.Nil(resp)
		}
	})

	t.Run("List", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		p := ptypes.TestPipeline(t, nil)
		err := s.PipelinePut(p)
		require.NoError(err)

		// List should return one pipeline
		{
			resp, err := s.PipelineList(&pb.Ref_Project{
				Project: "project",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 1)
		}

		// Another one
		p2 := ptypes.TestPipeline(t, &pb.Pipeline{
			Id:   "mario",
			Name: "mario",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
		})
		err = s.PipelinePut(p2)
		require.NoError(err)

		// a third
		p3 := ptypes.TestPipeline(t, &pb.Pipeline{
			Id:   "testtwo",
			Name: "testtwo",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "project",
				},
			},
		})
		err = s.PipelinePut(p3)
		require.NoError(err)

		// List should return three pipelines
		{
			resp, err := s.PipelineList(&pb.Ref_Project{
				Project: "project",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp, 3)
		}
	})
}
