package datasource

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type GitSource struct{}

func newGitSource() Sourcer { return &GitSource{} }

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
				Url:  cfg.Url,
				Path: cfg.Path,
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
	ui terminal.UI,
	raw *pb.Job_DataSource,
	baseDir string,
) (string, func() error, error) {
	source := raw.Source.(*pb.Job_DataSource_Git)

	// Some quick validation
	if p := source.Git.Path; p != "" {
		if filepath.IsAbs(p) {
			return "", nil, status.Errorf(codes.FailedPrecondition,
				"git path must be relative")
		}

		for _, part := range filepath.SplitList(p) {
			if part == ".." {
				return "", nil, status.Errorf(codes.FailedPrecondition,
					"git path may not contain '..'")
			}
		}
	}

	// Create a temporary directory where we will store the cloned data.
	td, err := ioutil.TempDir(baseDir, "waypoint")
	if err != nil {
		return "", nil, err
	}
	closer := func() error {
		return os.RemoveAll(td)
	}

	// Output
	ui.Output("Cloning data from Git", terminal.WithHeaderStyle())
	ui.Output("URL: %s", source.Git.Url, terminal.WithInfoStyle())
	if source.Git.Ref != "" {
		ui.Output("Ref: %s", source.Git.Ref, terminal.WithInfoStyle())
	}

	// Setup auth information
	var auth transport.AuthMethod
	switch authcfg := source.Git.Auth.(type) {
	case *pb.Job_Git_Basic_:
		ui.Output("Auth: username/password", terminal.WithInfoStyle())
		auth = &http.BasicAuth{
			Username: authcfg.Basic.Username,
			Password: authcfg.Basic.Password,
		}

	case *pb.Job_Git_Ssh:
		ui.Output("Auth: ssh", terminal.WithInfoStyle())
		auth, err = ssh.NewPublicKeys(
			authcfg.Ssh.User,
			authcfg.Ssh.PrivateKeyPem,
			authcfg.Ssh.Password,
		)
		if err != nil {
			return "", nil, status.Errorf(codes.FailedPrecondition,
				"Failed to load private key for Git auth: %s", err)
		}

	case nil:
		// Do nothing

	default:
		log.Warn("unknown auth configuration, ignoring: %T", source.Git.Auth)
	}

	// Clone
	var output bytes.Buffer
	repo, err := git.PlainCloneContext(ctx, td, false, &git.CloneOptions{
		URL:      source.Git.Url,
		Auth:     auth,
		Progress: &output,
	})
	if err != nil {
		closer()
		return "", nil, status.Errorf(codes.Aborted,
			"Git clone failed: %s", output.String())
	}

	// Checkout if we have a ref. If we don't have a ref we use the
	// default of whatever we got.
	if ref := source.Git.Ref; ref != "" {
		// We have to fetch all the refs so that ResolveRevisoin can find them.
		err = repo.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		})
		if err != nil {
			closer()
			return "", nil, status.Errorf(codes.Aborted,
				"Failed to fetch refspecs: %s", err)
		}

		// ResolveRevision will determine the hash of a short-hash, branch,
		// tag, etc. etc. basically anything "git checkout" accepts.
		hash, err := repo.ResolveRevision(plumbing.Revision(ref))
		if err != nil {
			closer()
			return "", nil, status.Errorf(codes.Aborted,
				"Failed to resolve revision for checkout: %s", err)
		} else if hash == nil {
			// should never happen but we don't want to panic if it does
			closer()
			return "", nil, status.Errorf(codes.Aborted,
				"Failed to resolve revision for checkout: nil hash")
		}

		wt, err := repo.Worktree()
		if err != nil {
			closer()
			return "", nil, status.Errorf(codes.Aborted,
				"Failed to load Git working tree: %s", err)
		}
		if err := wt.Checkout(&git.CheckoutOptions{Hash: *hash}); err != nil {
			closer()
			return "", nil, status.Errorf(codes.Aborted,
				"Git checkout failed: %s", err)
		}
	}

	// If we have a path, set it.
	result := td
	if p := source.Git.Path; p != "" {
		result = filepath.Join(result, p)
	}

	return result, closer, nil
}

type gitConfig struct {
	Url  string `hcl:"url,attr"`
	Path string `hcl:"path,optional"`
}

var _ Sourcer = (*GitSource)(nil)
