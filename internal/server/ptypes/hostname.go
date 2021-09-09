package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// ValidateCreateHostnameRequest
func ValidateCreateHostnameRequest(v *pb.CreateHostnameRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Hostname, validation.Required),
		validation.Field(&v.Target, validation.Required),
	))
}
