package state

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestJobsPrune(t *testing.T) {
	t.Run("removes only completed jobs", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Cancel it
		require.NoError(s.JobCancel("A", false))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		// Leave B running

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 0)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(1, cnt)
		require.Equal(1, s.indexedJobs)

		val, err := s.JobById("A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById("B", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("does nothing there are fewer than the maximum", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		// Create a build
		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		// Cancel it
		require.NoError(s.JobCancel("A", false))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		// Leave B running

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		require.Equal(2, s.indexedJobs)
		cnt, err := s.jobsPruneOld(memTxn, 10)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(0, cnt)
		require.Equal(2, s.indexedJobs)

		val, err := s.JobById("A", nil)
		require.NoError(err)
		require.NotNil(val)

		val, err = s.JobById("B", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("stops when the maximum are pruned", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		require.NoError(s.JobCancel("A", false))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		require.NoError(s.JobCancel("B", false))

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 1)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(1, cnt)
		require.Equal(1, s.indexedJobs)

		val, err := s.JobById("A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById("B", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("prunes according to the queue time", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		require.NoError(s.JobCancel("A", false))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		require.NoError(s.JobCancel("B", false))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
		})))

		require.NoError(s.JobCancel("C", false))

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 1)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(2, cnt)
		require.Equal(1, s.indexedJobs)

		val, err := s.JobById("A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById("B", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById("C", nil)
		require.NoError(err)
		require.NotNil(val)
	})

	t.Run("can prune all jobs", func(t *testing.T) {
		require := require.New(t)

		s := TestState(t)
		defer s.Close()

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "A",
		})))

		require.NoError(s.JobCancel("A", false))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "B",
		})))

		require.NoError(s.JobCancel("B", false))

		require.NoError(s.JobCreate(serverptypes.TestJobNew(t, &pb.Job{
			Id: "C",
		})))

		require.NoError(s.JobCancel("C", false))

		memTxn := s.inmem.Txn(true)
		defer memTxn.Abort()

		cnt, err := s.jobsPruneOld(memTxn, 0)
		require.NoError(err)

		memTxn.Commit()

		require.Equal(3, cnt)
		require.Equal(0, s.indexedJobs)

		val, err := s.JobById("A", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById("B", nil)
		require.NoError(err)
		require.Nil(val)

		val, err = s.JobById("C", nil)
		require.NoError(err)
		require.Nil(val)
	})
}
