package ptypes

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestGetWorkspaceRequest(t testing.T, src *pb.GetWorkspaceRequest) *pb.GetWorkspaceRequest {
	t.Helper()

	if src == nil {
		src = &pb.GetWorkspaceRequest{}
	}

	require.NoError(t, mergo.Merge(src, &pb.GetWorkspaceRequest{
		Workspace: &pb.Ref_Workspace{
			Workspace: "w_test",
		},
	}))

	return src
}

// ValidateGetWorkspaceRequest
func ValidateGetWorkspaceRequest(v *pb.GetWorkspaceRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Workspace, validation.Required),
		validationext.StructField(&v.Workspace, func() []*validation.FieldRules {
			return ValidateRefWorkspaceRules(v.Workspace)
		}),
	))
}
