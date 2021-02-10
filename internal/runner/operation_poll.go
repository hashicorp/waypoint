package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executePollOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
) (*pb.Job_Result, error) {
	// TODO
	return &pb.Job_Result{}, nil
}
