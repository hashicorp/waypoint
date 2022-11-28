package datasource

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mitchellh/mapstructure"
	cryptossh "golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type GitSource struct{}

func newGitSource() Sourcer { return &GitSource{} }

func (s *GitSource) RefToOverride(ref *pb.Job_DataSource_Ref) (map[string]string, error) {
	gitRef, ok := ref.Ref.(*pb.Job_DataSource_Ref_Git)
	if !ok {
		return nil, status.Errorf(codes.Internal, "ref is not a git ref: %T", ref.Ref)
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
		Url:                      cfg.Url,
		Path:                     cfg.Path,
		IgnoreChangesOutsidePath: cfg.IgnoreChangesOutsidePath,
		Ref:                      cfg.Ref,
		RecurseSubmodules:        cfg.RecurseSubmodules,
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

	// Output git info
	// NOTE(briancain): The leading whitespace here is to fit the formatting
	// for when we display the commit, timestamp, and message later on so that the
	// messages are all aligned.
	ui.Output("Cloning data from Git", terminal.WithHeaderStyle())
	ui.Output("       URL: %s", source.Git.Url, terminal.WithInfoStyle())
	if source.Git.Ref != "" {
		ui.Output("       Ref: %s", source.Git.Ref, terminal.WithInfoStyle())
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

		// Note: we don't set RecurseSubmodules here because if we're checking
		// out a ref without submodules or with different submodules, we
		// don't want to waste time recursing HEAD. We fetch submodules
		// later.
	})
	if err != nil {
		closer()

		return "", nil, nil, status.Errorf(codes.Aborted,
			"Git clone failed: %s\n\nOutput: %s", err, output.String())
	}

	// Checkout if we have a ref. If we don't have a ref we use the
	// default of whatever we got.
	if ref := source.Git.Ref; ref != "" {
		// We have to fetch all the refs so that ResolveRevisoin can find them.
		err = repo.Fetch(&git.FetchOptions{
			Auth:     auth,
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

	if depth := source.Git.RecurseSubmodules; depth > 0 {
		wt, err := repo.Worktree()
		if err != nil {
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to load Git working tree: %s", err)
		}

		sm, err := wt.Submodules()
		if err != nil {
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to load submodules: %s", err)
		}

		if err := sm.UpdateContext(ctx, &git.SubmoduleUpdateOptions{
			Init:              true,
			Auth:              auth,
			RecurseSubmodules: git.SubmoduleRescursivity(depth),
		}); err != nil {
			closer()
			return "", nil, nil, status.Errorf(codes.Aborted,
				"Failed to update submodules: %s", err)
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
	var commitTs *timestamppb.Timestamp
	if v := commit.Author.When; !v.IsZero() {
		commitTs = timestamppb.New(v)
	}

	// If we have a path, set it.
	result := td
	if p := source.Git.Path; p != "" {
		result = filepath.Join(result, p)
	}

	// Output additinoal git info
	ui.Output("Git Commit: %s", commit.Hash.String(), terminal.WithInfoStyle())
	ui.Output(" Timestamp: %s", commitTs.AsTime(), terminal.WithInfoStyle())
	ui.Output("   Message: %s", commit.Message, terminal.WithInfoStyle())

	return result, &pb.Job_DataSource_Ref{
		Ref: &pb.Job_DataSource_Ref_Git{
			Git: &pb.Job_Git_Ref{
				Commit:        commit.Hash.String(),
				Timestamp:     commitTs,
				CommitMessage: commit.Message,
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
	tempDir string,
) (*pb.Job_DataSource_Ref, bool, error) {
	source := raw.Source.(*pb.Job_DataSource_Git)

	// Build our auth mechanism
	auth, err := s.auth(log, ui, source)
	if err != nil {
		return nil, false, err
	}

	// Get our remote
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{source.Git.Url},
	})

	// Determine our target ref. If no Ref is specified by the user,
	// we default to HEAD which points to whatever the default branch is.
	targetRef := source.Git.Ref
	if targetRef == "" {
		targetRef = "HEAD"
	}

	// List our refs, equivalent to git ls-remote
	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return nil, false, err
	}

	// Build a map of our refs since we may have to do many lookups
	refMap := map[plumbing.ReferenceName]*plumbing.Reference{}
	for _, ref := range refs {
		// We should never get a duplicate ref, but if we do, let's log it
		// because this will PROBABLY result in some weird behavior.
		if _, ok := refMap[ref.Name()]; ok {
			log.Warn("duplicate ref in ls-remote, this shouldn't happen",
				"name", ref.Name())
		}

		refMap[ref.Name()] = ref
	}

	// tryRefs is the list of refs we will try to find in the ref list.
	// The first ref found that matches will be returned first. This contains
	// all the formats a user may enter for a ref. The order at the time
	// of writing is the following, where %s is replaced with the user-supplied
	// ref. This is the same logic `git rev-parse` uses.
	//
	//  - "%s",
	//  - "refs/%s",
	//  - "refs/tags/%s",
	//  - "refs/heads/%s",
	//  - "refs/remotes/%s",
	//  - "refs/remotes/%s/HEAD",
	//
	tryRefs := []string{targetRef}
	for _, rule := range plumbing.RefRevParseRules {
		tryRefs = append(tryRefs, fmt.Sprintf(rule, targetRef))
	}
	var foundRef *plumbing.Reference
	for _, tryRef := range tryRefs {
		ref, ok := refMap[plumbing.ReferenceName(tryRef)]
		if !ok {
			continue
		}

		// The limit prevents adversarial or buggy remote git repos that
		// might have a symbolic reference loop. The value 10 was chosen
		// arbitrarily, I've never seen a reference repeat more than 1 time.
		limit := 10

		// If the ref is a symbolic reference, we dereference until we find
		// the target. An example here is HEAD may point to refs/heads/main,
		// and so on.
		for limit > 0 && ref != nil && ref.Type() == plumbing.SymbolicReference {
			ref = refMap[ref.Target()]
			limit--
		}
		if limit == 0 {
			return nil, false, status.Errorf(codes.FailedPrecondition,
				"Infinite reference in hash lookup for target: %s", targetRef)
		}
		if ref == nil {
			continue
		}

		foundRef = ref
		break
	}

	if foundRef == nil {
		return nil, false, status.Errorf(codes.Internal, "Hash for target ref not found: %s", targetRef)
	}

	// We maybe have changes. We set the ignore value later if we
	// determine that the changes we have should be ignored.
	ignore := false

	// Compare
	if current != nil {
		currentRef := current.Ref.(*pb.Job_DataSource_Ref_Git).Git
		if currentRef.Commit == foundRef.Hash().String() {
			log.Trace("current ref matches last known ref, ignoring")
			return nil, false, nil
		}

		// If there is a subpath specified for the data, then we go one step
		// further and determine if there are any changes in this specific
		// path. To do that, we have no choice but to check out the whole repo.
		// We only do this if we had a previous value we used. If we don't,
		// we've never run before and we consider ANY new ref changes.
		if path := source.Git.Path; path != "" && source.Git.IgnoreChangesOutsidePath {
			log.Trace("subpath specified, we'll check for changes within the subpath")
			changes, err := s.changes(
				ctx,
				log,
				raw,
				tempDir,
				foundRef.Hash().String(),
				currentRef.Commit,
				path,
			)
			if err != nil {
				return nil, false, err
			}

			// We ignore if there are no changes in our subpath
			ignore = !changes
		}
	}

	return &pb.Job_DataSource_Ref{
		Ref: &pb.Job_DataSource_Ref_Git{
			Git: &pb.Job_Git_Ref{
				Commit: foundRef.Hash().String(),
			},
		},
	}, ignore, nil
}

func (s *GitSource) changes(
	ctx context.Context,
	log hclog.Logger,
	raw *pb.Job_DataSource,
	baseDir string,
	refNew, refCurrent string,
	subpath string,
) (bool, error) {
	source := raw.Source.(*pb.Job_DataSource_Git)

	// Normalize our subpath:
	//   ./foo => foo
	//   /foo  => foo
	//
	// Note that this SHOULD happen at the API level (we should reject
	// any weird values) but we introduced validation later so we do this
	// to clean up potentially old dirty values.
	subpath = filepath.ToSlash(subpath)
	if len(subpath) >= 2 && subpath[0] == '.' && subpath[1] == filepath.Separator {
		subpath = subpath[2:]
	} else if len(subpath) >= 1 && subpath[0] == filepath.Separator {
		subpath = subpath[1:]
	}

	// Ensure our path ends with a '/' for comparison purposes later.
	if !strings.HasSuffix(subpath, string(filepath.Separator)) {
		subpath += string(filepath.Separator)
	}

	// Create a temporary directory where we will store the cloned data.
	td, err := ioutil.TempDir(baseDir, "waypoint")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(td)

	// Setup auth information
	auth, err := s.auth(log, nil, source)
	if err != nil {
		return false, err
	}

	// Clone
	var output bytes.Buffer
	log.Trace("cloning repository", "url", source.Git.Url)
	repo, err := git.PlainCloneContext(ctx, td, false, &git.CloneOptions{
		URL:      source.Git.Url,
		Auth:     auth,
		Progress: &output,
	})
	if err != nil {
		return false, status.Errorf(codes.Aborted,
			"Git clone failed: %s", output.String())
	}

	// We have to fetch all the refs so that ResolveRevisoin can find them.
	log.Trace("fetching refs")
	err = repo.Fetch(&git.FetchOptions{
		Auth:     auth,
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		return false, status.Errorf(codes.Aborted,
			"Failed to fetch refspecs: %s", err)
	}

	// ResolveRevision will determine the hash of a short-hash, branch,
	// tag, etc. etc. basically anything "git checkout" accepts.
	hashCurrent, err := repo.ResolveRevision(plumbing.Revision(refCurrent))
	if err != nil {
		return false, status.Errorf(codes.Aborted,
			"Failed to resolve current revision for checkout: %s", err)
	} else if hashCurrent == nil {
		// should never happen but we don't want to panic if it does
		return false, status.Errorf(codes.Aborted,
			"Failed to resolve current revision for checkout: nil hash")
	}

	hashNew, err := repo.ResolveRevision(plumbing.Revision(refNew))
	if err != nil {
		return false, status.Errorf(codes.Aborted,
			"Failed to resolve new revision for checkout: %s", err)
	} else if hashNew == nil {
		// should never happen but we don't want to panic if it does
		return false, status.Errorf(codes.Aborted,
			"Failed to resolve new revision for checkout: nil hash")
	}

	// Get our trees so we can compare
	var treeCurrent, treeNew *object.Tree
	log.Trace("getting current tree object", "hash", *hashCurrent)
	commitCurrent, err := repo.CommitObject(*hashCurrent)
	if err == nil {
		treeCurrent, err = commitCurrent.Tree()
	}
	if err != nil {
		return false, status.Errorf(codes.Aborted,
			"Failed to get current tree object: %s", err)
	}
	log.Trace("getting new tree object", "hash", *hashNew)
	commitNew, err := repo.CommitObject(*hashNew)
	if err == nil {
		treeNew, err = commitNew.Tree()
	}
	if err != nil {
		return false, status.Errorf(codes.Aborted,
			"Failed to get new tree object: %s", err)
	}

	log.Trace("diffing")
	changes, err := treeCurrent.DiffContext(ctx, treeNew)
	if err != nil {
		return false, status.Errorf(codes.Aborted,
			"Failed to diff tree objects: %s", err)
	}

	// Go through each change in our diff and if it has anything to
	// do in our subpath (no matter what the action is) then we detect
	// changes.
	for _, change := range changes {
		// The slashes in the Git object I _think_ are always / but regardless
		// we want to make sure our slashes match our subpath.
		from := filepath.ToSlash(change.From.Name)
		to := filepath.ToSlash(change.To.Name)
		if strings.HasPrefix(from, subpath) || strings.HasPrefix(to, subpath) {
			log.Trace("detected change", "change", change.String())
			return true, nil
		}
	}

	// No changes
	return false, nil
}

func (s *GitSource) auth(
	log hclog.Logger,
	ui terminal.UI,
	source *pb.Job_DataSource_Git,
) (transport.AuthMethod, error) {
	switch authcfg := source.Git.Auth.(type) {
	case *pb.Job_Git_Basic_:
		if ui != nil {
			ui.Output("      Auth: username/password", terminal.WithInfoStyle())
		}
		return &http.BasicAuth{
			Username: authcfg.Basic.Username,
			Password: authcfg.Basic.Password,
		}, nil

	case *pb.Job_Git_Ssh:
		// Default the user to "git" which is typically what is used.
		user := authcfg.Ssh.User
		if user == "" {
			user = "git"
		}

		if ui != nil {
			ui.Output("      Auth: ssh", terminal.WithInfoStyle())
		}
		auth, err := ssh.NewPublicKeys(
			user,
			authcfg.Ssh.PrivateKeyPem,
			authcfg.Ssh.Password,
		)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition,
				"Failed to load private key for Git auth: %s", err)
		}

		// We do not do any host key verification for now.
		// NOTE(mitchellh): in the future we should expose a way to
		// configure enabling this in some way.
		auth.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()

		return auth, nil

	case nil:
		// Do nothing

	default:
		log.Warn("unknown auth configuration, ignoring: %T", source.Git.Auth)
	}

	return nil, nil
}

type gitConfig struct {
	Url                      string `hcl:"url,attr"`
	Path                     string `hcl:"path,optional"`
	Username                 string `hcl:"username,optional"`
	Password                 string `hcl:"password,optional"`
	SSHKey                   string `hcl:"key,optional"`
	SSHKeyPassword           string `hcl:"key_password,optional"`
	Ref                      string `hcl:"ref,optional"`
	IgnoreChangesOutsidePath bool   `hcl:"ignore_changes_outside_path,optional"`
	RecurseSubmodules        uint32 `hcl:"recurse_submodules,optional"`
}

var _ Sourcer = (*GitSource)(nil)
