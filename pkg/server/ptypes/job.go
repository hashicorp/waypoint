// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
// TODO: This still fails if the job passed in to be validated is nil
func ValidateJob(job *pb.Job) error {
	return validationext.Error(validation.ValidateStruct(job,
		ValidateJobRules(job)...,
	))
}

// ValidateJobRules
func ValidateJobRules(job *pb.Job) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&job.Id, validation.By(isEmpty)),
		validation.Field(&job.Application, validation.Required),
		validation.Field(&job.Workspace, validation.Required),
		validationext.StructField(&job.Workspace, func() []*validation.FieldRules {
			return ValidateJobWorkspaceRules(job.Workspace)
		}),
		validation.Field(&job.TargetRunner, validation.Required),
		validation.Field(&job.Operation, validation.Required),
		validationext.StructField(&job.DataSource, func() []*validation.FieldRules {
			return ValidateJobDataSourceRules(job.DataSource)
		}),
	}
}

// ValidateJobWorkspaceRules
func ValidateJobWorkspaceRules(v *pb.Ref_Workspace) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Workspace, validation.Required),
	}
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
		validation.Field(&v.Git.Path, validation.By(hasNoDotDot), validation.By(isGitPath)),

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

// ValidateListJobsRequest
func ValidateListJobsRequest(v *pb.ListJobsRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validationext.StructField(&v.Pagination, func() []*validation.FieldRules {
			return ValidatePaginationRequestRules(v.Pagination)
		}),
	))
}

func isEmpty(v interface{}) error {
	if reflect.ValueOf(v).IsZero() {
		return nil
	}

	return errors.New("must be empty")
}

// isGitPath validates the Git path.
func isGitPath(v interface{}) error {
	path := v.(string)
	if len(path) == 0 {
		return nil
	}

	if filepath.IsAbs(path) {
		return errors.New("must be relative")
	}

	// We do this so we can just assume that all slashes are filepath.Sep
	path = filepath.ToSlash(path)

	// Verify we don't start with ./ or .\
	if len(path) >= 2 && path[0] == '.' && path[1] == filepath.Separator {
		return errors.New("relative path shouldn't start with " + path[:2])
	}

	// Verify we don't have any '//' in there. This also catches anything
	// more than 2 since any grouping of 3 or more is also a grouping of at least 2
	multisep := strings.Repeat(string(filepath.Separator), 2)
	if strings.Contains(path, multisep) {
		return errors.New("path should not contain repeated separator characters such as '//'")
	}

	// We also don't want '..' anywhere in the path, but that
	// is validated with hasNoDotDot.

	// We also want paths to end with / but that seems overly
	// pedantic so that is something we'll add ourselves in our
	// data source anytime we need the path to end with a slash.

	return nil
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
	path := v.(string)
	path = filepath.ToSlash(path)
	for _, part := range strings.Split(path, string(filepath.Separator)) {
		if part == ".." {
			return errors.New("must not contain '..'")
		}
	}

	return nil
}
