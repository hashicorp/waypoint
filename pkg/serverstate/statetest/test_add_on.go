package statetest

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/pagination"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["add_on"] = []testFunc{
		TestAddOnFeatures,
		TestAddOnPagination,
	}
}

func TestAddOnFeatures(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	pn := "testProject"
	proj := serverptypes.TestProject(t, &pb.Project{
		Name: pn,
	})
	err := s.ProjectPut(ctx, proj)
	require.NoError(err)

	readme := []byte(strings.TrimSpace(`
My favorite add-on README.
`))

	tags := []string{
		"varset1",
		"varset2",
	}

	addOnDefinitionName := "postgres"
	testAddOnDefinition := &pb.AddOnDefinition{
		Name: addOnDefinitionName,
		TerraformNocodeModule: &pb.TerraformNocodeModule{
			Source:  "my/test/module",
			Version: "0.0.1",
		},
		ShortSummary:           "My short summary.",
		LongSummary:            "My very long summary.",
		ReadmeMarkdownTemplate: readme,
		Tags: []string{
			"tag",
			"you're",
			"it",
		},
		TfVariableSetIds: tags,
	}

	t.Run("Create, get, and delete Add-On definition", func(t *testing.T) {
		// Create
		addOnDefinition, err := s.AddOnDefinitionPut(ctx, testAddOnDefinition)
		require.NoError(err)
		require.Equal(testAddOnDefinition.Name, addOnDefinition.Name)
		require.NotNil(addOnDefinition.TerraformNocodeModule)
		require.Equal(testAddOnDefinition.TerraformNocodeModule.Source, addOnDefinition.TerraformNocodeModule.Source)
		require.Equal(testAddOnDefinition.TerraformNocodeModule.Version, addOnDefinition.TerraformNocodeModule.Version)
		require.Equal(testAddOnDefinition.Tags, addOnDefinition.Tags)
		require.Equal(testAddOnDefinition.TfVariableSetIds, addOnDefinition.TfVariableSetIds)
		require.Equal(testAddOnDefinition.ReadmeMarkdownTemplate, addOnDefinition.ReadmeMarkdownTemplate)
		require.Equal(testAddOnDefinition.ShortSummary, addOnDefinition.ShortSummary)
		require.Equal(testAddOnDefinition.LongSummary, addOnDefinition.LongSummary)

		// Get by ID
		actualAddOnDefinition, err := s.AddOnDefinitionGet(ctx, &pb.Ref_AddOnDefinition{
			Identifier: &pb.Ref_AddOnDefinition_Id{
				Id: addOnDefinition.Id,
			},
		})
		require.NoError(err)
		require.Equal(testAddOnDefinition.Name, actualAddOnDefinition.Name)
		require.NotNil(actualAddOnDefinition.TerraformNocodeModule)
		require.Equal(testAddOnDefinition.TerraformNocodeModule.Source, actualAddOnDefinition.TerraformNocodeModule.Source)
		require.Equal(testAddOnDefinition.TerraformNocodeModule.Version, actualAddOnDefinition.TerraformNocodeModule.Version)
		require.Equal(testAddOnDefinition.Tags, actualAddOnDefinition.Tags)
		require.Equal(testAddOnDefinition.TfVariableSetIds, actualAddOnDefinition.TfVariableSetIds)
		require.Equal(testAddOnDefinition.ReadmeMarkdownTemplate, actualAddOnDefinition.ReadmeMarkdownTemplate)
		require.Equal(testAddOnDefinition.ShortSummary, actualAddOnDefinition.ShortSummary)
		require.Equal(testAddOnDefinition.LongSummary, actualAddOnDefinition.LongSummary)

		// Delete Add-On definition
		err = s.AddOnDefinitionDelete(ctx, &pb.Ref_AddOnDefinition{
			Identifier: &pb.Ref_AddOnDefinition_Name{
				Name: testAddOnDefinition.Name,
			},
		})
		require.NoError(err)

		// Verify Add-On definition is deleted
		daod, err := s.AddOnDefinitionGet(ctx, &pb.Ref_AddOnDefinition{
			Identifier: &pb.Ref_AddOnDefinition_Name{
				Name: testAddOnDefinition.Name,
			},
		})
		// expecting a not found error
		require.Error(err)
		require.Nil(daod)

	})

	testUpdatedTestAddOnDefinition := &pb.AddOnDefinition{
		Name: "new-postgres", // new name
		TerraformNocodeModule: &pb.TerraformNocodeModule{
			Source:  "my/test/module",
			Version: "0.0.2",
		},
		ShortSummary:           "My super short summary.",
		LongSummary:            "My super long summary.",
		ReadmeMarkdownTemplate: readme,
		Tags: []string{
			"gotcha",
		},
		TfVariableSetIds: tags,
	}

	t.Run("Update Add-on Definition & get it by the new name", func(t *testing.T) {
		// Create an add-on definition
		aod, err := s.AddOnDefinitionPut(ctx, testAddOnDefinition)
		require.NoError(err)
		require.Equal(testAddOnDefinition.Name, aod.Name)

		// Update it
		updatedAddOnDefinition, err := s.AddOnDefinitionUpdate(ctx, testUpdatedTestAddOnDefinition, &pb.Ref_AddOnDefinition{
			Identifier: &pb.Ref_AddOnDefinition_Name{
				Name: testAddOnDefinition.Name,
			},
		})
		require.NoError(err)
		require.Equal(testUpdatedTestAddOnDefinition.Name, updatedAddOnDefinition.Name)
		require.NotNil(updatedAddOnDefinition.TerraformNocodeModule)
		require.Equal(testUpdatedTestAddOnDefinition.TerraformNocodeModule.Source, updatedAddOnDefinition.TerraformNocodeModule.Source)
		require.Equal(testUpdatedTestAddOnDefinition.TerraformNocodeModule.Version, updatedAddOnDefinition.TerraformNocodeModule.Version)
		require.Equal(testUpdatedTestAddOnDefinition.Tags, updatedAddOnDefinition.Tags)
		require.Equal(testUpdatedTestAddOnDefinition.TfVariableSetIds, updatedAddOnDefinition.TfVariableSetIds)
		require.Equal(testUpdatedTestAddOnDefinition.ReadmeMarkdownTemplate, updatedAddOnDefinition.ReadmeMarkdownTemplate)
		require.Equal(testUpdatedTestAddOnDefinition.ShortSummary, updatedAddOnDefinition.ShortSummary)
		require.Equal(testUpdatedTestAddOnDefinition.LongSummary, updatedAddOnDefinition.LongSummary)

		actualAddOnDefinition, err := s.AddOnDefinitionGet(ctx, &pb.Ref_AddOnDefinition{
			Identifier: &pb.Ref_AddOnDefinition_Name{
				Name: testUpdatedTestAddOnDefinition.Name,
			},
		})
		require.NoError(err)
		require.Equal(testUpdatedTestAddOnDefinition.Name, actualAddOnDefinition.Name)
		require.NotNil(updatedAddOnDefinition.TerraformNocodeModule)
		require.Equal(testUpdatedTestAddOnDefinition.TerraformNocodeModule.Source, actualAddOnDefinition.TerraformNocodeModule.Source)
		require.Equal(testUpdatedTestAddOnDefinition.TerraformNocodeModule.Version, actualAddOnDefinition.TerraformNocodeModule.Version)
		require.Equal(testUpdatedTestAddOnDefinition.Tags, actualAddOnDefinition.Tags)
		require.Equal(testUpdatedTestAddOnDefinition.TfVariableSetIds, actualAddOnDefinition.TfVariableSetIds)
		require.Equal(testUpdatedTestAddOnDefinition.ReadmeMarkdownTemplate, actualAddOnDefinition.ReadmeMarkdownTemplate)
		require.Equal(testUpdatedTestAddOnDefinition.ShortSummary, actualAddOnDefinition.ShortSummary)
		require.Equal(testUpdatedTestAddOnDefinition.LongSummary, actualAddOnDefinition.LongSummary)
	})

	testAddOn := &pb.AddOn{
		Name: "your friendly neighborhood add-on",
		Project: &pb.Ref_Project{
			Project: pn,
		},
		Definition: &pb.Ref_AddOnDefinition{
			Identifier: &pb.Ref_AddOnDefinition_Name{
				Name: testUpdatedTestAddOnDefinition.Name,
			},
		},
		ShortSummary: "My super short summary.",
		LongSummary:  "My super long summary.",
		TerraformNocodeModule: &pb.TerraformNocodeModule{
			Source:  "my/test/module",
			Version: "0.0.2",
		},
		ReadmeMarkdown: readme, // this does NOT test any rendering
		Tags:           tags,
		CreatedBy:      "foo@bar.com",
	}

	t.Run("Create, get, update, and delete Add-on", func(t *testing.T) {
		// Create an add-on definition
		addOnDefinition, err := s.AddOnDefinitionPut(ctx, testAddOnDefinition)
		require.NoError(err)
		require.Equal(testAddOnDefinition.Name, addOnDefinition.Name)

		// Create an add-on using the definition
		addOn, err := s.AddOnPut(ctx, testAddOn)
		require.NoError(err)
		require.Equal(testAddOn.Name, addOn.Name)

		actualAddOn, err := s.AddOnGet(ctx, &pb.Ref_AddOn{
			Identifier: &pb.Ref_AddOn_Name{
				Name: testAddOn.Name,
			},
		})
		require.NoError(err)
		require.Equal(testAddOn.Name, actualAddOn.Name)
		require.Equal(testAddOn.Tags, actualAddOn.Tags)
		require.Equal(testAddOn.Project, actualAddOn.Project)
		require.Equal(testAddOn.Definition, actualAddOn.Definition)
		require.Equal(testAddOn.ReadmeMarkdown, actualAddOn.ReadmeMarkdown)
		require.Equal(testAddOn.ShortSummary, actualAddOn.ShortSummary)
		require.Equal(testAddOn.LongSummary, actualAddOn.LongSummary)
		require.NotNil(actualAddOn.TerraformNocodeModule)
		require.Equal(testAddOn.TerraformNocodeModule.Source, actualAddOn.TerraformNocodeModule.Source)
		require.Equal(testAddOn.TerraformNocodeModule.Version, actualAddOn.TerraformNocodeModule.Version)

		updatedAddOnName := "your updated friendly neighborhood add-on"
		updatedAddOn, err := s.AddOnUpdate(ctx,
			&pb.AddOn{
				Name: updatedAddOnName,
			},
			&pb.Ref_AddOn{
				Identifier: &pb.Ref_AddOn_Name{
					Name: testAddOn.Name,
				},
			},
		)
		require.NoError(err)
		require.NotNil(updatedAddOn)
		require.Equal(updatedAddOnName, updatedAddOn.Name)

		err = s.AddOnDelete(ctx, &pb.Ref_AddOn{
			Identifier: &pb.Ref_AddOn_Name{
				Name: updatedAddOnName,
			},
		})
		require.NoError(err)

		// Verify Add-On is deleted
		actualAddOn, err = s.AddOnGet(ctx, &pb.Ref_AddOn{
			Identifier: &pb.Ref_AddOn_Name{
				Name: updatedAddOnName,
			},
		})
		// expecting a not found error
		require.Error(err)
		require.Nil(actualAddOn)
	})

	t.Run("Get Add-On by Id with Add-On definition deleted", func(t *testing.T) {
		// Create an add-on using the definition
		addOn, err := s.AddOnPut(ctx, testAddOn)
		require.NoError(err)
		require.Equal(testAddOn.Name, addOn.Name)
		require.Equal(testAddOn.Tags, addOn.Tags)
		require.Equal(testAddOn.Project, addOn.Project)
		require.Equal(testAddOn.Definition, addOn.Definition)
		require.Equal(testAddOn.ReadmeMarkdown, addOn.ReadmeMarkdown)

		// Delete Add-On definition
		err = s.AddOnDefinitionDelete(ctx, &pb.Ref_AddOnDefinition{
			Identifier: &pb.Ref_AddOnDefinition_Name{
				Name: testAddOnDefinition.Name,
			},
		})
		require.NoError(err)

		actualAddOn, err := s.AddOnGet(ctx, &pb.Ref_AddOn{
			Identifier: &pb.Ref_AddOn_Name{
				Name: testAddOn.Name,
			},
		})
		require.NoError(err)
		require.Equal(testAddOn.Name, actualAddOn.Name)
		require.Equal(testAddOn.Tags, actualAddOn.Tags)
		require.Equal(testAddOn.Project, actualAddOn.Project)
		require.Equal(testAddOn.Definition, actualAddOn.Definition)
		require.Equal(testAddOn.ReadmeMarkdown, actualAddOn.ReadmeMarkdown)
	})
}

func TestAddOnPagination(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	t.Run("List Add-On definitions", func(t *testing.T) {
		startChar := 'a'
		endChar := 'm'
		//addOnDefinitionsCount := endChar - startChar + 1
		var chars []string
		// Generate randomized add-on definitions
		for char := startChar; char <= endChar; char++ {
			chars = append(chars, fmt.Sprintf("%c", char))
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(chars), func(i, j int) {
			chars[i], chars[j] = chars[j], chars[i]
		})
		for _, char := range chars {
			aod, err := s.AddOnDefinitionPut(ctx, &pb.AddOnDefinition{
				Id:   char,
				Name: char,
			})
			require.NoError(err)
			require.NotNil(aod)
		}

		// list a-e
		aods, pr, err := s.AddOnDefinitionList(ctx, &pb.ListAddOnDefinitionsRequest{
			Pagination: serverptypes.TestPaginationRequest(t, &pb.PaginationRequest{
				PageSize: 5,
			}),
		})
		require.NoError(err)
		require.Len(aods, 5)
		expectedPageToken, _ := pagination.EncodeAndSerializePageToken("name", "e")
		require.Equal(expectedPageToken, pr.NextPageToken)
		require.Empty(pr.PreviousPageToken)

		// list f-j
		aods, pr, err = s.AddOnDefinitionList(
			ctx,
			&pb.ListAddOnDefinitionsRequest{
				Pagination: serverptypes.TestPaginationRequest(t, &pb.PaginationRequest{
					PageSize:      5,
					NextPageToken: pr.NextPageToken,
				}),
			},
		)
		require.NoError(err)
		require.Len(aods, 5)
		expectedPrevPageToken, _ := pagination.EncodeAndSerializePageToken("name", "f")
		require.Equal(expectedPrevPageToken, pr.PreviousPageToken)
		expectedNextPageToken, _ := pagination.EncodeAndSerializePageToken("name", "j")
		require.Equal(expectedNextPageToken, pr.NextPageToken)

		// list k-m
		aods, pr, err = s.AddOnDefinitionList(
			ctx,
			&pb.ListAddOnDefinitionsRequest{
				Pagination: serverptypes.TestPaginationRequest(t, &pb.PaginationRequest{
					PageSize:      5,
					NextPageToken: pr.NextPageToken,
				}),
			},
		)
		require.NoError(err)
		require.Len(aods, 3)
		expectedPrevPageToken, _ = pagination.EncodeAndSerializePageToken("name", "k")
		require.Equal(expectedPrevPageToken, pr.PreviousPageToken)
		require.Empty(pr.NextPageToken)
	})

	t.Run("List Add-Ons", func(t *testing.T) {
		pn := "my-test-project"
		proj := serverptypes.TestProject(t, &pb.Project{
			Name: pn,
		})
		err := s.ProjectPut(ctx, proj)
		require.NoError(err)

		startChar := 'a'
		endChar := 'm'
		//addOnDefinitionsCount := endChar - startChar + 1
		var chars []string
		// Generate randomized add-on definitions
		for char := startChar; char <= endChar; char++ {
			chars = append(chars, fmt.Sprintf("%c", char))
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(chars), func(i, j int) {
			chars[i], chars[j] = chars[j], chars[i]
		})
		for _, char := range chars {
			addOn, err := s.AddOnPut(ctx, &pb.AddOn{
				Name: char,
				// "a" - "m" were already created in the previous test, so just using definition "a" for all
				Definition: &pb.Ref_AddOnDefinition{Identifier: &pb.Ref_AddOnDefinition_Name{Name: "a"}},
				Project:    &pb.Ref_Project{Project: pn},
			})
			require.NoError(err)
			require.NotNil(addOn)
		}

		// list a-e
		addOns, pr, err := s.AddOnList(ctx, &pb.ListAddOnsRequest{
			Pagination: serverptypes.TestPaginationRequest(t, &pb.PaginationRequest{
				PageSize: 5,
			}),
		})
		require.NoError(err)
		require.Len(addOns, 5)
		expectedPageToken, _ := pagination.EncodeAndSerializePageToken("name", "e")
		require.Equal(expectedPageToken, pr.NextPageToken)
		require.Empty(pr.PreviousPageToken)

		// list f-j
		addOns, pr, err = s.AddOnList(
			ctx,
			&pb.ListAddOnsRequest{
				Pagination: serverptypes.TestPaginationRequest(t, &pb.PaginationRequest{
					PageSize:      5,
					NextPageToken: pr.NextPageToken,
				}),
			},
		)
		require.NoError(err)
		require.Len(addOns, 5)
		expectedPrevPageToken, _ := pagination.EncodeAndSerializePageToken("name", "f")
		require.Equal(expectedPrevPageToken, pr.PreviousPageToken)
		expectedNextPageToken, _ := pagination.EncodeAndSerializePageToken("name", "j")
		require.Equal(expectedNextPageToken, pr.NextPageToken)

		// list k-m
		addOns, pr, err = s.AddOnList(
			ctx,
			&pb.ListAddOnsRequest{
				Pagination: serverptypes.TestPaginationRequest(t, &pb.PaginationRequest{
					PageSize:      5,
					NextPageToken: pr.NextPageToken,
				}),
			},
		)
		require.NoError(err)
		require.Len(addOns, 3)
		expectedPrevPageToken, _ = pagination.EncodeAndSerializePageToken("name", "k")
		require.Equal(expectedPrevPageToken, pr.PreviousPageToken)
		require.Empty(pr.NextPageToken)
	})
}
