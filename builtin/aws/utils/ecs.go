// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/go-hclog"
)

func RegisterTaskDefinition(
	def *ecs.RegisterTaskDefinitionInput,
	ecsSvc *ecs.ECS,
	log hclog.Logger,
) (*ecs.TaskDefinition, error) {
	// AWS is eventually consistent so even though we probably created the
	// resources that are referenced by the task definition, it can error out if
	// we try to reference those resources too quickly. So we're forced to guard
	// actions which reference other AWS services with loops like this.
	var taskOut *ecs.RegisterTaskDefinitionOutput
	var err error
	for i := 0; i < 30; i++ {
		taskOut, err = ecsSvc.RegisterTaskDefinition(def)
		if err != nil {
			return nil, err
		}

		if taskOut != nil && taskOut.TaskDefinition != nil {
			break
		}

		log.Debug("error registering task definition, retrying", "error", err)
		// otherwise sleep and try again
		time.Sleep(2 * time.Second)
	}

	// The above loop could expire and never get a valid task definition, so
	// guard against a nil taskOut here. It's possible that the
	// response from RegisterTaskDefinition returns an error and a non-nil
	// RegisterTaskDefinitionOutput struct, so we need to verify that both the
	// output struct and its included TaskDefinition are both non-nil before
	// assuming success.
	if taskOut == nil || taskOut.TaskDefinition == nil {
		return nil, fmt.Errorf("error registering task definition, last error: %w", err)
	}
	return taskOut.TaskDefinition, nil
}
