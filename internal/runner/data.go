// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/datasource"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
	ui terminal.UI,
	source *pb.Job_DataSource,
	overrides map[string]string,
) (string, *pb.Job_DataSource_Ref, func() error, error) {
	sourcer, err := r.dataSourcer(ctx, log, source, overrides)
	if err != nil {
		return "", nil, nil, err
	}

	// Get data
	return sourcer.Get(ctx, log, ui, source, r.tempDir)
}

// dataSourcer returns the datasource implementation for a job,
// properly configured with any overrides.
func (r *Runner) dataSourcer(
	ctx context.Context,
	log hclog.Logger,
	source *pb.Job_DataSource,
	overrides map[string]string,
) (datasource.Sourcer, error) {
	if source == nil {
		return nil, status.Errorf(codes.Internal,
			"data source not set for job")
	}

	// Determine our sourcer
	typ := reflect.TypeOf(source.Source)
	factory, ok := datasource.FromType[typ]
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition,
			"invalid data source type: %s", typ.String())
	}
	sourcer := factory()

	// Apply any overrides
	if len(overrides) > 0 {
		if err := sourcer.Override(source, overrides); err != nil {
			return nil, status.Errorf(codes.FailedPrecondition,
				"error with data source overrides: %s", err)
		}
	}

	return sourcer, nil
}
