package utils

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func RegisterTaskDefinition(def *ecs.RegisterTaskDefinitionInput, ecsSvc *ecs.ECS) (*ecs.TaskDefinition, error) {
	// AWS is eventually consistent so even though we probably created the
	// resources that are referenced by the task definition, it can error out if
	// we try to reference those resources too quickly. So we're forced to guard
	// actions which reference other AWS services with loops like this.
	var taskOut *ecs.RegisterTaskDefinitionOutput
	var err error
	for i := 0; i < 30; i++ {
		taskOut, err = ecsSvc.RegisterTaskDefinition(def)
		if err == nil {
			break
		}

		// if we encounter an unrecoverable error, exit now.
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "ResourceConflictException" || aerr.Code() == "ClientException" {
				return nil, err
			}
		}

		// otherwise sleep and try again
		time.Sleep(2 * time.Second)
	}

	// the above loop could expire and never get a valid task definition, so
	// guard against a nil taskOut here
	if taskOut == nil {
		return nil, fmt.Errorf("error registering task definition, last error: %w", err)
	}

	return taskOut.TaskDefinition, nil
}
