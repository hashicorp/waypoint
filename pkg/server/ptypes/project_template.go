package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const (
	PROJECT_TEMPLATE_ID_LENGTH               = 64
	PROJECT_TEMPLATE_NAME_LENGTH             = 64
	PROJECT_TEMPLATE_TAG_LENGTH              = 64
	PROJECT_TEMPLATE_SUMMARY_LENGTH          = 64
	PROJECT_TEMPLATE_EXPANDED_SUMMARY_LENGTH = 512
	PROJECT_TEMPLATE_README_LENGTH           = 1024 ^ 2
	PROJECT_TEMPLATE_WAYPOINT_HCL_LENGTH     = 1024 ^ 2

	TERRAFORM_NOCODE_MODULE_SOURCE_LENGTH  = 4096
	TERRAFORM_NOCODE_MODULE_VERSION_LENGTH = 1024
)

func TestProjectTemplate(t testing.T, src *pb.ProjectTemplate) *pb.ProjectTemplate {
	t.Helper()

	if src == nil {
		src = &pb.ProjectTemplate{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Project{
		Name: "test",
	}))

	return src
}

func ValidateCreateProjectTemplateRequest(req *pb.CreateProjectTemplateRequest) error {
	return validationext.Error(validation.ValidateStruct(req,
		validation.Field(&req.ProjectTemplate, validation.Required),
		validationext.StructField(&req.ProjectTemplate, func() []*validation.FieldRules {
			return append(
				// Rules specific to creating a project template
				[]*validation.FieldRules{
					validation.Field(&req.ProjectTemplate.Name, validation.Required),
				},

				// General project template validation rules
				validateProjectTemplateRules(req.ProjectTemplate)...,
			)
		}),
	))
}

func ValidateUpdateProjectTemplateRequest(req *pb.UpdateProjectTemplateRequest) error {
	return validationext.Error(validation.ValidateStruct(req,
		validation.Field(&req.ProjectTemplate, validation.Required),
		validationext.StructField(&req.ProjectTemplate, func() []*validation.FieldRules {
			return append(
				// Rules specific to creating a project template
				[]*validation.FieldRules{
					// Require either Name or ID
					validation.Field(&req.ProjectTemplate.Id, validation.Required.When(req.ProjectTemplate.Name == "").Error("Either Name or ID is required.")),
					validation.Field(&req.ProjectTemplate.Name, validation.Required.When(req.ProjectTemplate.Id == "").Error("Either Name or ID is required.")),
				},

				// General project template validation rules
				validateProjectTemplateRules(req.ProjectTemplate)...,
			)
		}),
	))
}

// validateProjectTemplateRules validates the rules that must be true of any project template in any
// request or response.
func validateProjectTemplateRules(pt *pb.ProjectTemplate) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&pt.Id, validation.Length(0, PROJECT_TEMPLATE_ID_LENGTH)),
		validation.Field(&pt.Name, validation.Length(0, PROJECT_TEMPLATE_NAME_LENGTH)),

		validationext.StructField(&pt.TerraformNocodeModule, func() []*validation.FieldRules {
			return validateTerraformNocodeModule(pt.TerraformNocodeModule)
		}),

		validation.Field(&pt.Summary, validation.Length(0, PROJECT_TEMPLATE_SUMMARY_LENGTH)),
		validation.Field(&pt.ExpandedSummary, validation.Length(0, PROJECT_TEMPLATE_EXPANDED_SUMMARY_LENGTH)),
		validation.Field(&pt.ReadmeMarkdownTemplate, validation.Length(0, PROJECT_TEMPLATE_README_LENGTH)),

		validation.Field(&pt.Tags, validation.Each(
			validation.Length(1, PROJECT_TEMPLATE_TAG_LENGTH),
		)),

		validationext.StructField(&pt.WaypointProject, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(
					&pt.WaypointProject.WaypointHclTemplate,
					validation.Length(1, PROJECT_TEMPLATE_WAYPOINT_HCL_LENGTH),
				),
			}
		}),
	}
}

func validateTerraformNocodeModule(t *pb.ProjectTemplate_TerraformNocodeModule) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&t.Source,
			validation.Required,
			validation.Length(1, TERRAFORM_NOCODE_MODULE_SOURCE_LENGTH),
		),
		validation.Field(&t.Version,
			validation.Required,
			validation.Length(1, TERRAFORM_NOCODE_MODULE_VERSION_LENGTH),
		),
	}
}
