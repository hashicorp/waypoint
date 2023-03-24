package ptypes

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

// validatePathToken Validates a field that can be used as a GRPC Gateway path token
// Check the route in gateway.yml to see which fields are path tokens.
func validatePathToken(pathToken interface{}) error {
	s, _ := pathToken.(string)

	// A path token cannot contain ../, as grpc gateway will interpret that
	// as a path traversal.
	if strings.Contains(s, "../") {
		return errors.New("name cannot contain '../'")
	} else if !regexp.MustCompile(`^[a-zA-Z0-9_-]*$`).MatchString(s) {
		return fmt.Errorf("name %q must be alpha numeric", s)
	}
	return nil
}
