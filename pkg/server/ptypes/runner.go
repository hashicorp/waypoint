package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	serverpkg "github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestRunner(t testing.T, src *pb.Runner) *pb.Runner {
	t.Helper()

	if src == nil {
		src = &pb.Runner{}
	}

	id, err := serverpkg.Id()
	require.NoError(t, err)

	require.NoError(t, mergo.Merge(src, &pb.Runner{
		Id: id,
	}))

	return src
}

// ValidateAdoptRunnerRequest
func ValidateAdoptRunnerRequest(v *pb.AdoptRunnerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.RunnerId, validation.Required),
	))
}

// ValidateForgetRunnerRequest
func ValidateForgetRunnerRequest(v *pb.ForgetRunnerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.RunnerId, validation.Required),
	))
}
