package runner

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// downloadJobData takes the data source of the given job, gets the data,
// and returns the directory where the data is stored.
//
// This will also return a closer function that should be deferred to
// clean up any resources created by this. Note that the directory isn't
// always a temporary directory (such as for local data) so callers should
// NOT assume this and delete data. Use the returned closer.
func (r *Runner) downloadJobData(
	ctx context.Context,
	log hclog.Logger,
	source *pb.Job_DataSource,
) (string, func() error, error) {
	switch s := source.Source.(type) {
	case *pb.Job_DataSource_Local:
		// For local data, we just return empty.
		return "", nil, nil

	case *pb.Job_DataSource_Git:
		return r.downloadJobData_git(ctx, log, s)

	default:
		return "", nil, status.Errorf(codes.FailedPrecondition,
			"unknown job data source type: %T", source.Source)
	}
}

func (r *Runner) downloadJobData_git(
	ctx context.Context,
	log hclog.Logger,
	source *pb.Job_DataSource_Git,
) (string, func() error, error) {
	// Create a temporary directory where we will store the cloned data.
	td, err := ioutil.TempDir(r.tempDir, "waypoint")
	if err != nil {
		return "", nil, err
	}
	closer := func() error {
		return os.RemoveAll(td)
	}

	// Clone
	var output bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "clone", source.Git.Url, td)
	cmd.Stdout = &output
	cmd.Stderr = &output
	cmd.Stdin = nil
	if err := cmd.Run(); err != nil {
		closer()
		return "", nil, status.Errorf(codes.Aborted,
			"Git clone failed: %s", output.String())
	}

	// Checkout if we have a ref. If we don't have a ref we use the
	// default of whatever we got.
	if ref := source.Git.Ref; ref != "" {
		output.Reset()
		cmd := exec.CommandContext(ctx, "git", "checkout", ref)
		cmd.Dir = td
		cmd.Stdout = &output
		cmd.Stderr = &output
		cmd.Stdin = nil
		if err := cmd.Run(); err != nil {
			closer()
			return "", nil, status.Errorf(codes.Aborted,
				"Git checkout failed: %s", output.String())
		}
	}

	return td, nil, nil
}
