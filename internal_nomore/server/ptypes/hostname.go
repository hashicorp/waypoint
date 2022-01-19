package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/waypoint/internal_nomore/pkg/validationext"
	pb "github.com/hashicorp/waypoint/internal_nomore/server/gen"
)

// ValidateCreateHostnameRequest
func ValidateCreateHostnameRequest(v *pb.CreateHostnameRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Target, validation.Required),
	))
}
