package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	/*
		"google.golang.org/grpc/codes"
		"google.golang.org/grpc/status"
	*/

	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestApplication(t *testing.T) {
	t.Run("Put adds a new application", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))

		// Has no apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Empty(resp.Applications)
		}

		// Add
		app, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
		}))
		require.NoError(err)

		// Can read
		{
			resp, err := s.AppGet(&pb.Ref_Application{
				Project:     ref.Project,
				Application: app.Name,
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Has apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 1)
		}
	})

	t.Run("Put non-existent project", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}

		// Add
		app, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
		}))
		require.NoError(err)

		// Can read
		{
			resp, err := s.AppGet(&pb.Ref_Application{
				Project:     ref.Project,
				Application: app.Name,
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Has project
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 1)
		}
	})

	t.Run("Put appends to existing list of applications", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
			Applications: []*pb.Application{
				serverptypes.TestApplication(t, nil),
			},
		})))

		// Add
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    "next",
		}))
		require.NoError(err)

		// Has apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 2)
		}
	})

	t.Run("Put updates an existing application", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Write project
		ref := &pb.Ref_Project{Project: "foo"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
			Applications: []*pb.Application{
				serverptypes.TestApplication(t, &pb.Application{
					Name: "foo",
				}),
			},
		})))

		// Add
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    "foo",
		}))
		require.NoError(err)

		// Has apps
		{
			resp, err := s.ProjectGet(ref)
			require.NoError(err)
			require.NotNil(resp)
			require.Len(resp.Applications, 1)
		}
	})

	t.Run("reads file change signal upward", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		name := "abcde"
		// Set
		err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: name,
		}))
		require.NoError(err)

		_, err = s.AppPut(&pb.Application{
			Project: &pb.Ref_Project{Project: name},
			Name:    "app",
		})
		require.NoError(err)

		err = s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name:             name,
			FileChangeSignal: "HUP",
		}))
		require.NoError(err)

		sig, err := s.GetFileChangeSignal(&pb.Ref_Application{
			Project:     name,
			Application: "app",
		})
		require.NoError(err)

		require.Equal("HUP", sig)

		_, err = s.AppPut(&pb.Application{
			Project:          &pb.Ref_Project{Project: name},
			Name:             "app",
			FileChangeSignal: "TERM",
		})
		require.NoError(err)

		sig, err = s.GetFileChangeSignal(&pb.Ref_Application{
			Project:     name,
			Application: "app",
		})
		require.NoError(err)

		require.Equal("TERM", sig)
	})
}

func TestApplicationPollPeek(t *testing.T) {
	t.Run("returns nil if no values", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		v, _, err := s.ApplicationPollPeek(nil)
		require.NoError(err)
		require.Nil(v)
	})

	t.Run("returns next to poll", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		ref := &pb.Ref_Project{Project: "apple"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    ref.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled:  true,
				Interval: "30s",
			},
		}))
		require.NoError(err)

		// Set another later
		time.Sleep(10 * time.Millisecond)
		refOrg := &pb.Ref_Project{Project: "orange"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: refOrg.Project,
		})))
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    ref.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled:  true,
				Interval: "30s",
			},
		}))
		require.NoError(err)

		// Get exact
		{
			resp, t, err := s.ApplicationPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("apple", resp.Name)
			require.False(t.IsZero())
		}
	})

	t.Run("watchset triggers from empty to available", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ws := memdb.NewWatchSet()
		v, _, err := s.ApplicationPollPeek(ws)
		require.NoError(err)
		require.Nil(v)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Set
		ref := &pb.Ref_Project{Project: "apple"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    ref.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled:  true,
				Interval: "30s",
			},
		}))
		require.NoError(err)

		// Should be triggered.
		require.False(ws.Watch(time.After(100 * time.Millisecond)))

		// Get exact
		{
			resp, t, err := s.ApplicationPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("apple", resp.Name)
			require.False(t.IsZero())
		}
	})

	t.Run("watchset triggers when records change", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		ref := &pb.Ref_Project{Project: "apple"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    ref.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled:  true,
				Interval: "5s",
			},
		}))
		require.NoError(err)

		// Set another later
		refOrg := &pb.Ref_Project{Project: "orange"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: refOrg.Project,
		})))
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: refOrg,
			Name:    refOrg.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled:  true,
				Interval: "5m", // 5 MINUTES, longer than A
			},
		}))
		require.NoError(err)

		// Get applications
		pA, err := s.AppGet(&pb.Ref_Application{Application: "apple", Project: "apple"})
		require.NoError(err)
		require.NotNil(pA)
		pB, err := s.AppGet(&pb.Ref_Application{Application: "orange", Project: "orange"})
		require.NoError(err)
		require.NotNil(pB)

		// Complete both first
		now := time.Now()
		require.NoError(s.ApplicationPollComplete(pA, now))
		require.NoError(s.ApplicationPollComplete(pB, now))

		// Peek, we should get A
		ws := memdb.NewWatchSet()
		p, ts, err := s.ApplicationPollPeek(ws)
		require.NoError(err)
		require.NotNil(p)
		require.Equal("apple", p.Name)
		require.False(ts.IsZero())

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Set
		require.NoError(s.ApplicationPollComplete(pA, now.Add(1*time.Second)))

		// Should be triggered.
		require.False(ws.Watch(time.After(100 * time.Millisecond)))

		// Get exact
		{
			resp, t, err := s.ApplicationPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("apple", resp.Name)
			require.False(t.IsZero())
		}
	})
}

func TestApplicationPollComplete(t *testing.T) {
	t.Run("returns nil for application that doesn't exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		err := s.ApplicationPollComplete(&pb.Application{Name: "NOPE", Project: &pb.Ref_Project{}}, time.Now())
		require.NoError(err)
	})

	t.Run("does nothing for project that has polling disabled", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		ref := &pb.Ref_Project{Project: "apple"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    ref.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled: false,
			},
		}))
		require.NoError(err)

		// Get
		pA, err := s.AppGet(&pb.Ref_Application{Application: "apple", Project: "apple"})
		require.NoError(err)
		require.NotNil(pA)

		// No error
		require.NoError(s.ApplicationPollComplete(pA, time.Now()))

		// Peek does nothing
		v, _, err := s.ApplicationPollPeek(nil)
		require.NoError(err)
		require.Nil(v)
	})

	t.Run("schedules the next poll time", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		ref := &pb.Ref_Project{Project: "apple"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: ref.Project,
		})))
		_, err := s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: ref,
			Name:    ref.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled:  true,
				Interval: "5s",
			},
		}))
		require.NoError(err)

		// Set another later
		refOrg := &pb.Ref_Project{Project: "orange"}
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: refOrg.Project,
		})))
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Project: refOrg,
			Name:    refOrg.Project,
			StatusReportPoll: &pb.Application_Poll{
				Enabled:  true,
				Interval: "5m", // 5 MINUTES, longer than A
			},
		}))
		require.NoError(err)

		// Get applications
		pA, err := s.AppGet(&pb.Ref_Application{Application: "apple", Project: "apple"})
		require.NoError(err)
		require.NotNil(pA)
		pB, err := s.AppGet(&pb.Ref_Application{Application: "orange", Project: "orange"})
		require.NoError(err)
		require.NotNil(pB)

		// Complete both first
		now := time.Now()
		require.NoError(s.ApplicationPollComplete(pA, now))
		require.NoError(s.ApplicationPollComplete(pB, now))

		// Peek should return A, lower interval
		{
			resp, t, err := s.ApplicationPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("apple", resp.Name)
			require.False(t.IsZero())
		}

		// Complete again, a minute later. The result should be A again
		// because of the lower interval.
		{
			require.NoError(s.ApplicationPollComplete(pA, now.Add(1*time.Minute)))

			resp, t, err := s.ApplicationPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("apple", resp.Name)
			require.False(t.IsZero())
		}

		// Complete A, now 6 minutes later. The result should be B now.
		{
			require.NoError(s.ApplicationPollComplete(pA, now.Add(6*time.Minute)))

			resp, t, err := s.ApplicationPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("orange", resp.Name)
			require.False(t.IsZero())
		}
	})
}
