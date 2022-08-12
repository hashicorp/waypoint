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

		const projectName = "testproject"
		const appName = "testapp"
		// Create a project with one app
		require.NoError(s.ProjectPut(&pb.Project{
			Name: projectName,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName},
					Name:    appName,
				},
			},
		}))

		// Read it back
		projectBeforeDelete, err := s.ProjectGet(&pb.Ref_Project{Project: projectName})
		require.NoError(err)
		require.NotNil(projectBeforeDelete)

		// Create a build
		require.NoError(s.BuildPut(false, &pb.Build{
			Id: "testBuild",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(false, &pb.PushedArtifact{
			Id: "testArtifact",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.StatusReportPut(false, &pb.StatusReport{
			Id: "testStatusReport",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.WorkspacePut(&pb.Workspace{
			Name: "testWorkspace",
			Projects: []*pb.Workspace_Project{
				{
					Project: &pb.Ref_Project{Project: projectName},
				},
			},
		}))

		require.NoError(s.PipelinePut(&pb.Pipeline{
			Id:   "testPipeline",
			Name: "testPipeline",
			Owner: &pb.Pipeline_Project{
				Project: &pb.Ref_Project{
					Project: projectName,
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
			Id:        "testTrigger",
			Name:      "testTrigger",
			Project:   &pb.Ref_Project{Project: projectName},
			Workspace: &pb.Ref_Workspace{Workspace: "testWorkspace"},
		}))

		// Delete the project (this should also delete the build)
		err = s.ProjectDelete(&pb.Ref_Project{Project: projectName})
		require.NoError(err)

		// Attempt to get the project again (expected error)
		_, err = s.ProjectGet(&pb.Ref_Project{Project: projectName})
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

		_, err = s.TriggerGet(&pb.Ref_Trigger{Id: "testTrigger"})
		require.Error(err)

		_, err = s.PipelineGet(&pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Owner{&pb.Ref_PipelineOwner{
			Project:      &pb.Ref_Project{Project: projectName},
			PipelineName: "testPipeline",
		}}})
		require.Error(err)
	})
}
