// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestProject returns a valid project for tests.
func TestProject(t testing.T, src *pb.Project) *pb.Project {
	t.Helper()

	if src == nil {
		src = &pb.Project{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Project{
		Name: "test",
	}))

	return src
}

// Type wrapper around the proto type so that we can add some methods.
type Project struct{ *pb.Project }

// App returns the index of the app with the given name or -1 if its not found.
func (p *Project) App(n string) int {
	n = strings.ToLower(n)
	for i, app := range p.Applications {
		if strings.ToLower(app.Name) == n {
			return i
		}
	}

	return -1
}

// ValidateProject validates the project structure.
func ValidateProject(p *pb.Project) error {
	return validationext.Error(validation.ValidateStruct(p,
		ValidateProjectRules(p)...,
	))
}

// ValidateProjectRules
func ValidateProjectRules(p *pb.Project) []*validation.FieldRules {
	return []*validation.FieldRules{

		validation.Field(&p.Name,
			validation.Required,
			validation.By(validatePathToken),
		),

		validation.Field(&p.WaypointHcl, isWaypointHcl(p)),

		validationext.StructField(&p.DataSource, func() []*validation.FieldRules {
			return ValidateJobDataSourceRules(p.DataSource)
		}),

		validationext.StructField(&p.DataSourcePoll, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&p.DataSourcePoll.Interval, validationext.IsDuration),
			}
		}),
	}
}

// ValidateUpsertProjectRequest
func ValidateUpsertProjectRequest(v *pb.UpsertProjectRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Project, validation.Required),
		validationext.StructField(&v.Project, func() []*validation.FieldRules {
			return ValidateProjectRules(v.Project)
		}),
	))
}

// ValidateGetProjectRequest
func ValidateGetProjectRequest(v *pb.GetProjectRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Project, validation.Required),
	))
}

// ValidateListProjectsRequest
func ValidateListProjectsRequest(v *pb.ListProjectsRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validationext.StructField(&v.Pagination, func() []*validation.FieldRules {
			return ValidatePaginationRequestRules(v.Pagination)
		}),
	))
}

// ValidateUIListProjectsRequest
func ValidateUIListProjectsRequest(v *pb.UI_ListProjectsRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validationext.StructField(&v.Pagination, func() []*validation.FieldRules {
			return ValidatePaginationRequestRules(v.Pagination)
		}),
	))
}

// ValidateDestroyProjectRequest
func ValidateDestroyProjectRequest(v *pb.DestroyProjectRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Project, validation.Required),
	))
}

// ValidateUIGetProjectRequest
func ValidateUIGetProjectRequest(v *pb.UI_GetProjectRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Project, validation.Required),
	))
}

func isWaypointHcl(p *pb.Project) validation.Rule {
	return validation.By(func(_ interface{}) error {
		if len(p.WaypointHcl) == 0 {
			return nil
		}

		switch p.WaypointHclFormat {
		case pb.Hcl_HCL:
			_, diag := hclsyntax.ParseConfig(p.WaypointHcl, "<waypoint-hcl>", hcl.Pos{})
			if diag.HasErrors() {
				return diag
			}

			return nil
		case pb.Hcl_JSON:
			_, diag := hcljson.Parse(p.WaypointHcl, "<waypoint-hcl>")
			if diag.HasErrors() {
				return diag
			}

			return nil
		default:
			return errors.New("unknown format")
		}
	})
}
