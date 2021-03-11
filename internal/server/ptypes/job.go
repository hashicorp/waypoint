package ptypes

import (
	"errors"
	"path/filepath"
	"reflect"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
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
	return validationext.Error(validation.ValidateStruct(job,
		validation.Field(&job.Id, validation.By(isEmpty)),
		validation.Field(&job.Application, validation.Required),
		validation.Field(&job.Workspace, validation.Required),
		validation.Field(&job.TargetRunner, validation.Required),
		validation.Field(&job.Operation, validation.Required),
	))
}

// ValidateJobDataSourceRules
func ValidateJobDataSourceRules(v *pb.Job_DataSource) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Source, validation.Required),

		validationext.StructOneof(&v.Source, (*pb.Job_DataSource_Git)(nil),
			func() []*validation.FieldRules {
				v := v.Source.(*pb.Job_DataSource_Git)
				return validateJobDataSourceGitRules(v)
			}),
	}
}

// validateJobDataSourceGitRules
func validateJobDataSourceGitRules(v *pb.Job_DataSource_Git) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Git.Url, validation.Required),
		validation.Field(&v.Git.Path, validation.By(hasNoDotDot)),

		validationext.StructOneof(&v.Git.Auth, (*pb.Job_Git_Basic_)(nil),
			func() []*validation.FieldRules {
				v := v.Git.Auth.(*pb.Job_Git_Basic_)
				return []*validation.FieldRules{
					validation.Field(&v.Basic.Username, validation.Required),
					validation.Field(&v.Basic.Password, validation.Required),
				}
			}),

		validationext.StructOneof(&v.Git.Auth, (*pb.Job_Git_Ssh)(nil),
			func() []*validation.FieldRules {
				v := v.Git.Auth.(*pb.Job_Git_Ssh)
				return []*validation.FieldRules{
					validation.Field(&v.Ssh.PrivateKeyPem,
						validation.Required, isSSHKey(v)),
				}
			}),
	}
}

func isEmpty(v interface{}) error {
	if reflect.ValueOf(v).IsZero() {
		return nil
	}

	return errors.New("must be empty")
}

// isGitSSHKey validates the SSH key given.
func isSSHKey(v *pb.Job_Git_Ssh) validation.Rule {
	return validation.By(func(_ interface{}) error {
		if len(v.Ssh.PrivateKeyPem) == 0 {
			return nil
		}

		_, err := ssh.NewPublicKeys(
			"git",
			[]byte(v.Ssh.PrivateKeyPem),
			v.Ssh.Password,
		)

		return err
	})
}

func hasNoDotDot(v interface{}) error {
	for _, part := range filepath.SplitList(v.(string)) {
		if part == ".." {
			return errors.New("must not contain '..'")
		}
	}

	return nil
}
