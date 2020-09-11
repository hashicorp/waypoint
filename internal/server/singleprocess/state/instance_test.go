package state

import (
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestInstance_crud(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Create an instance
	rec := &Instance{Id: "A", DeploymentId: "B", Project: "C", Application: "D", Workspace: "E"}
	require.NoError(s.InstanceCreate(rec))

	// We should be able to find it
	found, err := s.InstanceById(rec.Id)
	require.NoError(err)
	require.Equal(rec, found)

	// Delete that instance
	require.NoError(s.InstanceDelete(rec.Id))

	// Delete again should be fine
	require.NoError(s.InstanceDelete(rec.Id))
}

func TestInstanceById_notFound(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// We should be able to find it
	found, err := s.InstanceById("nope")
	require.Error(err)
	require.Nil(found)
	require.Equal(codes.NotFound, status.Code(err))
}

func TestInstancesByApp(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Our ref we'll look for
	ref := &pb.Ref_Application{
		Project:     "A",
		Application: "B",
	}

	// Empty with nothing
	ws := memdb.NewWatchSet()
	list, err := s.InstancesByApp(ref, nil, ws)
	require.NoError(err)
	require.Empty(list)

	// Watch should block
	require.True(ws.Watch(time.After(10 * time.Millisecond)))

	// Create an instance
	rec := testInstance(t, &Instance{Project: ref.Project, Application: ref.Application})
	require.NoError(s.InstanceCreate(rec))

	// Should be triggered
	require.False(ws.Watch(time.After(100 * time.Millisecond)))

	// Should have values
	list, err = s.InstancesByApp(ref, nil, nil)
	require.NoError(err)
	require.Len(list, 1)

	// Should not for other app
	//nolint:copylocks
	ref2 := *ref
	ref2.Application = "NO"
	list, err = s.InstancesByApp(&ref2, nil, nil)
	require.NoError(err)
	require.Empty(list)
}

func TestInstancesByAppWorkspace(t *testing.T) {
	require := require.New(t)

	s := TestState(t)
	defer s.Close()

	// Our ref we'll look for
	ref := &pb.Ref_Application{
		Project:     "A",
		Application: "B",
	}

	refWs := &pb.Ref_Workspace{
		Workspace: "C",
	}

	// Empty with nothing
	ws := memdb.NewWatchSet()
	list, err := s.InstancesByApp(ref, refWs, ws)
	require.NoError(err)
	require.Empty(list)

	// Watch should block
	require.True(ws.Watch(time.After(10 * time.Millisecond)))

	// Create an instance
	rec := testInstance(t, &Instance{
		Project: ref.Project, Application: ref.Application, Workspace: refWs.Workspace})
	require.NoError(s.InstanceCreate(rec))

	// Should be triggered
	require.False(ws.Watch(time.After(100 * time.Millisecond)))

	// Should have values
	list, err = s.InstancesByApp(ref, refWs, nil)
	require.NoError(err)
	require.Len(list, 1)

	// Should not for other app
	//nolint:copylocks
	ref2 := *refWs
	ref2.Workspace = "NO"
	list, err = s.InstancesByApp(ref, &ref2, nil)
	require.NoError(err)
	require.Empty(list)
}

func testInstance(t *testing.T, v *Instance) *Instance {
	if v == nil {
		v = &Instance{}
	}

	require.NoError(t, mergo.Merge(v, &Instance{
		Id:           "A",
		DeploymentId: "B",
		Project:      "C",
		Application:  "D",
		Workspace:    "E",
	}))

	return v
}
