package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

var testHasGit bool

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
}

func TestRunnerAccept(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

func TestRunnerAccept_cancelContext(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Set a blocker
	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Cancel the context eventually. This isn't CI-sensitive cause
	// we'll block no matter what.
	time.AfterFunc(500*time.Millisecond, cancel)

	// Accept should complete with an error
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	require.Eventually(func() bool {
		job, err := client.GetJob(context.Background(), &pb.GetJobRequest{JobId: jobId})
		require.NoError(err)
		return job.State == pb.Job_ERROR
	}, 3*time.Second, 25*time.Millisecond)
}

func TestRunnerAccept_cancelJob(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Set a blocker
	noopCh := make(chan struct{})
	runner.noopCh = noopCh

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Cancel the context eventually. This isn't CI-sensitive cause
	// we'll block no matter what.
	time.AfterFunc(500*time.Millisecond, func() {
		_, err := client.CancelJob(ctx, &pb.CancelJobRequest{
			JobId: jobId,
		})
		require.NoError(err)
	})

	// Accept should complete with an error
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	require.Eventually(func() bool {
		job, err := client.GetJob(context.Background(), &pb.GetJobRequest{JobId: jobId})
		require.NoError(err)
		return job.State == pb.Job_ERROR
	}, 3*time.Second, 25*time.Millisecond)
}

func TestRunnerAccept_gitData(t *testing.T) {
	if !testHasGit {
		t.Skip("git not installed")
		return
	}

	require := require.New(t)
	ctx := context.Background()

	// Get a repo path
	path := testGitFixture(t, "git-noop")

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, &pb.Job{
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: path,
					},
				},
			},
		}),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
	require.NotNil(job.DataSourceRef)
}

func TestRunnerAccept_noConfig(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Change our directory to a new temp directory with no config file.
	testChdir(t, testTempDir(t))

	// Initialize our app
	singleprocess.TestApp(t, client, serverptypes.TestJobNew(t, nil).Application)

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_ERROR, job.State)
	require.NotNil(job.Error)

	st := status.FromProto(job.Error)
	require.Equal(codes.FailedPrecondition, st.Code())
}

func TestRunnerAccept_noConfig_serverHcl(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Setup our runner
	client := singleprocess.TestServer(t)
	runner := TestRunner(t, WithClient(client))
	require.NoError(runner.Start())

	// Change our directory to a new temp directory with no config file.
	testChdir(t, testTempDir(t))

	// Initialize our app
	ref := serverptypes.TestJobNew(t, nil).Application
	{
		_, err := client.UpsertProject(context.Background(), &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name:        ref.Project,
				WaypointHcl: []byte(configpkg.TestSource(t)),
			},
		})
		require.NoError(err)
	}

	{
		_, err := client.UpsertApplication(context.Background(), &pb.UpsertApplicationRequest{
			Project: &pb.Ref_Project{Project: ref.Project},
			Name:    ref.Application,
		})
		require.NoError(err)
	}

	// Queue a job
	queueResp, err := client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: serverptypes.TestJobNew(t, nil),
	})
	require.NoError(err)
	jobId := queueResp.JobId

	// Accept should complete
	require.NoError(runner.Accept(ctx))

	// Verify that the job is completed
	job, err := client.GetJob(ctx, &pb.GetJobRequest{JobId: jobId})
	require.NoError(err)
	require.Equal(pb.Job_SUCCESS, job.State)
}

// testGitFixture MUST be called before TestRunner since TestRunner
// changes our working directory.
func testGitFixture(t *testing.T, n string) string {
	t.Helper()

	// We need to get our working directory since the TestRunner call
	// changes it.
	wd, err := os.Getwd()
	require.NoError(t, err)
	wd, err = filepath.Abs(wd)
	require.NoError(t, err)
	path := filepath.Join(wd, "testdata", n)

	// Look for a DOTgit
	original := filepath.Join(path, "DOTgit")
	_, err = os.Stat(original)
	require.NoError(t, err)

	// Rename it
	newPath := filepath.Join(path, ".git")
	require.NoError(t, os.Rename(original, newPath))
	t.Cleanup(func() { os.Rename(newPath, original) })

	return path
}
