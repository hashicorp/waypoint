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
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
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

func (s *GitSource) RefToOverride(ref *pb.Job_DataSource_Ref) (map[string]string, error) {
	gitRef, ok := ref.Ref.(*pb.Job_DataSource_Ref_Git)
	if !ok {
		return nil, fmt.Errorf("ref is not a git ref: %T", ref.Ref)
	}

	return map[string]string{
		"ref": gitRef.Git.Commit,
	}, nil
}

func (s *GitSource) ProjectSource(body hcl.Body, ctx *hcl.EvalContext) (*pb.Job_DataSource, error) {
	// Decode
	var cfg gitConfig
	if diag := gohcl.DecodeBody(body, ctx, &cfg); len(diag) > 0 {
		return nil, diag
	}

	// Start building the result
	result := &pb.Job_Git{
		Url:  cfg.Url,
		Path: cfg.Path,
	}
	switch {
	case cfg.Username != "":
		result.Auth = &pb.Job_Git_Basic_{
			Basic: &pb.Job_Git_Basic{
				Username: cfg.Username,
				Password: cfg.Password,
			},
		}

	case cfg.SSHKey != "":
		// Validate the key
		if _, err := ssh.NewPublicKeys(
			"git",
			[]byte(cfg.SSHKey),
			cfg.SSHKeyPassword,
		); err != nil {
			return nil, fmt.Errorf("failed to load specified Git SSH key: %s", err)
		}

		result.Auth = &pb.Job_Git_Ssh{
			Ssh: &pb.Job_Git_SSH{
				PrivateKeyPem: []byte(cfg.SSHKey),
				Password:      cfg.SSHKeyPassword,
			},
		}
	}

	// Return the data source
	return &pb.Job_DataSource{
		Source: &pb.Job_DataSource_Git{
			Git: result,
		},
	}, nil
}

func (s *GitSource) Override(raw *pb.Job_DataSource, m map[string]string) error {
	src := raw.Source.(*pb.Job_DataSource_Git).Git

	// If we have a username set, then switch auth to basic auth
	if _, ok := m["username"]; ok {
		src.Auth = &pb.Job_Git_Basic_{
			Basic: &pb.Job_Git_Basic{
				Username: m["username"],
				Password: m["password"],
			},
		}

		delete(m, "username")
		delete(m, "password")
	}

	// If we have SSH key set, then change auth to SSH.
	if _, ok := m["key"]; ok {
		src.Auth = &pb.Job_Git_Ssh{
			Ssh: &pb.Job_Git_SSH{
				PrivateKeyPem: []byte(m["key"]),
				Password:      m["key_password"],
			},
		}

		delete(m, "key")
		delete(m, "key_password")
	}

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
) (string, *pb.Job_DataSource_Ref, func() error, error) {
	source := raw.Source.(*pb.Job_DataSource_Git)

	// Some quick validation
	if p := source.Git.Path; p != "" {
		if filepath.IsAbs(p) {
			return "", nil, nil, status.Errorf(codes.FailedPrecondition,
				"git path must be relative")
		}

		for _, part := range filepath.SplitList(p) {
			if part == ".." {
				return "", nil, nil, status.Errorf(codes.FailedPrecondition,
					"git path may not contain '..'")
			}
		}
	}

	// Create a temporary directory where we will store the cloned data.
	td, err := ioutil.TempDir(baseDir, "waypoint")
	if err != nil {
		return "", nil, nil, err
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
	auth, err := s.auth(log, ui, source)
	if err != nil {
		return "", nil, nil, err
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
		return "", nil, nil, status.Errorf(codes.Aborted,
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
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to fetch refspecs: %s", err)
		}

		// ResolveRevision will determine the hash of a short-hash, branch,
		// tag, etc. etc. basically anything "git checkout" accepts.
		hash, err := repo.ResolveRevision(plumbing.Revision(ref))
		if err != nil {
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to resolve revision for checkout: %s", err)
		} else if hash == nil {
			// should never happen but we don't want to panic if it does
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to resolve revision for checkout: nil hash")
		}

		wt, err := repo.Worktree()
		if err != nil {
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to load Git working tree: %s", err)
		}
		if err := wt.Checkout(&git.CheckoutOptions{Hash: *hash}); err != nil {
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Git checkout failed: %s", err)
		}
	}

	// Get our ref
	ref, err := repo.Head()
	if err != nil {
		closer()
		return "", nil, nil, status.Errorf(codes.Aborted,
			"Failed to determine Git HEAD: %s", err)
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		closer()
		return "", nil, nil, status.Errorf(codes.Aborted,
			"Failed to inspect commit information: %s", err)
	}
	var commitTs *timestamp.Timestamp
	if v := commit.Author.When; !v.IsZero() {
		commitTs, err = ptypes.TimestampProto(v)
		if err != nil {
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to inspect commit information: %s", err)
		}
	}

	// If we have a path, set it.
	result := td
	if p := source.Git.Path; p != "" {
		result = filepath.Join(result, p)
	}

	return result, &pb.Job_DataSource_Ref{
		Ref: &pb.Job_DataSource_Ref_Git{
			Git: &pb.Job_Git_Ref{
				Commit:    commit.Hash.String(),
				Timestamp: commitTs,
			},
		},
	}, closer, nil
}

func (s *GitSource) Changes(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	raw *pb.Job_DataSource,
	current *pb.Job_DataSource_Ref,
) (*pb.Job_DataSource_Ref, error) {
	source := raw.Source.(*pb.Job_DataSource_Git)

	// Build our auth mechanism
	auth, err := s.auth(log, ui, source)
	if err != nil {
		return nil, err
	}

	// Get our remote
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{source.Git.Url},
	})

	// List our refs, equivalent to git ls-remote
	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return nil, err
	}

	// Find the target for HEAD
	var headTarget plumbing.ReferenceName
	for _, ref := range refs {
		if ref.Type() == plumbing.SymbolicReference && ref.Name() == "HEAD" {
			headTarget = ref.Target()
			break
		}
	}
	if headTarget == "" {
		return nil, status.Errorf(codes.Internal, "HEAD symbolic ref not found")
	}

	// Get the actual hash for HEAD
	var headRef *plumbing.Reference
	for _, ref := range refs {
		if ref.Name() == headTarget {
			headRef = ref
			break
		}
	}
	if headRef == nil {
		return nil, status.Errorf(codes.Internal, "HEAD hash ref not found")
	}

	// Compare
	if current != nil {
		currentRef := current.Ref.(*pb.Job_DataSource_Ref_Git).Git
		if currentRef.Commit == headRef.Hash().String() {
			log.Trace("HEAD matches current ref, ignoring")
			return nil, nil
		}
	}

	return &pb.Job_DataSource_Ref{
		Ref: &pb.Job_DataSource_Ref_Git{
			Git: &pb.Job_Git_Ref{
				Commit: headRef.Hash().String(),
			},
		},
	}, nil
}

func (s *GitSource) auth(
	log hclog.Logger,
	ui terminal.UI,
	source *pb.Job_DataSource_Git,
) (transport.AuthMethod, error) {
	switch authcfg := source.Git.Auth.(type) {
	case *pb.Job_Git_Basic_:
		ui.Output("Auth: username/password", terminal.WithInfoStyle())
		return &http.BasicAuth{
			Username: authcfg.Basic.Username,
			Password: authcfg.Basic.Password,
		}, nil

	case *pb.Job_Git_Ssh:
		ui.Output("Auth: ssh", terminal.WithInfoStyle())
		auth, err := ssh.NewPublicKeys(
			authcfg.Ssh.User,
			authcfg.Ssh.PrivateKeyPem,
			authcfg.Ssh.Password,
		)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition,
				"Failed to load private key for Git auth: %s", err)
		}

		return auth, nil

	case nil:
		// Do nothing

	default:
		log.Warn("unknown auth configuration, ignoring: %T", source.Git.Auth)
	}

	return nil, nil
}

type gitConfig struct {
	Url            string `hcl:"url,attr"`
	Path           string `hcl:"path,optional"`
	Username       string `hcl:"username,optional"`
	Password       string `hcl:"password,optional"`
	SSHKey         string `hcl:"key,optional"`
	SSHKeyPassword string `hcl:"key_password,optional"`
}

var _ Sourcer = (*GitSource)(nil)
