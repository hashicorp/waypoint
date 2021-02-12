package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
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

	cs, err := app.Components(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range cs {
		defer c.Close()
	}

	var results []*pb.Job_AuthResult_Result
	for _, c := range cs {
		info := c.Info
		if info == nil {
			// Should never happen
			continue
		}

		// If we have a ref set for a component then we only auth ones that match.
		if op.Auth.Component != nil {
			ptypeC := serverptypes.Component{Component: info}
			if !ptypeC.Match(op.Auth.Component) {
				continue
			}
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

		// If we authed successfully or we're only checking, we're done.
		if result.CheckResult || op.Auth.CheckOnly {
			continue
		}

		// Attempt to authenticate
		L.Trace("attempting auth")
		authResult, err := app.Auth(ctx, c)
		if err != nil {
			st, _ := status.FromError(err)
			result.AuthError = st.Proto()
		}
		if authResult != nil {
			result.AuthCompleted = authResult.Authenticated
		}

		// If we did complete the auth, revalidate it.
		if result.AuthCompleted {
			err := app.ValidateAuth(ctx, c)
			result.CheckResult = err == nil
			if err != nil {
				st, _ := status.FromError(err)
				result.CheckError = st.Proto()
			}
		}
	}

	// If we referenced a component and have no results, then that component
	// wasn't found and this is an error.
	if op.Auth.Component != nil && len(results) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition,
			"component to auth was not found for this app")
	}

	return &pb.Job_Result{
		Auth: &pb.Job_AuthResult{
			Results: results,
		},
	}, nil
}
