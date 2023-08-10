// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestTask returns a valid user for tests.
func TestTask(t testing.T, src *pb.Task) *pb.Task {
	t.Helper()

	if src == nil {
		src = &pb.Task{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Task{
		Id: "test",
	}))

	return src
}

// ValidateTask validates the user structure.
func ValidateTask(v *pb.Task) error {
	return validationext.Error(validation.ValidateStruct(v,
		ValidateTaskRules(v)...,
	))
}

// ValidateTaskRules
func ValidateTaskRules(v *pb.Task) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&v.Id, validation.Required),
		validation.Field(&v.TaskJob, validation.Required),

		validationext.StructField(&v.TaskJob, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Id, validation.Required),
			}
		}),
	}
}

// ValidateUpsertTaskRequest
func ValidateUpsertTaskRequest(v *pb.UpsertTaskRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Task, validation.Required),
		validationext.StructField(&v.Task, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Task.TaskJob, validation.Required),
			}
		}),
	))
}

// ValidateGetTaskRequest
func ValidateGetTaskRequest(v *pb.GetTaskRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
	))
}

// ValidateCancelTaskRequest
func ValidateCancelTaskRequest(v *pb.CancelTaskRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
	))
}

// ValidateDeleteTaskRequest
func ValidateDeleteTaskRequest(v *pb.DeleteTaskRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
	))
}

// ValidateRefTaskRules
func ValidateRefTaskRules(v *pb.Ref_Task) []*validation.FieldRules {
	switch rv := v.Ref.(type) {
	case *pb.Ref_Task_Id:
		return []*validation.FieldRules{
			validation.Field(&rv.Id, validation.Required),
		}
	case *pb.Ref_Task_JobId:
		return []*validation.FieldRules{
			validation.Field(&rv.JobId, validation.Required),
		}
	default:
		return []*validation.FieldRules{}
	}
}
