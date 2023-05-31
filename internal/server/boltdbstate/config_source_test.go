package boltdbstate

import (
	"context"
	"testing"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"

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

func TestGlobalConfigSourceDelete(t *testing.T) {
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

	err = s.ConfigSourceDelete(ctx, &pb.ConfigSource{
		Scope:     &pb.ConfigSource_Global{Global: &pb.Ref_Global{}},
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
		Type:      "test",
	})
	require.NoError(err)

	source, err = s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
		Scope:     &pb.GetConfigSourceRequest_Global{Global: &pb.Ref_Global{}},
		Workspace: &pb.Ref_Workspace{Workspace: "default"},
		Type:      "test",
	})

	require.NoError(err)
	require.True(len(source) == 0)
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
	require.Equal(configSourceType, source[0].Type)
	require.Equal(&pb.ConfigSource_Project{
		Project: &pb.Ref_Project{
			Project: projectName,
		},
	}, source[0].Scope)
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
	require.Equal(configSourceType, source[0].Type)
	require.Equal(&pb.ConfigSource_Application{Application: &pb.Ref_Application{
		Application: appName,
		Project:     projectName,
	}}, source[0].Scope)
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
	require.Equal(configSourceType, source[0].Type)
	require.Equal(&pb.ConfigSource_Application{Application: &pb.Ref_Application{
		Application: appName,
		Project:     projectName,
	}}, source[0].Scope)
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
	require.Equal("test", source[0].Type)
	require.Equal(&pb.ConfigSource_Global{Global: &pb.Ref_Global{}}, source[0].Scope)

	// Verify that the second one matches the type we specified
	require.Equal("test2", source[1].Type)
	require.Equal(&pb.ConfigSource_Global{Global: &pb.Ref_Global{}}, source[1].Scope)
}

func TestGlobalProjectAndAppScope(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	const projectName = "testproject"
	const appName = "testapp"
	const workspaceName = "default"
	const configSourceType = "test"

	err := s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete: false,
		Scope: &pb.ConfigSource_Global{
			Global: &pb.Ref_Global{},
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: workspaceName,
		},
		Type: configSourceType,
	})

	err = s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete: false,
		Scope: &pb.ConfigSource_Project{
			Project: &pb.Ref_Project{
				Project: projectName,
			}},
		Workspace: &pb.Ref_Workspace{
			Workspace: workspaceName,
		},
		Type: configSourceType,
	})

	err = s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete: false,
		Scope: &pb.ConfigSource_Application{
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     projectName,
			}},
		Workspace: &pb.Ref_Workspace{
			Workspace: workspaceName,
		},
		Type: configSourceType,
	})

	require.NoError(err)

	source, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
		Scope: &pb.GetConfigSourceRequest_Application{
			Application: &pb.Ref_Application{
				Project:     projectName,
				Application: appName,
			}},
		Workspace: &pb.Ref_Workspace{
			Workspace: workspaceName,
		},
	})

	require.NoError(err)

	require.True(len(source) == 3)

	// Verify that the first one matches the type we specified
	require.Equal(configSourceType, source[0].Type)
	require.Equal(&pb.ConfigSource_Global{
		Global: &pb.Ref_Global{},
	}, source[0].Scope)

	// Verify that the second one matches the type we specified
	require.Equal(configSourceType, source[1].Type)
	require.Equal(&pb.ConfigSource_Project{Project: &pb.Ref_Project{
		Project: projectName,
	}}, source[1].Scope)

	require.Equal(configSourceType, source[2].Type)
	require.Equal(&pb.ConfigSource_Application{Application: &pb.Ref_Application{
		Application: appName,
		Project:     projectName,
	}}, source[2].Scope)
}

func TestProjectAndAppConfigSource(t *testing.T) {
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
		Scope: &pb.ConfigSource_Project{Project: &pb.Ref_Project{
			Project: projectName,
		}},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName},
		Type:      configSourceType,
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
	require.True(len(source) == 2)

	require.Equal(configSourceType, source[0].Type)
	require.Equal(&pb.ConfigSource_Project{
		Project: &pb.Ref_Project{
			Project: projectName,
		}}, source[0].Scope)

	require.Equal(configSourceType, source[1].Type)
	require.Equal(&pb.ConfigSource_Application{Application: &pb.Ref_Application{
		Application: appName,
		Project:     projectName,
	}}, source[1].Scope)
}

func TestMultipleWorkspaceApplicationConfigSources(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	const projectName = "testproject"
	const appName = "testapp"
	const workspaceName = "not-default"
	const workspaceName2 = "the-default"
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

	err = s.ConfigSourceSet(ctx, &pb.ConfigSource{
		Delete: false,
		Scope: &pb.ConfigSource_Application{Application: &pb.Ref_Application{
			Application: appName,
			Project:     projectName,
		}},
		Workspace: &pb.Ref_Workspace{Workspace: workspaceName2},
		Type:      configSourceType,
	})

	require.NoError(err)

	sources, err := s.ConfigSourceGet(ctx, &pb.GetConfigSourceRequest{
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
	require.True(len(sources) == 1)
	require.Equal(configSourceType, sources[0].Type)
	require.Equal(&pb.ConfigSource_Application{Application: &pb.Ref_Application{
		Application: appName,
		Project:     projectName,
	}}, sources[0].Scope)
}
