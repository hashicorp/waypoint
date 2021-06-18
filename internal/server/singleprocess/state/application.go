package state

import (
	"time"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

// AppPut creates or updates the application.
func (s *State) AppPut(app *pb.Application) (*pb.Application, error) {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.appPut(dbTxn, memTxn, app)
	})
	if err == nil {
		memTxn.Commit()
	}

	return app, err
}

// AppDelete deletes an application from a project. This will also delete
// all the operations associated with this application.
func (s *State) AppDelete(ref *pb.Ref_Application) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.appDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// AppGet retrieves the application..
func (s *State) AppGet(ref *pb.Ref_Application) (*pb.Application, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Application
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.appGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// AppPollPeek peeks at the next available project that will be polled against,
// and returns the projects application as the result along with the poll time.
// For more information on how ProjectPollPeek works, refer to the ProjectPollPeek
// docs.
func (s *State) ApplicationPollPeek(
	ws memdb.WatchSet,
) (*pb.Application, time.Time, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Application
	var pollTime time.Time
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, pollTime, err = s.appPollPeek(dbTxn, memTxn, ws)
		return err
	})

	return result, pollTime, err
}

// AppPollComplete sets the next poll time for a given project given the app
// reference along with the time interval "t".
func (s *State) ApplicationPollComplete(
	app *pb.Application,
	t time.Time,
) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.View(func(dbTxn *bolt.Tx) error {
		err := s.appPollComplete(dbTxn, memTxn, app, t)
		return err
	})

	return err
}

// GetFileChangeSignal checks the metadata for the given application and its
// project, returning the value of FileChangeSignal that is most relevent.
func (s *State) GetFileChangeSignal(scope *pb.Ref_Application) (string, error) {
	app, err := s.AppGet(scope)
	if err != nil {
		return "", err
	}

	if app.FileChangeSignal != "" {
		return app.FileChangeSignal, nil
	}

	project, err := s.ProjectGet(&pb.Ref_Project{Project: scope.Project})
	if err != nil {
		return "", err
	}

	return project.FileChangeSignal, nil
}

func (s *State) appPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Application,
) error {
	// Get the project
	p, err := s.projectGetOrCreate(dbTxn, memTxn, value.Project)
	if err != nil {
		return err
	}

	// If we have a matching app, then modify that that.
	pt := &serverptypes.Project{Project: p}
	if idx := pt.App(value.Name); idx >= 0 {
		p.Applications[idx] = value
		value = nil
	}

	// If we didn't have a matching app, insert it
	if value != nil {
		p.Applications = append(p.Applications, value)
	}

	return s.projectPut(dbTxn, memTxn, p)
}

func (s *State) appDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Application,
) error {
	// Get the project
	p, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: ref.Project,
	})
	if err != nil {
		// If the project doesn't exist then the app is deleted.
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// If we have a matching app, then modify that that.
	pt := &serverptypes.Project{Project: p}
	if i := pt.App(ref.Application); i >= 0 {
		s := p.Applications
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		p.Applications = s[:len(s)-1]
	}

	return nil
}

func (s *State) appGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Application,
) (*pb.Application, error) {
	// Get the project
	p, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: ref.Project,
	})
	if err != nil {
		return nil, err
	}

	// If we have a matching app, then modify that that.
	pt := &serverptypes.Project{Project: p}
	if i := pt.App(ref.Application); i >= 0 {
		return p.Applications[i], nil
	}

	return nil, status.Errorf(codes.NotFound, "application not found")
}

// appDefaultForRef returns a default pb.Application for a ref. This
// can be used in tandem with appCreateIfNotExist to create defaults.
func (s *State) appDefaultForRef(ref *pb.Ref_Application) *pb.Application {
	return &pb.Application{
		Name: ref.Application,
		Project: &pb.Ref_Project{
			Project: ref.Project,
		},
	}
}

func (s *State) appPollPeek(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
) (*pb.Application, time.Time, error) {
	// LowerBound doesn't support watches so we have to do a Get first
	// to get a valid watch channel on these fields.
	iter, err := memTxn.Get(
		projectIndexTableName,
		applIndexNextPollIndexName,
		true,            // polling enabled
		time.Unix(0, 0), // lowest next poll time
	)
	if err != nil {
		return nil, time.Time{}, err
	}
	ws.Add(iter.WatchCh())

	// Get the project with the lowest "next poll" time.
	iter, err = memTxn.LowerBound(
		projectIndexTableName,
		applIndexNextPollIndexName,
		true,            // polling enabled
		time.Unix(0, 0), // lowest next poll time
	)
	if err != nil {
		return nil, time.Time{}, err
	}

	// If we have no values, then return
	raw := iter.Next()
	if raw == nil {
		return nil, time.Time{}, nil
	}

	rec := raw.(*projectIndexRecord)
	if rec.ApplNextPoll.IsZero() {
		// This _shouldnt_ happen but let's protect against it anyways.
		return nil, time.Time{}, nil
	}

	var result *pb.Project
	err = s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
			Project: rec.Id,
		})

		return err
	})

	// TODO(briancain): what about projects that define multiple apps? So far I
	// don't think this workflow is really supported in Waypoint
	return result.Applications[0], rec.ApplNextPoll, err
}

func (s *State) appPollComplete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	app *pb.Application,
	t time.Time,
) error {
	// Get the project
	p, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: app.Project.Project,
	})
	if status.Code(err) == codes.NotFound {
		return nil
	}
	if err != nil {
		return err
	}

	raw, err := memTxn.First(
		projectIndexTableName,
		projectIndexIdIndexName,
		string(s.projectId(p)),
	)
	if err != nil {
		return err
	}
	if raw == nil {
		return nil
	}

	record := raw.(*projectIndexRecord)
	if !record.ApplPoll {
		// If this project doesn't have polling enabled, then do nothing.
		// This could happen if a project had polling when Peek was called,
		// then between Peek and Complete, polling was disabled.
		return nil
	}

	record = record.Copy()
	record.ApplLastPoll = t
	record.ApplNextPoll = t.Add(record.ApplPollInterval)
	if err := memTxn.Insert(projectIndexTableName, record); err != nil {
		return err
	}

	memTxn.Commit()
	return nil
}
