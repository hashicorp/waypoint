package runner

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

func TestRunnerProjectTemplateOp(t *testing.T) {
	// This test is not fully automated - it expects some setup, and
	// is really just a harness for manually testing executeProjectTemplateOp at this point
	t.SkipNow()

	require := require.New(t)

	ctx := context.Background()
	client := singleprocess.TestServer(t)
	log := testLog()

	// Create and start our runner
	runner, err := New(
		WithClient(client),
		WithLogger(log),
	)
	require.NoError(err)
	defer runner.Close()
	require.NoError(runner.Start(ctx))

	project := core.TestProject(t, core.WithClient(client))

	repoName := "deletemerepo1"
	job := &pb.Job{
		Operation: &pb.Job_TemplateProject{
			TemplateProject: &pb.Job_TemplateProjectOp{
				Req: &pb.UpsertProjectFromTemplateRequest{
					ProjectName: repoName,
					Description: repoName,
					SourceCodePlatformDestinationOptions: &pb.UpsertProjectFromTemplateRequest_Github{
						Github: &pb.ProjectTemplate_SourceCodePlatformGithub_Destination_Options{
							Owner: "izaaklauer",
						},
					},
					Template: &pb.Ref_ProjectTemplate{
						Name: "template-1",
					},
				},
			},
		},
	}

	res, err := runner.executeProjectTemplateOp(ctx, runner.logger, job, project)

	require.NoError(err)
	require.NotNil(res)
}

func TestRunner_replaceTokens(t *testing.T) {
	require := require.New(t)

	// Copy our test data
	wd, err := os.Getwd()
	require.NoError(err)
	wd, err = filepath.Abs(wd)
	require.NoError(err)
	referenceTestDataPath := filepath.Join(wd, "testdata", "template-repo")
	testDataPath, err := ioutil.TempDir("", "template-repo_")
	require.NoError(err)
	cmd := exec.Command("cp", "-r", referenceTestDataPath, testDataPath)
	require.NoError(cmd.Run())
	defer func() { require.NoError(os.RemoveAll(testDataPath)) }()

	err = replaceTokens(testDataPath, "hashicups", "description of test project")
	require.NoError(err)

	// TODO: examine file names and contents for tokens.
	// Without that, this function is just a fixture for running replaceTokens and
	// allowing you to verify manually
}
