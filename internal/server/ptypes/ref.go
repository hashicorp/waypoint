package ptypes

import (
	"github.com/go-ozzo/ozzo-validation/v4"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// ValidateRefWorkspaceRules
func ValidateRefWorkspaceRules(v *pb.Ref_Workspace) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Workspace, validation.Required),
	}
}
