package statetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/jsonpb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["project_template"] = []testFunc{
		TestProjectTemplate,
		TestProjectTemplateSetAllProperties,
	}
}

func TestProjectTemplate(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("Basic CRUD operations", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// No project templates initially
		{
			resp, _, err := s.ProjectTemplateList(ctx, nil)
			require.NoError(err)
			require.Empty(resp)
		}

		// Cannot write project template with no datasource
		require.Error(s.ProjectTemplatePut(ctx, &pb.ProjectTemplate{
			Name:            "foo",
			Description:     "test",
			ProjectSettings: nil,
		}))

		// Can write a project with a datasource
		pt1 := &pb.Ref_ProjectTemplate{
			Name: "pt1",
		}
		require.NoError(s.ProjectTemplatePut(ctx, &pb.ProjectTemplate{
			Name:        pt1.Name,
			Description: "desc",
			ProjectSettings: &pb.Project{
				DataSource: &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Git{
						Git: &pb.Job_Git{
							Url: "https://github.com/hashicorp/test",
						},
					},
				},
			},
		}))

		// Can get that project template
		{
			resp, err := s.ProjectTemplateGet(ctx, pt1)
			require.NoError(err)
			require.Equal(resp.Name, pt1.Name)
		}

		// Can modify that project template

		{
			newDesc := "desc2"
			require.NoError(s.ProjectTemplatePut(ctx, &pb.ProjectTemplate{
				Name:        pt1.Name,
				Description: newDesc,
			}))

			resp, err := s.ProjectTemplateGet(ctx, pt1)
			require.NoError(err)
			require.Equal(resp.Name, pt1.Name)
			require.Equal(resp.Description, newDesc)
		}

		// Can list project templates

		{
			resp, _, err := s.ProjectTemplateList(ctx, nil)
			require.NoError(err)
			require.Len(resp, 1)
		}

		// Can delete that project template

		{
			require.NoError(s.ProjectTemplateDelete(ctx, nil))

			_, getErr := s.ProjectTemplateGet(ctx, pt1)
			require.Error(getErr)

			listResp, _, err := s.ProjectTemplateList(ctx, nil)
			require.NoError(err)
			require.Len(listResp, 0)
		}
	})
}

func TestProjectTemplateSetAllProperties(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	// A project with all the properties set, excluding those not relevant
	// to template project settings
	initialProjectSettings := serverptypes.TestProjectAllProperties(t, nil)

	// Nil out settings that aren't relevant to template project settings
	initialProjectSettings.Name = ""
	initialProjectSettings.State = 0

	name := "test_template"

	projectTemplate := &pb.ProjectTemplate{
		Name:        "test",
		Description: "test",
		SourceCodePlatform: &pb.ProjectTemplate_Github{
			Github: &pb.ProjectTemplate_SourceCodePlatformGithub{
				Source: &pb.ProjectTemplate_SourceCodePlatformGithub_Source{
					Owner: "test",
					Repo:  name,
				},
				Destination: &pb.ProjectTemplate_SourceCodePlatformGithub_Destination{
					Private:            false,
					IncludeAllBranches: false,
				},
			},
		},
		// TODO: Tokens field test
		ProjectSettings: initialProjectSettings,
	}

	initialJsonBytes, err := jsonpb.Marshal(projectTemplate)
	require.NoError(err)
	initialJsonStr := string(initialJsonBytes)

	// Set
	err = s.ProjectTemplatePut(ctx, projectTemplate)
	require.NoError(err)

	// Get
	resp, err := s.ProjectGet(ctx, &pb.Ref_Project{
		Project: projectTemplate.Name,
	})
	require.NoError(err)
	require.NotNil(resp)

	// Compare the two
	respJsonBytes, err := jsonpb.Marshal(resp)
	require.NoError(err)
	respJsonStr := string(respJsonBytes)

	require.Equal(initialJsonStr, respJsonStr)
}
