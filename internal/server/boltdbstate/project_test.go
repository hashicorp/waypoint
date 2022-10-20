package boltdbstate

import (
	"context"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProject(t *testing.T) {
	ctx := context.Background()
	t.Run("create and get and delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		const projectName = "testproject"
		const appName = "testapp"
		// Create a project with one app
		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: projectName,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName},
					Name:    appName,
				},
			},
		}))

		// Read it back
		projectBeforeDelete, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: projectName})
		require.NoError(err)
		require.NotNil(projectBeforeDelete)

		// Create a build, artifact, deployment, release, trigger, workspace, and pipeline
		// Set a config at the project and app scope
		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id: "testBuild",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
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

		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Project{Project: &pb.Ref_Project{Project: projectName}},
			},
			Name:       "testProjectConfig",
			Value:      &pb.ConfigVar_Static{Static: "paladin"},
			Internal:   false,
			NameIsPath: false,
		}))

		require.NoError(s.ConfigSet(&pb.ConfigVar{
			Target: &pb.ConfigVar_Target{
				AppScope: &pb.ConfigVar_Target_Application{Application: &pb.Ref_Application{
					Project:     projectName,
					Application: appName,
				}},
			},
			Name:       "testAppConfig",
			Value:      &pb.ConfigVar_Static{Static: "devops"},
			Internal:   false,
			NameIsPath: false,
		}))

		require.NoError(s.WorkspacePut(&pb.Workspace{
			Name: "testWorkspace",
			Projects: []*pb.Workspace_Project{
				{
					Project: &pb.Ref_Project{Project: projectName},
				},
			},
		}))

		require.NoError(s.PipelinePut(ctx, &pb.Pipeline{
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
		require.NoError(s.ProjectDelete(ctx, &pb.Ref_Project{Project: projectName}))

		// Attempt to get the project again (expected error)
		_, err = s.ProjectGet(ctx, &pb.Ref_Project{Project: projectName})
		require.Error(err)

		// Verify that all builds, artifacts, deployments, releases, status reports,
		// triggers, pipelines and workspaces were deleted, and that configs were unset
		_, err = s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild"}})
		require.Error(err)

		_, err = s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact"}})
		require.Error(err)

		_, err = s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment"}})
		require.Error(err)

		_, err = s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease"}})
		require.Error(err)

		_, err = s.StatusReportGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testStatusReport"}})
		require.Error(err)

		configVars, err := s.ConfigGet(&pb.ConfigGetRequest{})
		require.NoError(err)
		require.Equal(0, len(configVars))

		_, err = s.WorkspaceGet("testWorkspace")
		require.Error(err)

		_, err = s.TriggerGet(&pb.Ref_Trigger{Id: "testTrigger"})
		require.Error(err)

		_, err = s.PipelineGet(ctx, &pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Owner{Owner: &pb.Ref_PipelineOwner{
			Project:      &pb.Ref_Project{Project: projectName},
			PipelineName: "testPipeline",
		}}})
		require.Error(err)
	})

	t.Run("delete project workspaces", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		const projectName1 = "testproject1"
		const appName1 = "testapp1"
		const projectName2 = "testproject2"
		const appName2 = "testapp2"
		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: projectName1,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName1},
					Name:    appName1,
				},
			},
		}))

		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: projectName2,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName2},
					Name:    appName2,
				},
			},
		}))

		require.NoError(s.WorkspacePut(&pb.Workspace{
			Name: "one-project-workspace",
			Projects: []*pb.Workspace_Project{
				{
					Project: &pb.Ref_Project{Project: projectName1},
				},
			},
			ActiveTime: nil,
		}))

		require.NoError(s.WorkspacePut(&pb.Workspace{
			Name: "two-project-workspace",
			Projects: []*pb.Workspace_Project{
				{
					Project: &pb.Ref_Project{Project: projectName1},
				},
				{
					Project: &pb.Ref_Project{Project: projectName2},
				},
			},
			ActiveTime: nil,
		}))

		require.NoError(s.ProjectDelete(ctx, &pb.Ref_Project{Project: projectName1}))

		workspaces, err := s.WorkspaceList()
		require.NoError(err)

		// After the project is deleted, only the 2nd workspace should exist
		require.Equal(1, len(workspaces))
	})

	t.Run("create multi app project and delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		const (
			projectName = "testproject"
			appName1    = "testapp"
			appName2    = "testapp2"
		)
		// Create a project with one app
		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: projectName,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName},
					Name:    appName1,
				},
				{
					Project: &pb.Ref_Project{Project: projectName},
					Name:    appName2,
				},
			},
		}))

		// Read it back
		projectBeforeDelete, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: projectName})
		require.NoError(err)
		require.NotNil(projectBeforeDelete)

		// Create multiple builds for each app
		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id:       "testBuild1App1",
			Sequence: 1,
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id:       "testBuild2App1",
			Sequence: 2,
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id: "testBuild1App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id: "testBuild2App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		// Create multiple artifacts for each app
		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Id: "testArtifact1",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Id: "testArtifact2",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Id: "testArtifact1",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Id: "testArtifact2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		// Create multiple deployments for each app
		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment1App1",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment2App1",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment1App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment2App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		// Create multiple releases for each app
		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease1App1",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease2App1",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease1App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease2App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		// Create multiple status reports for each app
		require.NoError(s.StatusReportPut(false, &pb.StatusReport{
			Id: "testStatusReport1App1",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.StatusReportPut(false, &pb.StatusReport{
			Id: "testStatusReport2App1",
			Application: &pb.Ref_Application{
				Application: appName1,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.StatusReportPut(false, &pb.StatusReport{
			Id: "testStatusReport1App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.StatusReportPut(false, &pb.StatusReport{
			Id: "testStatusReport2App2",
			Application: &pb.Ref_Application{
				Application: appName2,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		// Delete the project (this should also delete the other records)
		require.NoError(s.ProjectDelete(ctx, &pb.Ref_Project{Project: projectName}))

		// Attempt to get the project again (expected error)
		_, err = s.ProjectGet(ctx, &pb.Ref_Project{Project: projectName})
		require.Error(err)

		// Verify that all builds, artifacts, deployments, releases, and status reports were deleted with the project
		_, err = s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild1App1"}})
		require.Error(err)

		_, err = s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild2App1"}})
		require.Error(err)

		_, err = s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild1App2"}})
		require.Error(err)

		_, err = s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild2App2"}})
		require.Error(err)

		_, err = s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact1App1"}})
		require.Error(err)

		_, err = s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact2App1"}})
		require.Error(err)

		_, err = s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact1App2"}})
		require.Error(err)

		_, err = s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact2App2"}})
		require.Error(err)

		_, err = s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment1App1"}})
		require.Error(err)

		_, err = s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment2App1"}})
		require.Error(err)

		_, err = s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment1App2"}})
		require.Error(err)

		_, err = s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment2App2"}})
		require.Error(err)

		_, err = s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease1App1"}})
		require.Error(err)

		_, err = s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease2App1"}})
		require.Error(err)

		_, err = s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease1App2"}})
		require.Error(err)

		_, err = s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease2App2"}})
		require.Error(err)

		_, err = s.StatusReportGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testStatusReport1App1"}})
		require.Error(err)

		_, err = s.StatusReportGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testStatusReport2App1"}})
		require.Error(err)

		_, err = s.StatusReportGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testStatusReport1App2"}})
		require.Error(err)

		_, err = s.StatusReportGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testStatusReport2App2"}})
		require.Error(err)
	})

	t.Run("check sequence # is reset for app operation after deleting and re-initting a project", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		const projectName = "projecttestsequence"
		const appName = "apptestsequence"
		// Create a project with one app
		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: projectName,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName},
					Name:    appName,
				},
			},
		}))

		// Read it back
		projectBeforeDelete, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: projectName})
		require.NoError(err)
		require.NotNil(projectBeforeDelete)

		// Create a build, artifact, deployment, release, trigger, workspace, and pipeline
		// Set a config at the project and app scope
		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id: "testBuild",
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

		// Delete the project (this should also delete the build)
		require.NoError(s.ProjectDelete(ctx, &pb.Ref_Project{Project: projectName}))

		// Attempt to get the project again (expected error)
		_, err = s.ProjectGet(ctx, &pb.Ref_Project{Project: projectName})
		require.Error(err)

		// Re-create the project
		require.NoError(s.ProjectPut(ctx, &pb.Project{
			Name: projectName,
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: projectName},
					Name:    appName,
				},
			},
		}))

		// Read it back
		projectAfterReInit, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: projectName})
		require.NoError(err)
		require.NotNil(projectAfterReInit)

		// Create new build after the project is re-initialized
		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id: "testBuild1",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id: "testBuild2",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.BuildPut(ctx, false, &pb.Build{
			Id: "testBuild3",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Id: "testArtifact1",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Id: "testArtifact2",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ArtifactPut(ctx, false, &pb.PushedArtifact{
			Id: "testArtifact3",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment1",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment2",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.DeploymentPut(false, &pb.Deployment{
			Id: "testDeployment3",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease1",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease2",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		require.NoError(s.ReleasePut(false, &pb.Release{
			Id: "testRelease3",
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
			Workspace: &pb.Ref_Workspace{Workspace: "default"},
		}))

		// Verify that the app operation sequence is 1, since it is the 1st operation for the new (albeit
		// re-initialized it was deleted) project
		build1, err := s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild1"}})
		require.NoError(err)
		require.Equal(1, int(build1.Sequence))
		build2, err := s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild2"}})
		require.NoError(err)
		require.Equal(2, int(build2.Sequence))
		build3, err := s.BuildGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testBuild3"}})
		require.NoError(err)
		require.Equal(3, int(build3.Sequence))

		artifact1, err := s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact1"}})
		require.NoError(err)
		require.Equal(1, int(artifact1.Sequence))

		artifact2, err := s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact2"}})
		require.NoError(err)
		require.Equal(2, int(artifact2.Sequence))

		artifact3, err := s.ArtifactGet(ctx, &pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testArtifact3"}})
		require.NoError(err)
		require.Equal(3, int(artifact3.Sequence))

		deployment1, err := s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment1"}})
		require.NoError(err)
		require.Equal(1, int(deployment1.Sequence))

		deployment2, err := s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment2"}})
		require.NoError(err)
		require.Equal(2, int(deployment2.Sequence))

		deployment3, err := s.DeploymentGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testDeployment3"}})
		require.NoError(err)
		require.Equal(3, int(deployment3.Sequence))

		release1, err := s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease1"}})
		require.NoError(err)
		require.Equal(1, int(release1.Sequence))

		release2, err := s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease2"}})
		require.NoError(err)
		require.Equal(2, int(release2.Sequence))

		release3, err := s.ReleaseGet(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: "testRelease3"}})
		require.NoError(err)
		require.Equal(3, int(release3.Sequence))
	})
}
