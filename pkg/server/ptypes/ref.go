package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ValidateRefWorkspaceRules
func ValidateRefWorkspaceRules(v *pb.Ref_Workspace) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Workspace, validation.Required),
	}
}

// ValidateRefOperationRules
func ValidateRefOperationRules(v *pb.Ref_Operation) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Target, validation.Required),
	}
}
