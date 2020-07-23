package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeAuthOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	op, ok := job.Operation.(*pb.Job_Auth)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	var results []*pb.Job_AuthResult_Result
	for _, c := range app.Components() {
		info := app.ComponentProto(c)
		if info == nil {
			// Should never happen
			continue
		}

		L := log.With("type", info.Type.String(), "name", info.Name)
		L.Debug("checking auth")

		// Start building our result. We append it right away. Since we're
		// appending a pointer we can keep modifying it.
		var result pb.Job_AuthResult_Result
		results = append(results, &result)
		result.Component = info
		result.AuthSupported = app.CanAuth(c)

		// Validate the auth
		err := app.ValidateAuth(ctx, c)
		result.CheckResult = err == nil
		if err != nil {
			st, _ := status.FromError(err)
			result.CheckError = st.Proto()
		}

		L.Debug("auth result", "result", result.CheckResult, "error", result.CheckError)

		if op.Auth.CheckOnly {
			continue
		}

		// TODO
		return nil, status.Errorf(codes.Unimplemented, "CheckOnly is required")
	}

	return &pb.Job_Result{
		Auth: &pb.Job_AuthResult{
			Results: results,
		},
	}, nil
}
