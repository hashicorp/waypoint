package singleprocess

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess/mocks"
)

func TestPollQueuer_peek(t *testing.T) {
	// with this we're basically testing that the poller just doesn't crash
	// if Peek returns literally nothing.
	t.Run("zero result", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer wg.Wait()
		defer cancel()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)

		// Create our mock handler
		mockH := &mocks.PollHandler{}

		// Return zero values
		var counter uint32
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(nil, time.Time{}, nil).
			Run(func(args mock.Arguments) {
				// Count how many times we've peeked
				atomic.AddUint32(&counter, 1)
			})

		// Start
		wg.Add(1)
		go testServiceImpl(impl).runPollQueuer(ctx, &wg, mockH, nil)

		// What we're testing here is that we eventually call Peek
		// and that we call it a reasonable number of times. And we don't crash!
		require.Eventually(func() bool {
			count := atomic.LoadUint32(&counter)
			if count == 0 {
				return false
			}

			// We should poll exactly once cause we're stuck in a wait loop
			require.EqualValues(count, 1)
			return true
		}, 10*time.Second, 10*time.Millisecond)

		// Roughly test we never call PollJob. If we do AFTER this, its
		// okay, we have some other tests to verify some more. We don't
		// need an assertion cause the mock will fail.
		time.Sleep(100 * time.Millisecond)
	})

	// Test that if the watchset triggers while we're waiting, we re-peek
	t.Run("watchset trigger", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer wg.Wait()
		defer cancel()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)

		// Create our mock handler
		mockH := &mocks.PollHandler{}

		// On our first peek call, setup the watchset
		wsCh := make(chan struct{})
		peekCalled := make(chan struct{})
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(nil, time.Time{}, nil).
			Run(func(args mock.Arguments) {
				defer close(peekCalled)
				ws := args.Get(1).(memdb.WatchSet)
				ws.Add(wsCh)
			}).
			Once()

		// On our second peek call, return nothing
		peek2Called := make(chan struct{})
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(nil, time.Time{}, nil).
			Run(func(args mock.Arguments) {
				defer close(peek2Called)
			}).
			Once()

		// Start
		wg.Add(1)
		go testServiceImpl(impl).runPollQueuer(ctx, &wg, mockH, nil)

		// Wait for Peek to be called
		select {
		case <-peekCalled:
		case <-time.After(1 * time.Second):
			t.Fatal("never called peek")
		}

		// Trigger our watchset
		close(wsCh)

		// Peek #2 should be called
		select {
		case <-peekCalled:
		case <-time.After(1 * time.Second):
			t.Fatal("never called peek after watchset")
		}
	})

	// Test that watchset takes priority over a pending poll timer.
	t.Run("long poll with watchset trigger", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer wg.Wait()
		defer cancel()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)

		// Create our mock handler
		mockH := &mocks.PollHandler{}

		// On our first peek call, setup the watchset
		wsCh := make(chan struct{})
		peekCalled := make(chan struct{})
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(nil, time.Now().Add(1*time.Minute), nil).
			Run(func(args mock.Arguments) {
				defer close(peekCalled)
				ws := args.Get(1).(memdb.WatchSet)
				ws.Add(wsCh)
			}).
			Once()

		// On our second peek call, return nothing
		peek2Called := make(chan struct{})
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(nil, time.Time{}, nil).
			Run(func(args mock.Arguments) {
				defer close(peek2Called)
			}).
			Once()

		// Start
		wg.Add(1)
		go testServiceImpl(impl).runPollQueuer(ctx, &wg, mockH, nil)

		// Wait for Peek to be called
		select {
		case <-peekCalled:
		case <-time.After(1 * time.Second):
			t.Fatal("never called peek")
		}

		// Trigger our watchset
		close(wsCh)

		// Peek #2 should be called
		select {
		case <-peekCalled:
		case <-time.After(1 * time.Second):
			t.Fatal("never called peek after watchset")
		}
	})
}

func TestPollQueuer_queue(t *testing.T) {
	// If the poll item is nil, then we should repeek and not CPU saturate.
	t.Run("poll timeout with nil item", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer wg.Wait()
		defer cancel()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)

		// Create our mock handler
		mockH := &mocks.PollHandler{}

		// Return our fast poll result: ready in 1ms, but nil item
		var counter uint32
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(nil, time.Now().Add(1*time.Millisecond), nil).
			Run(func(args mock.Arguments) {
				atomic.AddUint32(&counter, 1)
			})

		// Start
		wg.Add(1)
		go testServiceImpl(impl).runPollQueuer(ctx, &wg, mockH, nil)

		// Let's just run this for awhile
		time.Sleep(250 * time.Millisecond)

		// We should only peek once
		require.EqualValues(1, atomic.LoadUint32(&counter))
	})

	// If the poll job errors, we should re-peek, but not too much.
	t.Run("poll job error", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer wg.Wait()
		defer cancel()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)

		// Create our mock handler
		mockH := &mocks.PollHandler{}

		// Return our fast poll result: ready in 1ms, but nil item
		var peekCounter, pollCounter uint32
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(42, time.Now().Add(1*time.Millisecond), nil).
			Run(func(args mock.Arguments) {
				atomic.AddUint32(&peekCounter, 1)
			})
		mockH.On("PollJob", mock.Anything, 42).
			Return(nil, errors.New("oh no")).
			Run(func(args mock.Arguments) {
				atomic.AddUint32(&pollCounter, 1)
			})

		// Start
		wg.Add(1)
		go testServiceImpl(impl).runPollQueuer(ctx, &wg, mockH, nil)

		// Let's just run this for awhile
		time.Sleep(250 * time.Millisecond)

		// We should peek more than once but less than a bunch. This is a rough
		// heuristic for a "reasonable amount" to not CPU saturate.
		require.True(atomic.LoadUint32(&peekCounter) > 0)
		require.True(atomic.LoadUint32(&peekCounter) < 100, "excessive peeking")
		require.True(atomic.LoadUint32(&pollCounter) > 0)
		require.True(atomic.LoadUint32(&pollCounter) < 100, "excessive polling")
	})

	// If the poll job errors, we should re-peek, but not too much.
	t.Run("queue job error", func(t *testing.T) {
		require := require.New(t)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer wg.Wait()
		defer cancel()

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)

		// Create our mock handler
		mockH := &mocks.PollHandler{}

		// Returning a queue job request that is empty should fail validation.
		var peekCounter, pollCounter uint32
		mockH.On("Peek", mock.Anything, mock.Anything).
			Return(42, time.Now().Add(1*time.Millisecond), nil).
			Run(func(args mock.Arguments) {
				atomic.AddUint32(&peekCounter, 1)
			})
		mockH.On("PollJob", mock.Anything, 42).
			Return(nil, errors.New("oh no")).
			Run(func(args mock.Arguments) {
				atomic.AddUint32(&pollCounter, 1)
			})
		mockH.On("Complete", mock.Anything, 42).
			Return(nil, nil).
			Run(func(args mock.Arguments) {
			})

		// Start
		wg.Add(1)
		go testServiceImpl(impl).runPollQueuer(ctx, &wg, mockH, nil)

		// Let's just run this for awhile
		time.Sleep(250 * time.Millisecond)

		// We should peek more than once but less than a bunch. This is a rough
		// heuristic for a "reasonable amount" to not CPU saturate.
		require.True(atomic.LoadUint32(&peekCounter) > 0)
		require.True(atomic.LoadUint32(&peekCounter) < 100, "excessive peeking")
		require.True(atomic.LoadUint32(&pollCounter) > 0)
		require.True(atomic.LoadUint32(&pollCounter) < 100, "excessive polling")
	})
}

func TestServicePollQueue(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create a project
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{
					Local: &pb.Job_Local{},
				},
			},
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "15ms",
			},
		}),
	})
	require.NoError(err)

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// Check for our condition, we do eventually here because if we're
	// in a slow environment then this may still be empty.
	require.Eventually(func() bool {
		// We should have a single poll job
		var jobs []*pb.Job
		raw, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
		for _, j := range raw {
			if j.State != pb.Job_ERROR {
				jobs = append(jobs, j)
			}
		}

		if err != nil {
			t.Logf("err: %s", err)
			return false
		}

		return len(jobs) == 1
	}, 5*time.Second, 50*time.Millisecond)

	// Cancel our poller to ensure it stops
	testServiceImpl(impl).Close()

	// Ensure we don't queue more jobs
	time.Sleep(100 * time.Millisecond)
	raw, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
	require.NoError(err)
	time.Sleep(100 * time.Millisecond)
	raw2, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
	require.NoError(err)
	require.Equal(len(raw), len(raw2))
}

func TestProjectPollHandler(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Create a project
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: serverptypes.TestProject(t, &pb.Project{
			Name: "Example",
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{
					Local: &pb.Job_Local{},
				},
			},
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "15ms",
			},
			StatusReportPoll: &pb.Project_AppStatusPoll{
				Enabled: false,
			},
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: "Example"},
					Name:    "apple-app",
				},
			},
		}),
	})
	require.NoError(err)

	// Grab next poll time
	state := testServiceImpl(impl).state(ctx)
	p, pollTime, err := state.ProjectPollPeek(ctx, nil)
	require.NoError(err)
	require.NotNil(p)
	require.NotNil(pollTime)

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// Check for our condition, we do eventually here because if we're
	// in a slow environment then this may still be empty.
	require.Eventually(func() bool {
		// We should have a single poll job
		var jobs []*pb.Job
		raw, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
		for _, j := range raw {
			if j.State != pb.Job_ERROR {
				jobs = append(jobs, j)
			}
		}

		if err != nil {
			t.Logf("err: %s", err)
			return false
		}

		return len(jobs) == 1
	}, 5*time.Second, 50*time.Millisecond)

	// Cancel our poller to ensure it stops
	testServiceImpl(impl).Close()

	// ensure the next poll is after the initial poll before waiting
	// next poll time gets set when a project poll is marked complete
	p, nextPollTime, err := state.ProjectPollPeek(ctx, nil)
	require.NoError(err)
	require.NotNil(p)
	require.NotNil(nextPollTime)
	require.True(nextPollTime.After(pollTime))
}

func TestApplicationPollHandler(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	appName := "apple-app"

	// Create a project with an application
	respProj, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: serverptypes.TestProject(t, &pb.Project{
			Name: "Example",
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{
					Local: &pb.Job_Local{},
				},
			},
			StatusReportPoll: &pb.Project_AppStatusPoll{
				Enabled:  false,
				Interval: "30ms",
			},
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: "Example"},
					Name:    appName,
				},
			},
		}),
	})
	require.NoError(err)
	project := respProj.Project

	// Grab next poll time
	state := testServiceImpl(impl).state(ctx)
	a, _, err := state.ApplicationPollPeek(ctx, nil)
	require.NoError(err)
	require.Nil(a) // Apps Next Poll should be 0 since not started yet

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// Do a deployment
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     "Example",
			},
		}),
	})
	require.NoError(err)
	require.NotNil(resp)

	// Update the app to start polling
	project.StatusReportPoll = &pb.Project_AppStatusPoll{
		Enabled:  true,
		Interval: "30ms",
	}
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: project,
	})
	require.NoError(err)

	// App poll time should be set
	a, pollTime, err := state.ApplicationPollPeek(ctx, nil)
	require.NoError(err)
	require.NotNil(pollTime)
	require.NotNil(a) // Apps Next Poll should be set

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// Check for our condition, we do eventually here because if we're
	// in a slow environment then this may still be empty.
	require.Eventually(func() bool {
		// We should have a single poll job
		var jobs []*pb.Job

		raw, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
		for _, j := range raw {
			if j.State != pb.Job_ERROR && j.SingletonId == appStatusPollSingletonId("default", "Example", appName, appStatusPollOperationTypeDeployment) {
				jobs = append(jobs, j)
			}
		}

		if err != nil {
			t.Logf("err: %s", err)
			return false
		}

		return len(jobs) == 1
	}, 5*time.Second, 50*time.Millisecond)

	// Cancel our poller to ensure it stops
	testServiceImpl(impl).Close()

	// ensure the next poll is after the initial poll before waiting
	// next poll time gets set when a app poll is marked complete
	a, nextPollTime, err := state.ApplicationPollPeek(ctx, nil)
	require.NoError(err)
	require.NotNil(a)
	require.NotNil(nextPollTime)
	t.Logf("nextPollTime %q should be after initial pollTime %q", nextPollTime.String(), pollTime.String())
	require.True(nextPollTime.After(pollTime))
}

func TestApplicationPollHandler_fullLifecycle(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	appName := "apple-app"

	// Create a project with an application
	respProj, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: serverptypes.TestProject(t, &pb.Project{
			Name: "Example",
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{
					Local: &pb.Job_Local{},
				},
			},
			//DataSourcePoll: &pb.Project_Poll{
			//	Enabled:  true,
			//	Interval: "15ms",
			//},
			StatusReportPoll: &pb.Project_AppStatusPoll{
				Enabled:  false,
				Interval: "15ms",
			},
			Applications: []*pb.Application{
				{
					Project: &pb.Ref_Project{Project: "Example"},
					Name:    appName,
				},
			},
		}),
	})
	require.NoError(err)
	project := respProj.Project

	// Grab next poll time
	state := testServiceImpl(impl).state(ctx)
	a, _, err := state.ApplicationPollPeek(ctx, nil)
	require.NoError(err)
	require.Nil(a) // Apps Next Poll should be 0 since not started yet

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// Do a deployment
	deployResp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			Component: &pb.Component{
				Name: "testapp",
			},
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     "Example",
			},
		}),
	})
	require.NoError(err)
	require.NotNil(deployResp)

	// Update the app to start polling
	project.StatusReportPoll = &pb.Project_AppStatusPoll{
		Enabled:  true,
		Interval: "15ms",
	}
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: project,
	})
	require.NoError(err)

	// App poll time should be set
	a, pollTime, err := state.ApplicationPollPeek(ctx, nil)
	require.NoError(err)
	require.NotNil(pollTime)
	require.NotNil(a) // Apps Next Poll should be set

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// Check for our condition, we do eventually here because if we're
	// in a slow environment then this may still be empty.
	require.Eventually(func() bool {
		// We should have a single poll job for just the deployment
		var jobs []*pb.Job
		raw, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
		for _, j := range raw {
			if j.State != pb.Job_ERROR &&
				j.SingletonId == appStatusPollSingletonId("default", "Example", appName, appStatusPollOperationTypeDeployment) {
				// App status polling should only have this singleton id
				jobs = append(jobs, j)
			}
		}

		if err != nil {
			t.Logf("err: %s", err)
			return false
		}

		return len(jobs) == 1
	}, 5*time.Second, 50*time.Millisecond)

	// Do a release
	releaseResp, err := client.UpsertRelease(ctx, &pb.UpsertReleaseRequest{
		Release: serverptypes.TestValidRelease(t, &pb.Release{
			Component: &pb.Component{
				Name: "testapp",
			},
			Application: &pb.Ref_Application{
				Application: appName,
				Project:     "Example",
			},
		}),
	})
	require.NoError(err)
	require.NotNil(releaseResp)

	// Wait a bit. The interval is so low that this should trigger
	// multiple loops through the poller. But we want to ensure we
	// have only one poll job queued.
	time.Sleep(50 * time.Millisecond)

	// Make sure we're polling on two status reports

	// Check for our condition, we do eventually here because if we're
	// in a slow environment then this may still be empty.
	require.Eventually(func() bool {
		// We should have a poll job for a deployment and another for a release

		raw, err := testServiceImpl(impl).state(ctx).JobList(ctx, &pb.ListJobsRequest{})
		releaseJobs := 0
		deployJobs := 0
		for _, j := range raw {
			t.Logf("Found job in state %s with id %s", j.State, j.SingletonId)
			if j.State != pb.Job_ERROR {
				if j.SingletonId == appStatusPollSingletonId("default", "Example", appName, appStatusPollOperationTypeRelease) {
					releaseJobs++
				}
				if j.SingletonId == appStatusPollSingletonId("default", "Example", appName, appStatusPollOperationTypeDeployment) {
					deployJobs++
				}
			}
		}

		if err != nil {
			t.Logf("err: %s", err)
			return false
		}

		return releaseJobs == 1 && deployJobs == 1
	}, 5*time.Second, 50*time.Millisecond)

	// Update the app to stop polling
	project.StatusReportPoll = &pb.Project_AppStatusPoll{
		Enabled: false,
	}
	_, err = client.UpsertProject(ctx, &pb.UpsertProjectRequest{
		Project: project,
	})
	require.NoError(err)

	// App poll time should be unset now
	a, _, err = state.ApplicationPollPeek(ctx, nil)
	require.NoError(err)
	require.Nil(a) // Apps Next Poll should be 0 since turned off

	// Cancel our poller to ensure it stops
	testServiceImpl(impl).Close()
}
