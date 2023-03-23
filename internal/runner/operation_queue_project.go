// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TODO: test. The reason there aren't unit tests at the moment is because
// at the time of writing there is no way to mock a plugin so we can't actually
// run any config that has a 'use' stanza (and that is required).
func (r *Runner) executeQueueProjectOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	op, ok := job.Operation.(*pb.Job_QueueProject)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}
	jobTemplate := op.QueueProject.JobTemplate

	// Go through each app in the project and queue a job.
	var queueApps []*pb.Job_QueueProjectResult_Application
	log.Debug("total number of apps in project", "len", len(project.Apps()))
	for _, name := range project.Apps() {
		log.Debug("queueing job", "app", name)

		app, err := project.App(name)
		if err != nil {
			return nil, err
		}

		// Build up our concrete job from the template.
		job := proto.Clone(jobTemplate).(*pb.Job)
		job.Application = app.Ref()

		resp, err := r.client.QueueJob(ctx, &pb.QueueJobRequest{
			Job: job,
		})
		if err != nil {
			log.Warn("error queueing job", "app", name, "err", err)
			return nil, err
		}

		queueApps = append(queueApps, &pb.Job_QueueProjectResult_Application{
			Application: job.Application,
			JobId:       resp.JobId,
		})
	}

	return &pb.Job_Result{
		QueueProject: &pb.Job_QueueProjectResult{
			Applications: queueApps,
		},
	}, nil
}
