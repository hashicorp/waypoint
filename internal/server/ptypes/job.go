package ptypes

import (
	"errors"
	"reflect"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestJobNew(t testing.T, src *pb.Job) *pb.Job {
	t.Helper()

	if src == nil {
		src = &pb.Job{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Job{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "w_test",
		},
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		},
		DataSource: &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Local{
				Local: &pb.Job_Local{},
			},
		},
		Operation: &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		},
	}))

	return src
}

// ValidateJob validates the job structure.
func ValidateJob(job *pb.Job) error {
	return validation.ValidateStruct(job,
		validation.Field(&job.Id, validation.By(isEmpty)),
		validation.Field(&job.Application, validation.Required),
		validation.Field(&job.Workspace, validation.Required),
		validation.Field(&job.TargetRunner, validation.Required),
		validation.Field(&job.Operation, validation.Required),
	)
}

func isEmpty(v interface{}) error {
	if reflect.ValueOf(v).IsZero() {
		return nil
	}

	return errors.New("must be empty")
}
