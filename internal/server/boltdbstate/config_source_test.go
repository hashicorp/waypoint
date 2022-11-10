package boltdbstate

import (
	"context"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobalConfigSource(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	err := s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete:    false,
		Scope:     &pb.ConfigSource_Global{Global: &pb.Ref_Global{}},
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
		Type:      "test",
	})

	require.NoError(err)

	source, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
		Scope:     &pb.GetConfigSourceRequest_Global{Global: &pb.Ref_Global{}},
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
		Type:      "test",
	})

	require.NoError(err)
	require.True(len(source) > 0)
	require.Equal(source[0].Type, "test")
	require.Equal(source[0].Scope, &pb.ConfigSource_Global{Global: &pb.Ref_Global{}})
}

func TestProjectConfigSource(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	const projectName = "testproject"
	const appName = "testapp"
	const workspaceName = "default"
	const configSourceType = "test"

	err := s.ProjectPut(ctx, &pb.Project{
		Name: "test-project",
		Applications: []*pb.Application{
			{
				Project: &pb.Ref_Project{Project: projectName},
				Name:    appName,
			},
		},
	})

	err = s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete:    false,
		Scope:     &pb.ConfigSource_Project{Project: &pb.Ref_Project{Project: projectName}},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName},
		Type:      configSourceType,
	})

	require.NoError(err)

	source, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
		Scope:     &pb.GetConfigSourceRequest_Project{Project: &pb.Ref_Project{Project: projectName}},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName},
		Type:      configSourceType,
	})

	require.NoError(err)
	require.True(len(source) > 0)
	require.Equal(source[0].Type, configSourceType)
	require.Equal(source[0].Scope, &pb.ConfigSource_Project{Project: &pb.Ref_Project{Project: projectName}})
}

func TestAppConfigSource(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	const projectName = "testproject"
	const appName = "testapp"
	const workspaceName = "default"
	const configSourceType = "test"

	err := s.ProjectPut(ctx, &pb.Project{
		Name: "test-project",
		Applications: []*pb.Application{
			{
				Project: &pb.Ref_Project{Project: projectName},
				Name:    appName,
			},
		},
	})

	err = s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete: false,
		Scope: &pb.ConfigSource_Application{Application: &pb.Ref_Application{
			Application: appName,
			Project:     projectName,
		}},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName},
		Type:      configSourceType,
	})

	require.NoError(err)

	source, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
		Scope: &pb.GetConfigSourceRequest_Application{
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
		},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName},
		Type:      configSourceType,
	})

	require.NoError(err)
	require.True(len(source) > 0)
	require.Equal(source[0].Type, configSourceType)
	require.Equal(source[0].Scope, &pb.ConfigSource_Application{Application: &pb.Ref_Application{
		Application: appName,
		Project:     projectName,
	}})
}

func TestWorkspaceProjectConfigSource(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	const projectName = "testproject"
	const appName = "testapp"
	const workspaceName = "not-default"
	const configSourceType = "test"

	err := s.ProjectPut(ctx, &pb.Project{
		Name: "test-project",
		Applications: []*pb.Application{
			{
				Project: &pb.Ref_Project{Project: projectName},
				Name:    appName,
			},
		},
	})

	err = s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete: false,
		Scope: &pb.ConfigSource_Application{Application: &pb.Ref_Application{
			Application: appName,
			Project:     projectName,
		}},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName},
		Type:      configSourceType,
	})

	require.NoError(err)

	source, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
		Scope: &pb.GetConfigSourceRequest_Application{
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			},
		},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName},
		Type:      configSourceType,
	})

	require.NoError(err)
	require.True(len(source) > 0)
	require.Equal(source[0].Type, configSourceType)
	require.Equal(source[0].Scope, &pb.ConfigSource_Application{Application: &pb.Ref_Application{
		Application: appName,
		Project:     projectName,
	}})
}

func TestMultipleGlobalConfigSources(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	err := s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete:    false,
		Scope:     &pb.ConfigSource_Global{Global: &pb.Ref_Global{}},
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
		Type:      "test",
	})

	err = s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete:    false,
		Scope:     &pb.ConfigSource_Global{Global: &pb.Ref_Global{}},
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
		Type:      "test2",
	})

	require.NoError(err)

	source, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
		Scope:     &pb.GetConfigSourceRequest_Global{Global: &pb.Ref_Global{}},
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
	})

	require.NoError(err)

	// Ensure that since we created 2 global config sources, we get 2 back
	require.True(len(source) == 2)

	// Verify that the first one matches the type we specified
	require.Equal(source[0].Type, "test")
	require.Equal(source[0].Scope, &pb.ConfigSource_Global{Global: &pb.Ref_Global{}})

	// Verify that the second one matches the type we specified
	require.Equal(source[1].Type, "test2")
	require.Equal(source[1].Scope, &pb.ConfigSource_Global{Global: &pb.Ref_Global{}})
}
