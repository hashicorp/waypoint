package statetest

import (
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["runner"] = []testFunc{
		TestRunner_crud,
		TestRunnerOffline_new,
		TestRunnerAdopt,
		TestRunnerAdopt_changeLabels,
		TestRunnerById_notFound,
	}
}

func TestRunner_crud(t *testing.T, factory Factory, restartF RestartFactory) {
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	// List should be empty
	list, err := s.RunnerList()
	require.NoError(err)
	require.Len(list, 0)

	// Create an instance
	rec := &pb.Runner{Id: "A"}
	require.NoError(s.RunnerCreate(rec))

	// We should be able to find it
	found, err := s.RunnerById(rec.Id, nil)
	require.NoError(err)
	require.Equal(rec.Id, found.Id)
	require.Equal(pb.Runner_NEW, found.AdoptionState)

	// List should include it
	list, err = s.RunnerList()
	require.NoError(err)
	require.Len(list, 1)

	// Delete that instance
	require.NoError(s.RunnerDelete(rec.Id))

	// We should not find it
	found, err = s.RunnerById(rec.Id, nil)
	require.Error(err)
	require.Nil(found)
	require.Equal(codes.NotFound, status.Code(err))

	// List should be empty again
	list, err = s.RunnerList()
	require.NoError(err)
	require.Len(list, 0)

	// Delete again should be fine
	require.NoError(s.RunnerDelete(rec.Id))
}

// New runners that are unadopted should just get deleted when they go offline.
func TestRunnerOffline_new(t *testing.T, factory Factory, restartF RestartFactory) {
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	// List should be empty
	list, err := s.RunnerList()
	require.NoError(err)
	require.Len(list, 0)

	// Create an instance
	rec := &pb.Runner{Id: "A"}
	require.NoError(s.RunnerCreate(rec))

	// List should include it
	list, err = s.RunnerList()
	require.NoError(err)
	require.Len(list, 1)

	// Offline that instance
	require.NoError(s.RunnerOffline(rec.Id))

	// We should not find it
	found, err := s.RunnerById(rec.Id, nil)
	require.Error(err)
	require.Nil(found)
	require.Equal(codes.NotFound, status.Code(err))

	// List should be empty again
	list, err = s.RunnerList()
	require.NoError(err)
	require.Len(list, 0)

	// Delete again should be fine
	require.NoError(s.RunnerDelete(rec.Id))
}

func TestRunnerAdopt(t *testing.T, factory Factory, restartF RestartFactory) {
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	// Create an instance
	rec := &pb.Runner{
		Id: "A",
		Kind: &pb.Runner_Remote_{
			Remote: &pb.Runner_Remote{},
		},
	}
	require.NoError(s.RunnerCreate(rec))

	// Should be new
	ws := memdb.NewWatchSet()
	found, err := s.RunnerById(rec.Id, ws)
	require.NoError(err)
	require.Equal(pb.Runner_NEW, found.AdoptionState)

	// Watch should block
	require.True(ws.Watch(time.After(10 * time.Millisecond)))

	// Adopt that instance
	require.NoError(s.RunnerAdopt(rec.Id, false))

	// Should be triggered. This is a very important test because
	// we need to ensure that the watchers can detect adoption changes.
	require.False(ws.Watch(time.After(100 * time.Millisecond)))

	// Should be adopted
	ws = memdb.NewWatchSet()
	{
		found, err := s.RunnerById(rec.Id, ws)
		require.NoError(err)
		require.Equal(pb.Runner_ADOPTED, found.AdoptionState)
	}

	// Offline that instance, then bring it back.
	require.NoError(s.RunnerOffline(rec.Id))
	require.NoError(s.RunnerCreate(rec))

	// Should be triggered.
	require.False(ws.Watch(time.After(100 * time.Millisecond)))

	// Should still be adopted
	{
		found, err := s.RunnerById(rec.Id, nil)
		require.NoError(err)
		require.Equal(pb.Runner_ADOPTED, found.AdoptionState)
	}

	// Delete that instance, then bring it back.
	require.NoError(s.RunnerDelete(rec.Id))
	require.NoError(s.RunnerCreate(rec))

	// Should NOT be adopted
	{
		found, err := s.RunnerById(rec.Id, nil)
		require.NoError(err)
		require.Equal(pb.Runner_NEW, found.AdoptionState)
	}
}

func TestRunnerAdopt_changeLabels(t *testing.T, factory Factory, restartF RestartFactory) {
	require := require.New(t)

	t.Run("zero to N labels", func(t *testing.T) {
		s := factory(t)
		defer s.Close()

		// Create an instance
		rec := &pb.Runner{
			Id: "A",
			Kind: &pb.Runner_Remote_{
				Remote: &pb.Runner_Remote{},
			},
		}
		require.NoError(s.RunnerCreate(rec))
		require.NoError(s.RunnerAdopt(rec.Id, false))

		// Should be adopted
		{
			found, err := s.RunnerById(rec.Id, nil)
			require.NoError(err)
			require.Equal(pb.Runner_ADOPTED, found.AdoptionState)
		}

		// Offline that instance, then bring it back but with labels.
		require.NoError(s.RunnerOffline(rec.Id))

		// Change labels
		rec.Labels = map[string]string{"A": "B"}
		require.NoError(s.RunnerCreate(rec))

		// Should no longer be adopted
		{
			found, err := s.RunnerById(rec.Id, nil)
			require.NoError(err)
			require.Equal(pb.Runner_NEW, found.AdoptionState)
		}
	})

	t.Run("N to N (matching) labels", func(t *testing.T) {
		s := factory(t)
		defer s.Close()

		// Create an instance
		rec := &pb.Runner{
			Id:     "A",
			Labels: map[string]string{"A": "B"},
			Kind: &pb.Runner_Remote_{
				Remote: &pb.Runner_Remote{},
			},
		}
		require.NoError(s.RunnerCreate(rec))
		require.NoError(s.RunnerAdopt(rec.Id, false))

		// Should be adopted
		{
			found, err := s.RunnerById(rec.Id, nil)
			require.NoError(err)
			require.Equal(pb.Runner_ADOPTED, found.AdoptionState)
		}

		// Offline that instance, then bring it back but with labels.
		require.NoError(s.RunnerOffline(rec.Id))
		require.NoError(s.RunnerCreate(rec))

		// Should no longer be adopted
		{
			found, err := s.RunnerById(rec.Id, nil)
			require.NoError(err)
			require.Equal(pb.Runner_ADOPTED, found.AdoptionState)
		}
	})

	t.Run("N to 0 labels", func(t *testing.T) {
		s := factory(t)
		defer s.Close()

		// Create an instance
		rec := &pb.Runner{
			Id:     "A",
			Labels: map[string]string{"A": "B"},
			Kind: &pb.Runner_Remote_{
				Remote: &pb.Runner_Remote{},
			},
		}
		require.NoError(s.RunnerCreate(rec))
		require.NoError(s.RunnerAdopt(rec.Id, false))

		// Should be adopted
		{
			found, err := s.RunnerById(rec.Id, nil)
			require.NoError(err)
			require.Equal(pb.Runner_ADOPTED, found.AdoptionState)
		}

		// Offline that instance, then bring it back but with labels.
		require.NoError(s.RunnerOffline(rec.Id))
		rec.Labels = nil
		require.NoError(s.RunnerCreate(rec))

		// Should no longer be adopted
		{
			found, err := s.RunnerById(rec.Id, nil)
			require.NoError(err)
			require.Equal(pb.Runner_NEW, found.AdoptionState)
		}
	})
}

func TestRunnerById_notFound(t *testing.T, factory Factory, restartF RestartFactory) {
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	// We should be able to find it
	found, err := s.RunnerById("nope", nil)
	require.Error(err)
	require.Nil(found)
	require.Equal(codes.NotFound, status.Code(err))
}
