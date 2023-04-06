// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"context"
	"github.com/hashicorp/go-hclog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeInitOp(
	ctx context.Context,
	log hclog.Logger,
	project *core.Project,
) (*pb.Job_Result, error) {
	client := project.Client()
	// UpsertWorkspace called to ensure the workspace exists before the
	// project/app is initialized.
	resp, err := client.UpsertWorkspace(ctx, &pb.UpsertWorkspaceRequest{
		Workspace: &pb.Workspace{
			Name: project.WorkspaceRef().Workspace,
		},
	})
	if err != nil {
		return nil, err
	}

	// this is unlikely to happen with a nil error above, but added here to be
	// defensive.
	if resp.Workspace == nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"unable to verify workspace (%s) exists", project.WorkspaceRef().Workspace)
	}

	// This operation upserts apps defined in the project’s waypoint.hcl
	// into the server’s database. This is important for projects that use
	// the GitOps flow without polling, as otherwise the project appears
	// empty and a manual CLI init step is required.
	for _, name := range project.Apps() {
		_, err := client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
			Project: project.Ref(),
			Name:    name,
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.Job_Result{
		Init: &pb.Job_InitResult{},
	}, nil
}
