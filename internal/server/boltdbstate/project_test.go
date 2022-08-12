package boltdbstate

import (
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProject(t *testing.T) {
	t.Run("create and get and delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a project with one app
		require.NoError(s.ProjectPut(&pb.Project{
			Name: "test",
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: "test"},
					Name:    "testApp",
				},
			},
		}))

		// Read it back
		projectBeforeDelete, err := s.ProjectGet(&pb.Ref_Project{Project: "test"})
		require.NoError(err)
		require.NotNil(projectBeforeDelete)

		// Create a build
		require.NoError(s.BuildPut(false, &pb.Build{
			Id: "testBuild",
			Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "testProject",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(false, &pb.PushedArtifact{
			Id: "testArtifact",
			Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "testProject",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment",
			Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "testProject",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease",
			Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "testProject",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.StatusReportPut(false, &pb.StatusReport{
			Id: "testStatusReport",
			Application: &pb.Ref_Application{
				Application: "testApp",
				Project:     "testProject",
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.WorkspacePut(&pb.Workspace{
			Name: "testWorkspace",
			Projects: []*pb.Workspace_Project{
				{
					Project: &pb.Ref_Project{Project: "testProject"},
				},
			},
		}))

		require.NoError(s.PipelinePut(&pb.Pipeline{
			Id:   "testPipeline",
			Name: "testPipeline",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: "testProject",
				},
			},
			Steps: map[string]*pb.Pipeline_Step{
				"testStep": {
					Name: "testStep",
					Kind: &pb.Pipeline_Step_Up_{
						Up: &pb.Pipeline_Step_Up{},
					},
				},
			},
		}))

		require.NoError(s.TriggerPut(&pb.Trigger{
			Id:      "testTrigger",
			Name:    "testTrigger",
			Project: &pb.Ref_Project{Project: "testProject"},
		}))

		// Delete the project (this should also delete the build)
		err = s.ProjectDelete(&pb.Ref_Project{Project: "testProject"})
		require.NoError(err)

		// Attempt to get the project again (expected error)
		_, err = s.ProjectGet(&pb.Ref_Project{Project: "testProject"})
		require.Error(err)

		_, err = s.BuildGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild"}})
		require.Error(err)

		_, err = s.ArtifactGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact"}})
		require.Error(err)

		_, err = s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment"}})
		require.Error(err)

		_, err = s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease"}})
		require.Error(err)

		_, err = s.StatusReportGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testStatusReport"}})
		require.Error(err)

		_, err = s.WorkspaceGet("testWorkspace")
		require.Error(err)

		_, err = s.PipelineGet(&pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{Id: &pb.Ref_PipelineId{Id: "testPipeline"}}})
		require.Error(err)
	})
}
