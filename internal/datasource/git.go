package datasource

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type GitSource struct{}

func newGitSource() interface{} { return &GitSource{} }

func (s *GitSource) ProjectSource(body hcl.Body, ctx *hcl.EvalContext) (*pb.Job_DataSource, error) {
	// Decode
	var cfg gitConfig
	if diag := gohcl.DecodeBody(body, ctx, &cfg); len(diag) > 0 {
		return nil, diag
	}

	// Return the data source
	return &pb.Job_DataSource{
		Source: &pb.Job_DataSource_Git{
			Git: &pb.Job_Git{
				Url: cfg.Url,
			},
		},
	}, nil
}

func (s *GitSource) Override(raw *pb.Job_DataSource, m map[string]string) error {
	src := raw.Source.(*pb.Job_DataSource_Git).Git

	var md mapstructure.Metadata
	if err := mapstructure.DecodeMetadata(m, src, &md); err != nil {
		return err
	}

	if len(md.Unused) > 0 {
		return fmt.Errorf("invalid override keys: %v", md.Unused)
	}

	return nil
}

func (s *GitSource) Get(
	ctx context.Context,
	log hclog.Logger,
	raw *pb.Job_DataSource,
	baseDir string,
) (string, func() error, error) {
	source := raw.Source.(*pb.Job_DataSource_Git)

	// Create a temporary directory where we will store the cloned data.
	td, err := ioutil.TempDir(baseDir, "waypoint")
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

type gitConfig struct {
	Url string `hcl:"url,attr"`
}

var _ Sourcer = (*GitSource)(nil)
