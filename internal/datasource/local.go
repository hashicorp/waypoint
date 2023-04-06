// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type LocalSource struct{}

func newLocalSource() Sourcer { return &LocalSource{} }

func (s *LocalSource) ProjectSource(body hcl.Body, ctx *hcl.EvalContext) (*pb.Job_DataSource, error) {
	// Return the data source
	return &pb.Job_DataSource{
		Source: &pb.Job_DataSource_Local{
			Local: &pb.Job_Local{},
		},
	}, nil
}

func (s *LocalSource) Override(raw *pb.Job_DataSource, m map[string]string) error {
	if len(m) > 0 {
		return fmt.Errorf("overrides not allowed for local data source")
	}

	return nil
}

func (s *LocalSource) RefToOverride(*pb.Job_DataSource_Ref) (map[string]string, error) {
	return nil, nil
}

func (s *LocalSource) Get(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	raw *pb.Job_DataSource,
	baseDir string,
) (string, *pb.Job_DataSource_Ref, func() error, error) {
	pwd, err := os.Getwd()
	if err == nil && !filepath.IsAbs(pwd) {
		// This should never happen because os.Getwd I believe always
		// returns an absolute path but we want to be absolutely sure
		// so we'll make it abs here.
		pwd, err = filepath.Abs(pwd)
	}

	return pwd, nil, nil, err
}

func (s *LocalSource) Changes(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	source *pb.Job_DataSource,
	current *pb.Job_DataSource_Ref,
	tempDir string,
) (*pb.Job_DataSource_Ref, bool, error) {
	// Never any changes.
	return nil, false, nil
}

var _ Sourcer = (*LocalSource)(nil)
