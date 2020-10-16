package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

func (s *LocalSource) Get(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	raw *pb.Job_DataSource,
	baseDir string,
) (string, func() error, error) {
	return "", nil, nil
}

var _ Sourcer = (*LocalSource)(nil)
