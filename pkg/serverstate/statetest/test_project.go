package statetest

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/pkg/jsonpb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["project"] = []testFunc{
		TestProject,
		TestProjectPollPeek,
		TestProjectPollComplete,
		TestProjectListWorkspaces,
		TestProjectGetSetAllProperties,
		TestProjectGetSetAllPropertiesSansVariables,
		TestProjectCanTransitionDataSource,
	}
}

func TestProject(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		_, err := s.ProjectGet(ctx, &pb.Ref_Project{
			Project: "foo",
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Basic Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		name := "AbCdE"

		// Set
		err := s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: name,
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
				Project: name,
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get case insensitive
		{
			resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
				Project: strings.ToLower(name),
			})
			require.NoError(err, "unable to use case insensitive name for: %s", name)
			require.NotNil(resp)
		}

		// Create another one so we're sure that List can see more than one.

		// Set
		err = s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: name + "2",
		}))
		require.NoError(err)

		// List
		{
			resp, err := s.ProjectList(ctx)
			require.NoError(err)
			require.Len(resp, 2)
		}
	})

	t.Run("Put does not modify applications", func(t *testing.T) {
		require := require.New(t)

		const name = "AbCdE"
		ref := &pb.Ref_Project{Project: name}

		s := factory(t)
		defer s.Close()

		// Set
		proj := serverptypes.TestProject(t, &pb.Project{Name: name})
		err := s.ProjectPut(ctx, proj)
		require.NoError(err)
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Name:    "test",
			Project: ref,
		}))
		require.NoError(err)
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Name:    "test2",
			Project: ref,
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.False(resp.RemoteEnabled)
			require.Len(resp.Applications, 2)
		}

		// Update the project
		proj.RemoteEnabled = true
		require.NoError(s.ProjectPut(ctx, proj))

		// Get exact
		{
			resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.True(resp.RemoteEnabled)
			require.Len(resp.Applications, 2)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "AbCdE",
		}))
		require.NoError(err)

		// Read
		resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
			Project: "AbCdE",
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.ProjectDelete(ctx, &pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.ProjectGet(ctx, &pb.Ref_Project{
				Project: "AbCdE",
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.ProjectList(ctx)
			require.NoError(err)
			require.Len(resp, 0)
		}
	})
}

func TestProjectGetSetAllPropertiesSansVariables(t *testing.T, f Factory, rf RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := f(t)
	defer s.Close()

	// A project with all the properties set
	initialProject := &pb.Project{
		Name: "complex project",
		Applications: []*pb.Application{{
			Name:    "complex project",
			Project: &pb.Ref_Project{Project: "complex project"},
		}},
		RemoteEnabled: true,
		DataSource: &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Git{
				Git: &pb.Job_Git{
					Url:                      "https://github.com/hashicorp/test",
					Ref:                      "main",
					Path:                     "/test",
					IgnoreChangesOutsidePath: true,
					RecurseSubmodules:        1,
					Auth: &pb.Job_Git_Ssh{
						Ssh: &pb.Job_Git_SSH{
							PrivateKeyPem: []byte("private key"),
							Password:      "password",
							User:          "user",
						},
					},
				},
			},
		},
		DataSourcePoll: &pb.Project_Poll{
			Enabled:  true,
			Interval: "1h",
		},
		WaypointHcl:       []byte("hcl bytes"),
		WaypointHclFormat: pb.Hcl_JSON,
		FileChangeSignal:  "HUP",
		StatusReportPoll: &pb.Project_AppStatusPoll{
			Enabled:  true,
			Interval: "1h",
		},
	}

	initialJsonBytes, err := jsonpb.Marshal(initialProject)
	require.NoError(err)
	initialJsonStr := string(initialJsonBytes)

	// Set
	err = s.ProjectPut(ctx, initialProject)
	require.NoError(err)

	// Get
	resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
		Project: initialProject.Name,
	})
	require.NoError(err)
	require.NotNil(resp)

	// Compare the two
	respJsonBytes, err := jsonpb.Marshal(resp)
	require.NoError(err)
	respJsonStr := string(respJsonBytes)

	require.Equal(initialJsonStr, respJsonStr)
}

func TestProjectCanTransitionDataSource(t *testing.T, f Factory, rf RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := f(t)
	defer s.Close()

	project := &pb.Project{
		Name: "testProject",
	}

	// Can post initially
	{
		// Set
		err := s.ProjectPut(ctx, project)
		require.NoError(err)

		// Get
		resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
			Project: project.Name,
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Nil(resp.DataSource)
	}

	// TODO(izaak) Reenable this when all server implementations support an explicit local datasource state
	// rather than assuming that the nil state is datasource.

	// Can set ds to local
	//{
	//	project.DataSource = &pb.Job_DataSource{
	//		Source: &pb.Job_DataSource_Local{Local: &pb.Job_Local{}},
	//	}
	//
	//	// Set
	//	err := s.ProjectPut(ctx, project)
	//	require.NoError(err)
	//
	//	// Get
	//	resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
	//		Project: project.Name,
	//	})
	//	require.NoError(err)
	//	require.NotNil(resp)
	//	require.NotNil(resp.DataSource)
	//	require.IsType(&pb.Job_DataSource_Local{}, resp.DataSource.Source)
	//}

	// Can set ds to git
	{
		project.DataSource = &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Git{Git: &pb.Job_Git{
				Url:                      "test",
				Ref:                      "test",
				Path:                     "test",
				IgnoreChangesOutsidePath: true,
				RecurseSubmodules:        1,
				Auth:                     nil,
			}},
		}

		// Set
		err := s.ProjectPut(ctx, project)
		require.NoError(err)

		// Get
		resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
			Project: project.Name,
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotNil(resp.DataSource)
		require.IsType(&pb.Job_DataSource_Git{}, resp.DataSource.Source)
	}

	// Can set ds back to nil
	{
		project.DataSource = nil

		// Set
		err := s.ProjectPut(ctx, project)
		require.NoError(err)

		// Get
		resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
			Project: project.Name,
		})
		require.NoError(err)
		require.NotNil(resp)
		require.Nil(resp.DataSource)
	}

	// NOTE(izaak): we should probably make it possible to set DS to nil, but in practice nil is interpreted
	// as local by the cli and likely other parts of the server, so setting it to local explicitly
	// rather than nil is what we do.
}

func TestProjectGetSetAllProperties(t *testing.T, f Factory, rf RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := f(t)
	defer s.Close()

	// A project with all the properties set
	initialProject := &pb.Project{
		Name: "complex project",
		Applications: []*pb.Application{{
			Name:    "complex project",
			Project: &pb.Ref_Project{Project: "complex project"},
		}},
		RemoteEnabled: true,
		DataSource: &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Git{
				Git: &pb.Job_Git{
					Url:                      "https://github.com/hashicorp/test",
					Ref:                      "main",
					Path:                     "/test",
					IgnoreChangesOutsidePath: true,
					RecurseSubmodules:        1,
					Auth: &pb.Job_Git_Ssh{
						Ssh: &pb.Job_Git_SSH{
							PrivateKeyPem: []byte("private key"),
							Password:      "password",
							User:          "user",
						},
					},
				},
			},
		},
		DataSourcePoll: &pb.Project_Poll{
			Enabled:  true,
			Interval: "1h",
		},
		WaypointHcl:       []byte("hcl bytes"),
		WaypointHclFormat: pb.Hcl_JSON,
		FileChangeSignal:  "HUP",
		Variables: []*pb.Variable{{
			Name: "test-variable",
			Value: &pb.Variable_Str{
				Str: "variable-value",
			},
			Source: &pb.Variable_Vcs{
				Vcs: &pb.Variable_VCS{
					FileName: "test.file",
					HclRange: &pb.Variable_HclRange{
						Filename: "test.file",
						Start: &pb.Variable_HclPos{
							Line:   1,
							Column: 2,
							Byte:   3,
						},
						End: &pb.Variable_HclPos{
							Line:   4,
							Column: 5,
							Byte:   6,
						},
					},
				},
			},
			FinalValue: &pb.Variable_FinalValue{
				Value: &pb.Variable_FinalValue_Sensitive{
					Sensitive: "sensitive-value",
				},
				Source: pb.Variable_FinalValue_DEFAULT,
			},
			Sensitive: true,
		}},
		StatusReportPoll: &pb.Project_AppStatusPoll{
			Enabled:  true,
			Interval: "1h",
		},
	}

	initialJsonBytes, err := jsonpb.Marshal(initialProject)
	require.NoError(err)
	initialJsonStr := string(initialJsonBytes)

	// Set
	err = s.ProjectPut(ctx, initialProject)
	require.NoError(err)

	// Get
	resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
		Project: initialProject.Name,
	})
	require.NoError(err)
	require.NotNil(resp)

	// Compare the two
	respJsonBytes, err := jsonpb.Marshal(resp)
	require.NoError(err)
	respJsonStr := string(respJsonBytes)

	require.Equal(initialJsonStr, respJsonStr)

	t.Run("can delete all input vars", func(t *testing.T) {
		initialProject.Variables = []*pb.Variable{}
		err = s.ProjectPut(ctx, initialProject)
		require.NoError(err)

		resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
			Project: initialProject.Name,
		})
		require.NoError(err)
		require.NotNil(resp)

		require.Empty(resp.Variables)
	})

}

func TestProjectPollPeek(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("returns nil if no values", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		v, _, err := s.ProjectPollPeek(ctx, nil)
		require.NoError(err)
		require.Nil(v)
	})

	t.Run("returns next to poll", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "10s",
			},
		})))

		// Set another later
		time.Sleep(10 * time.Millisecond)
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "B",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "10s",
			},
		})))

		// Get exact
		{
			resp, t, err := s.ProjectPollPeek(ctx, nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
			require.False(t.IsZero())
		}
	})

	t.Run("watchset triggers from empty to available", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		ws := memdb.NewWatchSet()
		v, _, err := s.ProjectPollPeek(ctx, ws)
		require.NoError(err)
		require.Nil(v)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Set
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "10s",
			},
		})))

		// Should be triggered.
		require.False(ws.Watch(time.After(2 * time.Second)))

		// Get exact
		{
			resp, t, err := s.ProjectPollPeek(ctx, nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
			require.False(t.IsZero())
		}
	})

	t.Run("watchset triggers when records change", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5s",
			},
		})))

		// Set another later
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "B",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5m", // 5 MINUTES, longer than A
			},
		})))

		// Get
		pA, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: "A"})
		require.NoError(err)
		require.NotNil(pA)
		pB, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: "B"})
		require.NoError(err)
		require.NotNil(pB)

		// Complete both first
		now := time.Now()
		require.NoError(s.ProjectPollComplete(ctx, pA, now))
		require.NoError(s.ProjectPollComplete(ctx, pB, now))

		// Peek, we should get A
		ws := memdb.NewWatchSet()
		p, ts, err := s.ProjectPollPeek(ctx, ws)
		require.NoError(err)
		require.NotNil(p)
		require.Equal("A", p.Name)
		require.False(ts.IsZero())

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Set
		require.NoError(s.ProjectPollComplete(ctx, pA, now.Add(1*time.Second)))

		// Should be triggered.
		require.False(ws.Watch(time.After(2 * time.Second)))

		// Get exact
		{
			resp, t, err := s.ProjectPollPeek(ctx, nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
			require.False(t.IsZero())
		}
	})
}

func TestProjectPollComplete(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("returns nil for project that doesn't exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		require.NoError(s.ProjectPollComplete(ctx, &pb.Project{Name: "NOPE"}, time.Now()))
	})

	t.Run("does nothing for project that has polling disabled", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled: false,
			},
		})))

		// Get
		p, err := s.ProjectGet(ctx, &pb.Ref_Project{
			Project: "A",
		})
		require.NoError(err)
		require.NotNil(p)

		// No error
		require.NoError(s.ProjectPollComplete(ctx, p, time.Now()))

		// Peek does nothing
		v, _, err := s.ProjectPollPeek(ctx, nil)
		require.NoError(err)
		require.Nil(v)
	})

	t.Run("schedules the next poll time", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5s",
			},
		})))

		// Set another later
		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
			Name: "B",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5m", // 5 MINUTES, longer than A
			},
		})))

		// Get
		pA, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: "A"})
		require.NoError(err)
		require.NotNil(pA)
		pB, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: "B"})
		require.NoError(err)
		require.NotNil(pB)

		// Complete both first
		now := time.Now()
		require.NoError(s.ProjectPollComplete(ctx, pA, now))
		require.NoError(s.ProjectPollComplete(ctx, pB, now))

		// Peek should return A, lower interval
		{
			resp, t, err := s.ProjectPollPeek(ctx, nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
			require.False(t.IsZero())
		}

		// Complete again, a minute later. The result should be A again
		// because of the lower interval.
		{
			require.NoError(s.ProjectPollComplete(ctx, pA, now.Add(1*time.Minute)))

			resp, t, err := s.ProjectPollPeek(ctx, nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
			require.False(t.IsZero())
		}

		// Complete A, now 6 minutes later. The result should be B now.
		{
			require.NoError(s.ProjectPollComplete(ctx, pA, now.Add(6*time.Minute)))

			resp, t, err := s.ProjectPollPeek(ctx, nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("B", resp.Name)
			require.False(t.IsZero())
		}
	})
}

func TestProjectListWorkspaces(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("empty for non-existent project", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		result, err := s.ProjectListWorkspaces(ctx, &pb.Ref_Project{Project: "nope"})
		require.NoError(err)
		require.Empty(result)
	})

	t.Run("returns only the workspaces a project is in", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create a build
		require.NoError(s.BuildPut(false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "1",
			Workspace: &pb.Ref_Workspace{
				Workspace: "A",
			},
		})))
		require.NoError(s.BuildPut(false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "2",
			Workspace: &pb.Ref_Workspace{
				Workspace: "B",
			},
		})))
		require.NoError(s.BuildPut(false, serverptypes.TestValidBuild(t, &pb.Build{
			Id: "3",
			Application: &pb.Ref_Application{
				Application: "B",
				Project:     "B",
			},
		})))

		// Create some other resources
		require.NoError(s.DeploymentPut(false, serverptypes.TestValidDeployment(t, &pb.Deployment{
			Id: "1",
		})))

		// Workspace list should only list one
		{
			result, err := s.ProjectListWorkspaces(ctx, &pb.Ref_Project{Project: "B"})
			require.NoError(err)
			require.Len(result, 1)
			require.NotNil(result[0].Workspace)
		}
	})
}
