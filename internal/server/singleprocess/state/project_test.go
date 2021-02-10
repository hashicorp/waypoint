package state

import (
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestProject(t *testing.T) {
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		_, err := s.ProjectGet(&pb.Ref_Project{
			Project: "foo",
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "AbCdE",
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get case insensitive
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "abcDe",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// List
		{
			resp, err := s.ProjectList()
			require.NoError(err)
			require.Len(resp, 1)
		}
	})

	t.Run("Put does not modify applications", func(t *testing.T) {
		require := require.New(t)

		const name = "AbCdE"
		ref := &pb.Ref_Project{Project: name}

		s := TestState(t)
		defer s.Close()

		// Set
		proj := serverptypes.TestProject(t, &pb.Project{Name: name})
		err := s.ProjectPut(proj)
		require.NoError(err)
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Name:    "test",
			Project: ref,
		}))
		require.NoError(err)
		_, err = s.AppPut(serverptypes.TestApplication(t, &pb.Application{
			Name:    "test2",
			Project: ref,
		}))
		require.NoError(err)

		// Get exact
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.False(resp.RemoteEnabled)
			require.Len(resp.Applications, 2)
		}

		// Update the project
		proj.RemoteEnabled = true
		require.NoError(s.ProjectPut(proj))

		// Get exact
		{
			resp, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.True(resp.RemoteEnabled)
			require.Len(resp.Applications, 2)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		err := s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "AbCdE",
		}))
		require.NoError(err)

		// Read
		resp, err := s.ProjectGet(&pb.Ref_Project{
			Project: "AbCdE",
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.ProjectDelete(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.ProjectGet(&pb.Ref_Project{
				Project: "AbCdE",
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.ProjectList()
			require.NoError(err)
			require.Len(resp, 0)
		}
	})
}

func TestProjectPollPeek(t *testing.T) {
	t.Run("returns nil if no values", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		v, err := s.ProjectPollPeek(nil)
		require.NoError(err)
		require.Nil(v)
	})

	t.Run("returns next to poll", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "10s",
			},
		})))

		// Set another later
		time.Sleep(10 * time.Millisecond)
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "B",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "10s",
			},
		})))

		// Get exact
		{
			resp, err := s.ProjectPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
		}
	})

	t.Run("watchset triggers from empty to available", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		ws := memdb.NewWatchSet()
		v, err := s.ProjectPollPeek(ws)
		require.NoError(err)
		require.Nil(v)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Set
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "10s",
			},
		})))

		// Should be triggered.
		require.False(ws.Watch(time.After(100 * time.Millisecond)))

		// Get exact
		{
			resp, err := s.ProjectPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
		}
	})

	t.Run("watchset triggers when records change", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5s",
			},
		})))

		// Set another later
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "B",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5m", // 5 MINUTES, longer than A
			},
		})))

		// Get
		pA, err := s.ProjectGet(&pb.Ref_Project{Project: "A"})
		require.NoError(err)
		require.NotNil(pA)
		pB, err := s.ProjectGet(&pb.Ref_Project{Project: "B"})
		require.NoError(err)
		require.NotNil(pB)

		// Complete both first
		now := time.Now()
		require.NoError(s.ProjectPollComplete(pA, now))
		require.NoError(s.ProjectPollComplete(pB, now))

		// Peek, we should get A
		ws := memdb.NewWatchSet()
		p, err := s.ProjectPollPeek(ws)
		require.NoError(err)
		require.NotNil(p)
		require.Equal("A", p.Name)

		// Watch should block
		require.True(ws.Watch(time.After(10 * time.Millisecond)))

		// Set
		require.NoError(s.ProjectPollComplete(pA, now.Add(1*time.Second)))

		// Should be triggered.
		require.False(ws.Watch(time.After(100 * time.Millisecond)))

		// Get exact
		{
			resp, err := s.ProjectPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
		}
	})
}

func TestProjectPollComplete(t *testing.T) {
	t.Run("returns nil for project that doesn't exist", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		require.NoError(s.ProjectPollComplete(&pb.Project{Name: "NOPE"}, time.Now()))
	})

	t.Run("does nothing for project that has polling disabled", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled: false,
			},
		})))

		// Get
		p, err := s.ProjectGet(&pb.Ref_Project{
			Project: "A",
		})
		require.NoError(err)
		require.NotNil(p)

		// No error
		require.NoError(s.ProjectPollComplete(p, time.Now()))

		// Peek does nothing
		v, err := s.ProjectPollPeek(nil)
		require.NoError(err)
		require.Nil(v)
	})

	t.Run("schedules the next poll time", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Set
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "A",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5s",
			},
		})))

		// Set another later
		require.NoError(s.ProjectPut(serverptypes.TestProject(t, &pb.Project{
			Name: "B",
			DataSourcePoll: &pb.Project_Poll{
				Enabled:  true,
				Interval: "5m", // 5 MINUTES, longer than A
			},
		})))

		// Get
		pA, err := s.ProjectGet(&pb.Ref_Project{Project: "A"})
		require.NoError(err)
		require.NotNil(pA)
		pB, err := s.ProjectGet(&pb.Ref_Project{Project: "B"})
		require.NoError(err)
		require.NotNil(pB)

		// Complete both first
		now := time.Now()
		require.NoError(s.ProjectPollComplete(pA, now))
		require.NoError(s.ProjectPollComplete(pB, now))

		// Peek should return A, lower interval
		{
			resp, err := s.ProjectPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
		}

		// Complete again, a minute later. The result should be A again
		// because of the lower interval.
		{
			require.NoError(s.ProjectPollComplete(pA, now.Add(1*time.Minute)))

			resp, err := s.ProjectPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("A", resp.Name)
		}

		// Complete A, now 6 minutes later. The result should be B now.
		{
			require.NoError(s.ProjectPollComplete(pA, now.Add(6*time.Minute)))

			resp, err := s.ProjectPollPeek(nil)
			require.NoError(err)
			require.NotNil(resp)
			require.Equal("B", resp.Name)
		}
	})
}
