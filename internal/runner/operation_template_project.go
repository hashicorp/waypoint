package runner

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v48/github"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var (
	PROJECT_NAME_CONTENTS_REGEX       = regexp.MustCompile("%%wp_project%%")
	PROJECT_NAME_CONTENTS_UPPER_REGEX = regexp.MustCompile("%%Wp_project%%")

	PROJECT_NAME_PATH_REGEX       = regexp.MustCompile("__wp_project__")
	PROJECT_NAME_PATH_UPPER_REGEX = regexp.MustCompile("__Wp_project__")

	PROJECT_DESCRIPTION_REGEX = regexp.MustCompile("%%wp_project_description%%")
)

// executeProjectTemplateOp creates a project based on a template
// For a github template, it creates a new repository based on the github repo template,
// clones it, replaces tokens in file names or contents, pushes the changes,
// and finally updates the project in waypoint to have the correct datasource ref.
//
// It isn't great that the runner upserts the project at the last minute. It would
// be better if the server could handle that after this job is complete, but alas,
// there is no mechanism. The server _could_ predict the final datasource at request time
// if we needed to roll this back to narrow down permissions.
func (r *Runner) executeProjectTemplateOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	op, ok := job.Operation.(*pb.Job_TemplateProject)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type for template project")
	}
	templateReq := op.TemplateProject.Req

	projectTemplateResp, err := r.client.GetProjectTemplate(ctx, &pb.GetProjectTemplateRequest{
		ProjectTemplate: templateReq.Template,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project template named %q", templateReq.Template)
	}
	projectTemplate := projectTemplateResp.ProjectTemplate

	githubTemplate, ok := projectTemplate.SourceCodePlatform.(*pb.ProjectTemplate_Github)
	if !ok {
		return nil, errors.Errorf("unsupported project template type %t", projectTemplate.SourceCodePlatform)
	}

	githubDestOptions, ok := templateReq.SourceCodePlatformDestinationOptions.(*pb.UpsertProjectFromTemplateRequest_Github)
	if !ok {
		return nil, errors.Errorf("unsupport project options type %t", templateReq.SourceCodePlatformDestinationOptions)
	}

	githubUsername := os.Getenv("GITHUB_USERNAME")
	if githubUsername == "" {
		return nil, errors.Errorf("template project op requires runner to have GITHUB_USERNAME env var set")
	}

	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if githubAccessToken == "" {
		return nil, errors.Errorf("template project op requires runner to have GITHUB_ACCESS_TOKEN env var set")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	githubClient := github.NewClient(oauth2.NewClient(ctx, ts))

	project.UI.Output("Creating repo from template...")
	githubRepo, _, err := githubClient.Repositories.CreateFromTemplate(ctx,
		githubTemplate.Github.Source.Owner,
		githubTemplate.Github.Source.Repo,
		&github.TemplateRepoRequest{
			// Options for all template types
			Name:               &templateReq.ProjectName,
			Description:        &templateReq.Description,
			Owner:              &githubDestOptions.Github.Owner,
			IncludeAllBranches: &githubTemplate.Github.Destination.IncludeAllBranches,
			Private:            &githubTemplate.Github.Destination.Private,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create github repo from template")
	}

	//ptr := func(i string) *string { return &i }
	//fmt.Println(githubTemplate, githubDestOptions, githubClient)
	//githubRepo := &github.Repository{
	//	URL:      ptr("https://api.github.com/repos/izaaklauer/deletemerepo1"),
	//	CloneURL: ptr("https://github.com/izaaklauer/deletemerepo1.git"),
	//	GitURL:   ptr("git://github.com/izaaklauer/deletemerepo1.git"),
	//}

	// Clone the repo

	td, err := ioutil.TempDir("", fmt.Sprintf("waypoint-%s", templateReq.ProjectName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create tmp dir for cloning newly created repo")
	}
	defer os.RemoveAll(td)

	githubAuth := &http.BasicAuth{
		Username: githubUsername,
		Password: githubAccessToken,
	}

	// In practice, we seem to need to give github a moment before cloning.
	time.Sleep(time.Second * 5)

	project.UI.Output("Creating new repo...")
	var output bytes.Buffer
	clonedRepo, err := git.PlainCloneContext(ctx, td, false, &git.CloneOptions{
		URL:      *githubRepo.CloneURL,
		Auth:     githubAuth,
		Progress: &output,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to git clone")
	}

	// Render template
	project.UI.Output("Rendering template...")
	if err := replaceTokens(td, templateReq.ProjectName, templateReq.Description); err != nil {
		return nil, errors.Wrapf(err, "failed replacing tokens")
	}

	worktree, err := clonedRepo.Worktree()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get worktree for cloned repo")
	}

	_, err = worktree.Add(".")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to git add after replacing tokens")
	}
	// Unfortunately, the above doesn't add deleted files: https://github.com/go-git/go-git/issues/113#issuecomment-758084471
	// worktree.Move also doesn't work on directories, so I can't `git mv` them instead of renaming.
	// The sad workaround:
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = worktree.Filesystem.Root()
	err = cmd.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get add deleted refs after replacing tokens")
	}

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name: "waypoint",
			When: time.Now(),
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to commit after replacing tokens")
	}

	project.UI.Output("Pushing rendering changes...")
	var pushOutput bytes.Buffer
	err = clonedRepo.PushContext(ctx, &git.PushOptions{
		Auth:     githubAuth,
		Progress: &pushOutput,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to push after replacing tokens")
	}

	// Create the final version of this project, which includes
	pbProject := projectTemplate.ProjectSettings
	pbProject.Name = templateReq.ProjectName
	pbProject.DataSource = &pb.Job_DataSource{
		Source: &pb.Job_DataSource_Git{
			Git: &pb.Job_Git{
				Url:  *githubRepo.CloneURL,
				Ref:  *githubRepo.DefaultBranch,
				Path: "",
				// TODO: figure out auth
			},
		},
	}

	// Clear out saved hack waypoint.hcl
	pbProject.WaypointHcl = []byte{}

	project.UI.Output("Upserting final project...")
	_, err = r.client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: pbProject,
	})

	project.UI.Output("Waiting for project init...")
	time.Sleep(time.Second * 15)

	return &pb.Job_Result{
		TemplateProjectResult: &pb.Job_TemplateProjectResult{},
	}, nil
}

func replaceTokens(path string, projectName string, projectDescription string) error {

	// Collect files and dirs to rename
	var renameFilesAndDirs []string
	if err := filepath.WalkDir(path, func(filePath string, f os.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed while walking")
		}
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		if PROJECT_NAME_PATH_REGEX.FindString(f.Name()) != "" ||
			PROJECT_NAME_PATH_UPPER_REGEX.FindString(f.Name()) != "" {
			renameFilesAndDirs = append(renameFilesAndDirs, filePath)
		}

		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to collect files and dirs")
	}

	// Rename files and dirs, iterate backwards to rename leaf nodes first
	for i := len(renameFilesAndDirs) - 1; i >= 0; i-- {
		oldPath := renameFilesAndDirs[i]

		name := filepath.Base(oldPath)

		newName := PROJECT_NAME_PATH_REGEX.ReplaceAllString(name, projectName)
		newName = PROJECT_NAME_PATH_UPPER_REGEX.ReplaceAllString(newName, strings.Title(projectName))

		if newName != name {
			newPath := filepath.Join(filepath.Dir(oldPath), newName)
			if err := os.Rename(oldPath, newPath); err != nil {
				// Assumes file names aren't secret
				return errors.Wrapf(err, "failed to rename file %q to %q", name, newName)
			}
		}
	}

	// Apply templating
	err := filepath.WalkDir(path, func(filePath string, f os.DirEntry, err error) error {
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}

		if f.IsDir() {
			return nil
		}
		info, err := f.Info()
		if err != nil {
			return errors.Wrapf(err, "failed to get info on file %q", f.Name())
		}

		contents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return errors.Wrapf(err, "failed to read file %q", f.Name())
		}

		contents = PROJECT_NAME_CONTENTS_REGEX.ReplaceAll(contents, []byte(projectName))
		contents = PROJECT_NAME_CONTENTS_UPPER_REGEX.ReplaceAll(contents, []byte(strings.Title(projectName)))

		contents = PROJECT_DESCRIPTION_REGEX.ReplaceAll(contents, []byte(projectDescription))

		if err := ioutil.WriteFile(filePath, contents, info.Mode()); err != nil {
			return errors.Wrapf(err, "failed to write new contents of file %q", f.Name())
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to template files")
	}

	return nil
}
